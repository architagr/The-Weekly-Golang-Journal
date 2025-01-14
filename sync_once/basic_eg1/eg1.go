package basiceg1

import (
	"fmt"
	"sync"
)

func initLogger() {
	fmt.Println("Logger initialized")
}

func Run() {
	// 1️⃣ Channel to signal all Goroutines simultaneously
	ch := make(chan bool)

	// 2️⃣ sync.Once object for single-use execution
	once := &sync.Once{}

	// 3️⃣ WaitGroup to wait for all Goroutines to finish
	wg := &sync.WaitGroup{}

	// 4️⃣ Launch 10 Goroutines waiting on the same signal
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int, c <-chan bool, waitgroup *sync.WaitGroup) {
			defer waitgroup.Done()
			// Unblocks when the channel is closed
			<-c
			fmt.Println("goroutine index ", index, "unblocked")

			// 5️⃣ sync.Once controlling the initialization
			once.Do(initLogger)
		}(i, ch, wg)
	}
	fmt.Println("added goroutine that will run the initilize code using sync.Once")

	// 6️⃣ Signal all Goroutines to proceed
	close(ch)
	wg.Wait()
}
