package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/kru/soundify/core/database"
	"github.com/kru/soundify/core/helper"
	"github.com/kru/soundify/core/middleware"
)

type requestBody struct {
	Link string `json:"link"`
}

func handleLink(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(database.UserContextKey).(*database.User)
	if !ok {
		log.Println("invalid user")
		w.WriteHeader(http.StatusBadRequest)
	}

	var reqBody requestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Incoming request body: %+v\n", reqBody)
	resp, err := http.Get(string(reqBody.Link))
	if err != nil {
		http.Error(w, "Youtube page not found", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Can not process youtube page", http.StatusNotFound)
		return
	}

	// Step 2: Extract the video URL (simplified)
	// This is a basic regex and might not work for all cases
	re := regexp.MustCompile(`"url":"(https://[^"]+)"`)
	matches := re.FindStringSubmatch(string(htmlData))
	if len(matches) < 2 {
		fmt.Println("Could not find video URL")
		return
	}
	videoDownloadURL := matches[1]

	// You may need to decode URL-encoded characters
	videoDownloadURL = helper.Unescape(videoDownloadURL)

	fileID, err := database.CreateFile(user.Id, reqBody.Link, videoDownloadURL)
	if err != nil {
		http.Error(w, "Unable to process given url", http.StatusInternalServerError)
		return
	}
	log.Println("fileID: ", fileID)

	// Step 3: Download the video
	// put inside worker
	// outFile, err := os.Create("kris-test.mp4")
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)
	// 	return
	// }
	// defer outFile.Close()
	//
	// videoResp, err := http.Get(videoDownloadURL)
	// if err != nil {
	// 	fmt.Println("Error downloading video:", err)
	// 	return
	// }
	// defer videoResp.Body.Close()
	//
	// _, err = io.Copy(outFile, videoResp.Body)
	// if err != nil {
	// 	fmt.Println("Error saving video:", err)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response": "We're processing your request, you can download it via email later"}`))
}

func main() {

	err := helper.LoadEnv(".env")
	if err != nil {
		log.Fatalf("error while loading env file %v", err)
		return
	}
	// start db connection
	database.Init()
	defer database.DB.Close()

	if err != nil {
		log.Fatalf("error while querying users %v", err)
		return
	}

	router := http.NewServeMux()

	router.HandleFunc("POST /links", handleLink)

	middlewares := middleware.CombineMiddleware(
		middleware.Logger,
		middleware.Auth,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
	}

	fmt.Println("Server listening to port 8080")
	log.Fatal(server.ListenAndServe())
}
