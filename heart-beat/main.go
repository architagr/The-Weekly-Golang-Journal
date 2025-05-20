package main

import (
	"fmt"
	"time"
)

func dowork(done <-chan interface{}, puseInterval time.Duration) (<-chan interface{}, <-chan struct{}) {
	// create a channel to send heartbeats and a channel to send results
	heartbreater := make(chan interface{})
	result := make(chan struct{})
	go func() {
		defer close(result)
		defer close(heartbreater)
		// create a ticker to send heartbeats and a ticker to send results (just for example)
		pulse := time.NewTicker(puseInterval)
		workGen := time.NewTicker(3 * puseInterval)

		defer pulse.Stop()
		defer workGen.Stop()
		sendPulse := func() {
			// send a heartbeat if the channel is not full
			// this is a non-blocking send
			// if the channel is full, it will not block and will not send
			// this is just an example, in a real application you would
			// probably want to handle the case where the channel is full
			// and either drop the heartbeat or wait for it to be read
			// in this example we just drop the heartbeat
			// and do not send it
			select {
			case heartbreater <- struct{}{}:
			default:
			}
		}
		sendResult := func(res struct{}) {
			// send a result if the channel is not full
			// this is a blocking send
			// if the channel is full, it will block and will wait to send result
			// just like done we should handle the case to validate the pulse when sending or receiving.
			for {
				select {
				case <-done:
					return
				case <-pulse.C:
					sendPulse()
				case result <- res:
					return
				}
			}
		}
		for {
			select {
			case <-done:
				return
			case <-pulse.C:
				sendPulse()
			case <-workGen.C:
				sendResult(struct{}{})
			}
		}
	}()
	return heartbreater, result
}

func main() {
	done := make(chan interface{})

	// simulate a done channel that will be closed after 10 seconds
	// this is just for example, in a real application you would
	// probably want to use a context with a timeout or a wait group
	// to signal when the work is done
	time.AfterFunc(10*time.Second, func() { close(done) })

	pulseInterval := 1 * time.Second
	heartbreater, result := dowork(done, pulseInterval)

	// simulate work
	go func() {
		for {
			select {

			case _, ok := <-heartbreater:
				if !ok {
					fmt.Println("worker heartbeat stoped")
					return
				}
				fmt.Println("worker heartbeat")
			case _, ok := <-result:
				if !ok {
					return
				}
				fmt.Println("worker completed work")
				// handle result
				// in this example we just print it
				// in a real application you would probably want to do something
				// more useful with the result
				// like send it to a channel or process it in some way
			}
		}
	}()

	time.Sleep(20 * time.Second)
}
