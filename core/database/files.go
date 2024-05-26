package database

import "log"

type FileStatus string

// Define the possible statuses
const (
	New        FileStatus = "new"
	Processing FileStatus = "processing"
	Processed  FileStatus = "processed"
	Failed     FileStatus = "failed"
)

func CreateFile(userID int, sourceURL string, downloadURL string) (int, error) {
	var fileID int
	err := DB.QueryRow(
		"INSERT INTO files (user_id, source_url, download_url, status) VALUES ($1, $2, $3, $4) RETURNING id",
		userID, sourceURL, downloadURL, New,
	).Scan(&fileID)

	if err != nil {
		log.Printf("CreateFile err: %v", err)
		return 0, err
	}

	return fileID, nil
}
