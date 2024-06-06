package database

import (
	"database/sql"
	"fmt"
	"log"
)

type FileStatus string

type File struct {
	Id          int            `json:"id"`
	UserID      int            `json:"user_id"`
	Status      string         `json:"status"`
	SourceUrl   string         `json:"source_url"`
	AudioUrl    sql.NullString `json:"audio_url"`
	DownloadUrl sql.NullString `json:"download_url"`
	Name        string         `json:"name"`
	CreatedAt   string         `json:"created_at"`
}

// Define the possible statuses
const (
	New        FileStatus = "new"
	Processing FileStatus = "processing"
	Processed  FileStatus = "processed"
	Failed     FileStatus = "failed"
)

func CreateFile(userID int, sourceURL string, downloadURL string, name string) (int, error) {
	var fileID int
	err := DB.QueryRow(
		"INSERT INTO files (user_id, source_url, download_url, status, name) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userID, sourceURL, downloadURL, New, name,
	).Scan(&fileID)

	if err != nil {
		log.Printf("CreateFile err: %v", err)
		return 0, err
	}

	return fileID, nil
}

func GetFiles() ([]File, error) {
	var file File
	var files []File
	rows, err := DB.Query(`SELECT id, user_id, status, 
		source_url, audio_url, download_url, 
		name, created_at FROM files WHERE status = $1`, New)
	if err != nil {
		return files, fmt.Errorf("Err executing select query on files: %v\n", err)
	}

	for rows.Next() {
		err := rows.Scan(&file.Id, &file.UserID, &file.Status,
			&file.SourceUrl, &file.AudioUrl, &file.DownloadUrl,
			&file.Name, &file.CreatedAt)
		if err != nil {
			return files, fmt.Errorf("Err scanning next files row: %v\n", err)
		}
		files = append(files, file)
	}

	fmt.Printf("files: %+v\n", files)

	// Check for errors from iterating over rows
	err = rows.Err()
	if err != nil {
		return files, fmt.Errorf("Err iterating over rows: %v\n", err)
	}

	return files, nil
}

func UpdateFile(status FileStatus, updatedAt string, fileID int) error {
	result, err := DB.Exec(
		"UPDATE files SET status = $1, updated_at = $2 WHERE id = $3",
		status, updatedAt, fileID,
	)

	if err != nil {
		return fmt.Errorf("Error while updating file ID: %d error %v", fileID, err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error file ID: %d not affecting any row %v", fileID, err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("no record found with ID %d", fileID)
	}

	return nil
}
