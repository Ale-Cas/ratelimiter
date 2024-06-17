package src

import (
	"time"

	"github.com/gin-gonic/gin"
)

type TokenBucket struct {
	ipAddr string // client's IP address
	// tokens and capacity can only be positive (or zero) integers
	tokens           uint64
	capacity         uint64
	lastTokenAddedAt time.Time // last time a token was added to the bucket
}

// refill adds tokens to the bucket every second
// until it reaches its capacity.
func (tb *TokenBucket) refill() {
	if tb.tokens < tb.capacity {
		seconds := time.Since(tb.lastTokenAddedAt).Seconds()
		if seconds > 0 {
			tokens := tb.tokens + uint64(seconds)
			tb.tokens = min(tokens, tb.capacity)
			tb.lastTokenAddedAt = time.Now()
		}
	}
}

var buckets []TokenBucket

// TokenBucketLimiter is a middleware that acts as a rate limiter
// based on the token bucket algorithm.
// It takes an integer as an argument, which represents the number
// of requests per second that the server will allow.
func TokenBucketLimiter(reqPerSec uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the client's IP address is in the buckets slice
		for i := range buckets {
			// Add a new token every second
			buckets[i].refill()
			// Check if the client's IP address is in the buckets slice
			if buckets[i].ipAddr == c.ClientIP() {
				// If it is, decrement the number of tokens in the bucket
				if buckets[i].tokens > 0 {
					buckets[i].tokens--
					return
				} else {
					c.JSON(429, gin.H{
						"message": "Too Many Requests",
					})
					c.Abort()
					return
				}
			}
		}
		// If the client's IP address is not in the buckets slice,
		// create a new bucket for that client
		bucket := TokenBucket{
			ipAddr:           c.ClientIP(),
			tokens:           reqPerSec,
			capacity:         reqPerSec,
			lastTokenAddedAt: time.Now(),
		}
		buckets = append(buckets, bucket)
	}
}
