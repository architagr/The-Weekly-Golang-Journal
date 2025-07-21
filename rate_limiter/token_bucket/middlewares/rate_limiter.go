package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	bucket         []int64
	rate           int
	bucketCapacity int
	nextRefillTime time.Time
	mu             sync.Mutex
)

func refillBucket() {
	d := time.Second * time.Duration(rate)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		mu.Lock()
		refilled := 0
		for i := 0; i < rate && len(bucket) < bucketCapacity; i++ {
			bucket = append(bucket, time.Now().Unix())
			refilled++
		}
		nextRefillTime = time.Now().UTC().Add(d)
		log.Printf("[refill] Added %d tokens. New bucket size: %d. Next refill at: %s\n", refilled, len(bucket), nextRefillTime.Format(time.RFC1123))
		mu.Unlock()
	}
}

func RateLimiterMiddleware(r, cap int) gin.HandlerFunc {
	bucketCapacity = cap
	rate = r
	bucket = make([]int64, bucketCapacity)
	for i := 0; i < bucketCapacity; i++ {
		bucket[i] = time.Now().Unix()
	}
	nextRefillTime = time.Now().UTC().Add(time.Second * (time.Duration(rate)))

	log.Printf("[init] Token bucket initialized with capacity=%d tokens and refillRate=%d second\n", bucketCapacity, rate)

	go refillBucket()

	return func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		log.Printf("[request] Incoming request %s %s", c.Request.Method, c.Request.URL.Path)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", bucketCapacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", len(bucket)))

		if len(bucket) > 0 {
			bucket = bucket[1:]
			log.Printf("[allow] Request allowed. Remaining tokens: %d", len(bucket))
			c.Next()
		} else {
			resetIn := time.Until(nextRefillTime).Seconds()
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", nextRefillTime.Unix()))
			c.Header("Retry-After", fmt.Sprintf("%.0f", resetIn))
			log.Printf("[deny] Too many requests. Retry after %.0f seconds (next refill at %s)", resetIn, nextRefillTime.Format(time.RFC1123))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
		}
	}
}
