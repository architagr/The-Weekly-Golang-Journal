package unidirectional

import (
	"fmt"
	"sync"
	processor "sync_pool/message_processor"
	stringsanatization "sync_pool/string_sanatization"
)

type messageProcessor interface {
	Push(data string)
}

func messageProducer(wg *sync.WaitGroup) <-chan string {
	data := []string{
		"data 1- test  data",
		"data 2- test  \tdata",
		"data 3- test \t \tdata",
		"data 4- test  \t \t  data",
		"data 5- test  \t \t  data   ",
		"   data 6- test  \t \t  data  ",
	}
	sanatizationStream := make(chan string, 3)
	wg.Add(1)
	go func(inputStream chan<- string) {
		defer wg.Done()
		defer close(inputStream)
		for _, dataStr := range data {
			inputStream <- dataStr
			fmt.Println("messageProducer pushed message:", dataStr)
		}
		fmt.Println("messageProducer: closing channel")
	}(sanatizationStream)
	return sanatizationStream
}

func dataSanitization(consumerId int, inputStream <-chan string, processor messageProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range inputStream {
		fmt.Printf("dataSanitization(%d) received message: %s\n", consumerId, data)
		sanatizedData := stringsanatization.Sanatize(data)
		processor.Push(sanatizedData)
	}
}

func Run() {
	dbProcessor := processor.InitDatabaseMessageProcessor("dbConn string")
	wg := new(sync.WaitGroup)
	sanatizationStream := messageProducer(wg)
	numberOfConsumers := 2
	for i := 1; i <= numberOfConsumers; i++ {
		wg.Add(1)
		go dataSanitization(i, sanatizationStream, dbProcessor, wg)
	}
	wg.Wait()
}
