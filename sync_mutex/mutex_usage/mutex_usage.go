package mutexusage

import "sync"

var (
	// defined a mutex
	mu          *sync.Mutex
	bankBalance int
)

func init() {
	// initilized a mutex
	mu = new(sync.Mutex)
}

func addBalance(amount int) {
	// a goroutine accuires a lock
	mu.Lock()
	bankBalance += amount
	// goroutine release the lock
	// so that another goroutine can access this
	mu.Unlock()
}
