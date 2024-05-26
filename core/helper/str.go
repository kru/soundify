package helper

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// Helper function to unescape URL-encoded characters
func Unescape(s string) string {
	repl := map[string]string{
		`\u0026`:  "&",
		`\\u0026`: "&",
	}
	for old, new := range repl {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}

func GenerateRandomStr(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Convert bytes to base64 string
	randomString := base64.URLEncoding.EncodeToString(bytes)
	// Trim any padding characters from the string
	randomString = randomString[:length]

	return randomString, nil
}
