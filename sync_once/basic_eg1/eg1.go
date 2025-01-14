package basiceg1

import (
	"fmt"
	"sync"
)

func initLogger() {
	fmt.Println("Logger initialized")
}

func Run() {
	// 1️⃣ Created this channel so that we can signal
	// multiple goroutines at same time
	ch := make(chan bool)

	// 2️⃣ initilized the object for sync.Once
	once := &sync.Once{}

	// 3️⃣ a wait group to wait till all goroutine have completed
	wg := &sync.WaitGroup{}

	// 4️⃣ have 10 goroutine that wants to run the initilize function
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int, c <-chan bool, waitgroup *sync.WaitGroup) {
			defer waitgroup.Done()
			<-c
			fmt.Println("goroutine index ", index, "unblocked")

			// 5️⃣ added a function to object of sync.Once using the Do function,
			// and this is the task for the sync.Once object
			once.Do(initLogger)
		}(i, ch, wg)
	}
	fmt.Println("added goroutine that will run the initilize code using sync.Once")

	// 6️⃣ closed the channel to signal all the goroutine to be unblocked
	close(ch)
	wg.Wait()
}
