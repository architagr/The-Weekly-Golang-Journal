package mutexusage

import (
	"sync"
	"testing"
)

func TestAddBalance(t *testing.T) {
	countOfConcurrentThreads := 100
	amount := 100
	wg := &sync.WaitGroup{}
	for i := 0; i < countOfConcurrentThreads; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup) {
			defer wg.Done()
			addBalance(amount)
		}(wg)
	}
	wg.Wait()
	expectedBalance := countOfConcurrentThreads * amount
	if bankBalance != expectedBalance {
		t.Errorf("expected bank balance to be %d, but got %d", expectedBalance, bankBalance)
	}
}
