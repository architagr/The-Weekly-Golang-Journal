package fileprocessor

import (
	"log"
	"sync"
	"time"
)

func ProcessContent(done <-chan bool, fileStream <-chan []byte) <-chan string {
	contentStream := make(chan string)

	go func() {
		defer close(contentStream)
		outFileName := "output.zip"
		// create a zip archive

		for fileContent := range fileStream {
			// add file to a zip archive
			time.Sleep(1 * time.Second)
			log.Printf("Adding file %s to zip archive: %s \n", fileContent, outFileName)
			// write file to zip
		}
		contentStream <- outFileName
	}()
	return contentStream
}

func Merge(done <-chan bool, contentStreams ...<-chan []byte) <-chan []byte {
	mergedStream := make(chan []byte)
	wg := sync.WaitGroup{}
	// This function takes single content stream and multiplexes it to mergedStream
	multiplex := func(done <-chan bool, contentStream <-chan []byte) {
		defer wg.Done()
		for content := range contentStream {
			// merge content
			select {
			case <-done:
				return
			case mergedStream <- content:
			}
		}
	}
	wg.Add(len(contentStreams))
	for _, contentStream := range contentStreams {
		// Fan-In
		// here we are starting multiple goroutines to multiplex content from multiple content streams to mergedStream
		go multiplex(done, contentStream)
	}
	// Wait for all multiplexing goroutines to finish
	go func() {
		wg.Wait()
		close(mergedStream)
	}()
	return mergedStream
}
