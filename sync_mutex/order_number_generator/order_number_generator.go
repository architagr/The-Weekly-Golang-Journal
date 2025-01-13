package ordernumbergenerator

import (
	"fmt"
	"sync"
	"time"
)

type orderNumberGenerator struct {
	lastSequence int
	// declare a lock that can be use to prevent race conditions
	lock *sync.Mutex
}

func InitOrderNumberGenerator(start int) *orderNumberGenerator {
	return &orderNumberGenerator{
		lastSequence: start,
		lock:         &sync.Mutex{},
	}
}

func (gen *orderNumberGenerator) GenerateOrderNumber() string {
	// Acquire lock by a goroutine
	gen.lock.Lock()
	// Release lock by a goroutine
	defer gen.lock.Unlock()
	gen.lastSequence++
	now := time.Now()
	return fmt.Sprintf("or-%d-ind-%d", now.Year(), gen.lastSequence)
}
