package filedownloader

import (
	"log"
	"time"
)

func DownloadFiles(done <-chan bool, fileUrlsStream <-chan string) <-chan []byte {
	fileStream := make(chan []byte, 10)

	go func() {
		defer close(fileStream)
		for url := range fileUrlsStream {
			log.Println("Downloading file from url: ", url)
			// download file
			time.Sleep(2 * time.Second)
			fileStream <- []byte("file content " + url)
		}
	}()
	return fileStream
}
