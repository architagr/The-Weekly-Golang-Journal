package ordernumbergenerator

import (
	"sync"
	"testing"
)

func TestGenerateOrderNumber(t *testing.T) {
	countOfConcurrentThreads := 100
	gen := InitOrderNumberGenerator(0)

	wg := &sync.WaitGroup{}
	orderNumbers := sync.Map{}
	for i := 0; i < countOfConcurrentThreads; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			defer wg.Done()
			orderNumbers.Store(gen.GenerateOrderNumber(), true)
		}(wg)
	}
	wg.Wait()
	count := 0
	orderNumbers.Range(func(key, value any) bool {
		count++
		return true
	})
	if count != countOfConcurrentThreads {
		t.Errorf("some order number are re used")
	}
}
