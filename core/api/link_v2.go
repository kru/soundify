package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/kru/soundify/core/database"
)

type PlayerResponse struct {
	StreamingData struct {
		Formats []struct {
			URL string `json:"url"`
		} `json:"formats"`
	} `json:"streamingData"`
}

// Extract the player response JSON from the HTML
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func writeToFile(str string) {
	data := []byte(str)

	fileName := "matches.txt"

	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}

	fmt.Println("sucessfully wrote data to", fileName)

}

func extractPlayerResponse(html string) (string, error) {

	re := regexp.MustCompile(`ytInitialPlayerResponse\s*=\s*({.+?});`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		writeToFile(matches[1])
		return matches[1], nil
	}

	return "", fmt.Errorf("no player response found")
}

func parsePlayerResponse(playerResponse string) ([]string, error) {
	var response PlayerResponse
	err := json.Unmarshal([]byte(playerResponse), &response)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, format := range response.StreamingData.Formats {
		if format.URL != "" {
			urls = append(urls, format.URL)
		}
	}

	return urls, nil
}

func HandleLinkV2(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(database.UserContextKey).(*database.User)
	fmt.Printf("user %+v\n", user)
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

	html, err := fetchHTML(reqBody.Link)
	if err != nil {
		log.Printf("failed to fetch the YouTube page: %v", err)
		return
	}

	playerResponse, err := extractPlayerResponse(html)
	if err != nil {
		log.Printf("failed to extract the player response: %v", err)
		return
	}

	downloadURLs, err := parsePlayerResponse(playerResponse)
	if err != nil {
		log.Printf("failed to parse the player response: %v", err)
		return
	}

	if len(downloadURLs) == 0 {
		fmt.Println("no downloadable URLs found")
	} else {
		fmt.Println("Download URLs:")
		for _, url := range downloadURLs {
			fmt.Println(url)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"response": "we are processing"}`))
}
