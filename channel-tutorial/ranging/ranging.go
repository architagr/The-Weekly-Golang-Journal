package ranging

import (
	processor "channel-tutorial/message_processor"
	stringsanatization "channel-tutorial/string_sanatization"
	"fmt"
	"os"
	"sync"
)

type messageProcessor interface {
	Push(data string)
}

func Run() {
	stdPrintProcessor := processor.InitPrintMessage(os.Stdout)
	sanatizationStream := make(chan string, 3)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go messageProducer(sanatizationStream, wg)

	numberOfConsumers := 2
	for i := 1; i <= numberOfConsumers; i++ {
		wg.Add(1)
		go dataSanitization(i, sanatizationStream, stdPrintProcessor, wg)
	}
	wg.Wait()
}

func messageProducer(inputStream chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	data := []string{
		"data 1- test  data",
		"data 2- test  \tdata",
		"data 3- test \t \tdata",
		"data 4- test  \t \t  data",
		"data 5- test  \t \t  data   ",
		"   data 6- test  \t \t  data  ",
	}

	for _, dataStr := range data {
		inputStream <- dataStr
		fmt.Println("messageProducer pushed message:", dataStr)
	}
	fmt.Println("messageProducer: closing channel")
	close(inputStream)
}
func dataSanitization(consumerId int, inputStream chan string, processor messageProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range inputStream {
		fmt.Printf("dataSanitization(%d) received message: %s\n", consumerId, data)
		sanatizedData := stringsanatization.Sanatize(data)
		processor.Push(sanatizedData)
	}
}
