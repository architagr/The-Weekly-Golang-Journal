package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SlidingWindowCounter struct {
	mu           sync.Mutex
	capacity     int
	window       time.Duration
	subWindow    time.Duration
	subWindows   map[int64]int
	lastAccessed int64
}

func NewSlidingWindowCounter(capacity int, window time.Duration, precision int) *SlidingWindowCounter {
	subWindow := time.Duration(int64(window) / int64(precision))
	log.Printf("[rate-limiter] 🎯 Sliding Window Counter initialized with capacity=%d, window=%v, subWindow=%v", capacity, window, subWindow)

	return &SlidingWindowCounter{
		capacity:   capacity,
		window:     window,
		subWindow:  subWindow,
		subWindows: make(map[int64]int),
	}
}

func (s *SlidingWindowCounter) allow() (bool, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	currBucket := now.Unix() / int64(s.subWindow.Seconds())

	// Clean up expired buckets
	expiredCutoff := now.Add(-s.window).Unix() / int64(s.subWindow.Seconds())
	for ts := range s.subWindows {
		if ts <= expiredCutoff {
			delete(s.subWindows, ts)
		}
	}

	// Count total requests in window
	var total int
	for _, count := range s.subWindows {
		total += count
	}

	if total < s.capacity {
		s.subWindows[currBucket]++
		return true, time.Time{}
	}

	// Denied: Estimate next reset time
	var earliestBucket int64 = currBucket
	for ts := range s.subWindows {
		if ts < earliestBucket {
			earliestBucket = ts
		}
	}
	retryAt := time.Unix(earliestBucket*int64(s.subWindow.Seconds()), 0).Add(s.window)
	return false, retryAt
}

func SlidingWindowCounterMiddleware(capacity int, window time.Duration, precision int) gin.HandlerFunc {
	limiter := NewSlidingWindowCounter(capacity, window, precision)

	return func(c *gin.Context) {
		allowed, retryAt := limiter.allow()

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", capacity))
		// approximation
		c.Header("X-RateLimit-Remaining", "N/A")

		if allowed {
			log.Printf("[rate-limiter] ✅ Allowed (Counter) — %s", c.Request.URL.Path)
			c.Next()
			return
		}

		retryAfter := time.Until(retryAt).Seconds()
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", retryAt.Unix()))
		c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

		log.Printf("[rate-limiter] ⛔ Denied (Counter) — Path: %s | Retry after %.0f sec", c.Request.URL.Path, retryAfter)
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error":       "Too Many Requests",
			"retry_after": fmt.Sprintf("%.0f seconds", retryAfter),
		})
	}
}
