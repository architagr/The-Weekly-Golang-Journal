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
	sanatizationStream := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go messageProducer(sanatizationStream, wg)
	wg.Add(1)
	go dataSanitization(sanatizationStream, stdPrintProcessor, wg)
	wg.Wait()
}

func messageProducer(inputStream chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	data := []string{
		"data 1- test  data",
		"data 2- test  data\t with new line ",
		"data 3- test  data\t with\tnew line and tab ",
		"data 4-   test  data\t with\tnew line and tab    ",
		"data 5-   test  213v data\t with\tnew line and tab    234234   ",
		"data 6-   test  213v data\t with\tnew line and tab    234234   test",
	}

	for _, dataStr := range data {
		fmt.Println("messageProducer:", dataStr)
		inputStream <- dataStr
	}
	fmt.Println("messageProducer: closing channel")
	close(inputStream)
}
func dataSanitization(inputStream chan string, processor messageProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range inputStream {
		fmt.Println("dataSanitization:", data)
		sanatizedData := stringsanatization.Sanatize(data)
		processor.Push(sanatizedData)
	}
}
