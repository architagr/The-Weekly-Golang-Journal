package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LeakyBucket struct {
	capacity   int
	queue      chan struct{}
	leakRate   time.Duration
	lastLeaked time.Time
}

func NewLeakyBucket(capacity int, leakRatePerSec int) *LeakyBucket {
	bucket := &LeakyBucket{
		capacity: capacity,
		queue:    make(chan struct{}, capacity),
		leakRate: time.Second / time.Duration(leakRatePerSec),
	}
	go bucket.startLeaking()
	log.Printf("[init] Leaky bucket created — capacity=%d, leak rate=%v/sec", capacity, bucket.leakRate)
	return bucket
}

func (lb *LeakyBucket) startLeaking() {
	ticker := time.NewTicker(lb.leakRate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-lb.queue:
			lb.lastLeaked = time.Now().UTC()
			log.Printf("[leak] Processed 1 request. Remaining in queue: %d", len(lb.queue))
		default:
			// Nothing to leak right now
		}
	}
}

func LeakyBucketMiddleware(capacity int, leakRatePerSec int) gin.HandlerFunc {
	bucket := NewLeakyBucket(capacity, leakRatePerSec)

	return func(c *gin.Context) {
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", bucket.capacity))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", bucket.capacity-len(bucket.queue)))

		select {
		case bucket.queue <- struct{}{}:
			log.Printf("[allow] Request enqueued. Current queue size: %d", len(bucket.queue))
			c.Next()
		default:
			nextRefillTime := bucket.lastLeaked.Add(bucket.leakRate)
			resetIn := time.Until(nextRefillTime).Seconds()
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", nextRefillTime.Unix()))
			c.Header("Retry-After", fmt.Sprintf("%.0f", resetIn))
			log.Printf("[deny] Bucket full. Dropping request to %s", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests — bucket full",
			})
		}
	}
}
