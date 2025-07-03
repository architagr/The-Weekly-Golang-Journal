package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type FixedWindowCounter struct {
	capacity  int
	counter   int
	window    time.Duration
	lastReset time.Time
	mu        sync.Mutex
}

func NewFixedWindowCounter(capacity int, window time.Duration) *FixedWindowCounter {
	bucket := &FixedWindowCounter{
		capacity:  capacity,
		counter:   0,
		window:    window,
		lastReset: time.Now().UTC(),
	}
	go bucket.startWindow()
	log.Printf("[rate-limiter] 🪟 Initialized with capacity=%d, window=%v", capacity, window)
	return bucket
}

// Periodically resets counter at window boundary
func (f *FixedWindowCounter) startWindow() {
	ticker := time.NewTicker(f.window)
	defer ticker.Stop()

	for range ticker.C {
		f.mu.Lock()
		f.lastReset = time.Now().UTC()
		f.counter = 0
		log.Printf("[rate-limiter] 🔄 Window reset at %v", f.lastReset)
		f.mu.Unlock()
	}
}

// Gin middleware
func FixedWindowCounterMiddleware(capacity int, window time.Duration) gin.HandlerFunc {
	limiter := NewFixedWindowCounter(capacity, window)

	return func(c *gin.Context) {
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		remaining := limiter.capacity - limiter.counter
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if limiter.counter < limiter.capacity {
			limiter.counter++
			log.Printf("[rate-limiter] ✅ Allowed - %s | Count: %d/%d", c.Request.URL.Path, limiter.counter, limiter.capacity)
			c.Next()
		} else {
			resetAt := limiter.lastReset.Add(limiter.window)
			retryAfter := time.Until(resetAt).Seconds()

			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt.Unix()))
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

			log.Printf("[rate-limiter] ⛔ Rate limit exceeded for %s | IP: %s", c.Request.URL.Path, c.ClientIP())
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too Many Requests",
				"retry_after": fmt.Sprintf("%.0f seconds", retryAfter),
			})
		}
	}
}
