package main

import (
	"fmt"
	"net/http"
)

func main() {

	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://example.com/file.txt", "htp://a"}

	for response := range downloadMultipleFiles(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
		// save the file to disk or send it to another channel to be processed.
	}

}

func downloadMultipleFiles(done <-chan interface{}, fileUrls ...string) <-chan *http.Response {
	saveFiles := make(chan *http.Response)
	go func() {
		defer close(saveFiles)
		for _, file := range fileUrls {
			downloadFile(done, file, saveFiles)
		}
	}()
	return saveFiles
}
func downloadFile(done <-chan interface{}, url string, responseChan chan<- *http.Response) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error reaching server: %v\n", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return
	}
	select {
	case <-done:
		return
	case responseChan <- resp:
	}

}
