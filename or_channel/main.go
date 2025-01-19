package main

import (
	"fmt"
	"time"
)

func main() {
	var or func(doneChannels ...<-chan interface{}) <-chan interface{}
	or = func(doneChannels ...<-chan interface{}) <-chan interface{} {
		switch len(doneChannels) {
		case 0:
			return nil
		case 1:
			return doneChannels[0]
		}

		done := make(chan interface{})
		go func() {
			defer close(done)
			switch len(doneChannels) {
			case 2:
				select {
				case <-doneChannels[0]:
				case <-doneChannels[1]:
				}
			default:
				select {
				case <-doneChannels[0]:
				case <-doneChannels[1]:
				case <-doneChannels[2]:
				case <-or(append(doneChannels[3:], done)...):
				}
			}
		}()
		return done
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
