package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/kru/soundify/core/database"
	"github.com/kru/soundify/core/helper"
)

type requestBody struct {
	Link string `json:"link"`
}

func HandleLink(w http.ResponseWriter, r *http.Request) {
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

	htmlSource, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Can not process youtube page", http.StatusNotFound)
		return
	}

	// Step 2: Extract the video URL (simplified)
	// This is a basic regex and might not work for all cases
	// pattern := `"url":"(https:[^"]+googlevideo\.com[^"]+)"`
	pattern := `"url":"(https://[^"]+googlevideo\.com[^"]+videoplayback[^"]+)"`
	// 	re := regexp.MustCompile(`"url":"(https:[^"]+googlevideo\.com[^"]+)"`)
	// matches := re.FindStringSubmatch(html)
	// if len(matches) > 1 {
	// 	unescapedURL := strings.ReplaceAll(matches[1], `\u0026`, "&")
	// 	return unescapedURL, nil
	// }
	// return "", fmt.Errorf("no download URL found")
	// pattern := `"url":"(https://[^"]*googlevideo\.com[^"]*)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(htmlSource))
	if len(matches) < 2 {
		fmt.Println("Could not find video URL")
		return
	}
	videoDownloadURL := matches[1]

	// You may need to decode URL-encoded characters
	videoDownloadURL = helper.Unescape(videoDownloadURL)

	re2 := regexp.MustCompile(`\"title\":\"([^\"]{0,255})`)
	matches2 := re2.FindStringSubmatch(string(htmlSource))

	var videoTitle = ""
	if len(matches2) < 2 {
		fmt.Println("Could not find video title, fallback to default naming")
		videoTitle = videoDownloadURL[10:12]
	}
	videoTitle = matches2[1]
	fmt.Println("Video title: ", videoTitle)

	fileID, err := database.CreateFile(user.Id, reqBody.Link, videoDownloadURL, videoTitle)
	if err != nil {
		http.Error(w, "Unable to process given url", http.StatusInternalServerError)
		return
	}
	log.Println("fileID: ", fileID)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response": "We're processing your request, you can download it via email later"}`))
}
