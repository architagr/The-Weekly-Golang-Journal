package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always call cancel to release resources

	go func() {
		// Simulate a long-running task
		time.Sleep(10 * time.Second)
		if ctx.Err() != nil {
			return
		}
		fmt.Println("Task completed")
	}()

	// Simulate a long-running task
	select {
	case <-ctx.Done():
		fmt.Println("Task cancelled:", ctx.Err())
	}
}
