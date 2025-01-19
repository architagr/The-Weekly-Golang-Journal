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
		if response.Err != nil {
			fmt.Println("error in downloding the file, Error is:", response.Err)
		} else {
			fmt.Printf("Response status: %v\n", response.Response.Status)
			response.Response.Body.Close()
		}
	}
}

type CustomDownloadResponse struct {
	Response *http.Response
	Err      error
}

func downloadMultipleFiles(done <-chan interface{}, fileUrls ...string) <-chan *CustomDownloadResponse {
	saveFiles := make(chan *CustomDownloadResponse)
	go func() {
		defer close(saveFiles)
		for _, file := range fileUrls {
			downloadFile(done, file, saveFiles)
		}
	}()
	return saveFiles
}
func downloadFile(done <-chan interface{}, url string, responseChan chan<- *CustomDownloadResponse) {
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("download status: %v", resp.StatusCode)
	}
	result := &CustomDownloadResponse{
		Response: resp,
		Err:      err,
	}
	select {
	case <-done:
		return
	default:
		responseChan <- result
	}
}
