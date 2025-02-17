package main

import (
	filedownloader "fan-ot-fan-in/file_downloader"
	fileprocessor "fan-ot-fan-in/file_processor"
	"log"
	"runtime"
	"time"
)

func main() {
	fileUrls := []string{
		"http://example.com/file1",
		"http://example.com/file2",
		"http://example.com/file3",
		"http://example.com/file4",
		"http://example.com/file5",
		"http://example.com/file6",
		"http://example.com/file7",
		"http://example.com/file8",
		"http://example.com/file9",
		"http://example.com/file10",
	}
	done := make(chan bool)
	defer close(done)
	startSynchronousProcessingStage(done, fileUrls)
	startAsynchronousProcessingStage(done, fileUrls)
}

func startAsynchronousProcessingStage(done <-chan bool, fileUrls []string) {
	log.Println("Starting asynchronously processing stage")
	startTime := time.Now()
	urlStream := fileUrlStreamGenerator(fileUrls)

	numberOfWorkers := 4
	runtime.GOMAXPROCS(numberOfWorkers)
	downloadStreamArr := make([]<-chan []byte, numberOfWorkers)
	// Fan-Out
	for i := 0; i < numberOfWorkers; i++ {
		downloadStreamArr[i] = filedownloader.DownloadFiles(done, urlStream)
	}
	// Fan-In
	mergedStream := fileprocessor.Merge(done, downloadStreamArr...)
	zipFileName := <-fileprocessor.ProcessContent(done, mergedStream)
	log.Println("Zip file created: ", zipFileName)
	log.Println("Time taken in asynchronously processing stage: ", time.Since(startTime))

}
func startSynchronousProcessingStage(done <-chan bool, fileUrls []string) {
	log.Println("Starting synchronous processing stage")
	startTime := time.Now()
	urlStream := fileUrlStreamGenerator(fileUrls)

	downloadFileStream := filedownloader.DownloadFiles(done, urlStream)
	zipFileName := <-fileprocessor.ProcessContent(done, downloadFileStream)
	log.Println("Zip file created: ", zipFileName)
	log.Println("Time taken in synchronously processing stage: ", time.Since(startTime))
}
func fileUrlStreamGenerator(fileUrls []string) <-chan string {
	fileUrlStream := make(chan string, len(fileUrls))
	go func() {
		defer close(fileUrlStream)
		for _, fileUrl := range fileUrls {
			fileUrlStream <- fileUrl
		}
	}()
	return fileUrlStream
}
