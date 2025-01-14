package basiceg1

import (
	"fmt"
	"sync"
)

func initilize() {
	fmt.Println("just initilized")
}

func Run() {
	ch := make(chan bool)
	once := &sync.Once{}
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int, c <-chan bool, waitgroup *sync.WaitGroup) {
			defer waitgroup.Done()
			<-c
			fmt.Println("goroutine index ", index, "unblocked")
			once.Do(initilize)
		}(i, ch, wg)
	}
	fmt.Println("added goroutine that will run the initilize code using sync.Once")
	close(ch)
	wg.Wait()
}
