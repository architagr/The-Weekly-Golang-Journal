package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SlidingWindowLog struct {
	mu       sync.Mutex
	logs     []int64
	capacity int
	window   time.Duration
}

func NewSlidingWindowLog(capacity int, window time.Duration) *SlidingWindowLog {
	log.Printf("[rate-limiter] 🪵 Sliding Window Log initialized with capacity=%d and window=%v", capacity, window)
	return &SlidingWindowLog{
		logs:     make([]int64, 0),
		capacity: capacity,
		window:   window,
	}
}

func (s *SlidingWindowLog) allow() (bool, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Unix()
	s.logs = append(s.logs, now)

	cutoff := now - int64(s.window.Seconds())
	// Clean up old logs
	newLogs := make([]int64, 0, len(s.logs))
	for _, ts := range s.logs {
		if ts > cutoff {
			newLogs = append(newLogs, ts)
		}
	}
	s.logs = newLogs

	if len(s.logs) < s.capacity {
		return true, time.Time{}
	}

	// Denied — find when the earliest timestamp will expire
	earliest := time.Unix(s.logs[0], 0).Add(s.window)
	return false, earliest
}

func SlidingWindowLogMiddleware(capacity int, window time.Duration) gin.HandlerFunc {
	limiter := NewSlidingWindowLog(capacity, window)

	return func(c *gin.Context) {
		allowed, retryAt := limiter.allow()

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", capacity-len(limiter.logs)))

		if allowed {
			log.Printf("[rate-limiter] ✅ Allowed (Log) — %s", c.Request.URL.Path)
			c.Next()
			return
		}

		retryAfter := time.Until(retryAt).Seconds()
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", retryAt.Unix()))
		c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

		log.Printf("[rate-limiter] ⛔ Denied (Log) — Path: %s | Retry after %.0f sec", c.Request.URL.Path, retryAfter)
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error":       "Too Many Requests",
			"retry_after": fmt.Sprintf("%.0f seconds", retryAfter),
		})
	}
}
