package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kru/soundify/core/database"
)

func Run() {
	fmt.Println("init worker")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("worker started")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		go func(ctx context.Context) {
			files, err := database.GetFiles()
			if err != nil {
				log.Fatalf("error while querying files table %v", err)
			}

			// TODO: Before processing file check expiry time in the url, if expire fetch new url
			for i := 0; i < len(files); i++ {
				file := files[i]
				fmt.Println("processing file: ", file.Name)

				fileName := strings.ReplaceAll(file.Name, " ", "-")
				outFile, err := os.Create(fmt.Sprintf("%s.mp4", fileName))
				if err != nil {
					fmt.Println("error while creating file", err)
					continue
				}

				defer outFile.Close()

				err = database.UpdateFile(database.Processing, time.Now().Format(time.RFC3339), file.Id)
				if err != nil {
					fmt.Println("error while updating status to processing", err)
					continue
				}

				video, err := http.Get(file.DownloadUrl.String)
				if err != nil {
					fmt.Println("error downloading video:", err)
					continue
				}
				defer video.Body.Close()

				_, err = io.Copy(outFile, video.Body)
				if err != nil {
					fmt.Println("error saving video:", err)
					continue
				}
				// update status to processing or processed
				database.UpdateFile(database.Processed, time.Now().Format(time.RFC3339), file.Id)

				//TODO: upload video to external bucket, maybe S3, then save the link to DB to audio_url
				//TODO: after saved, send the audio_url via email to the user
				//TODO: save extention to create file using mimetype while parsing the video download url
			}

			select {
			case <-ctx.Done():
				fmt.Println("Worker timeout after 15 minutes")
				return
			default:
				fmt.Println("Worker completed successfully")
			}
		}(ctx)

	}
}
