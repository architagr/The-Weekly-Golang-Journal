package singlemessage

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
	// close(sanatizationStream)
}

func messageProducer(inputStream chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	data := "  test  213v data with\tnew line and tab    234234   "
	fmt.Println("messageProducer:", data)
	inputStream <- data
}
func dataSanitization(inputStream chan string, processor messageProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	data := <-inputStream
	fmt.Println("dataSanitization:", data)
	sanatizedData := stringsanatization.Sanatize(data)
	processor.Push(sanatizedData)
}
