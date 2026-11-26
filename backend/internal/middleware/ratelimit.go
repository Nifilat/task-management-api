package middleware

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket allows bursts up to capacity and refills one token every refillEvery.
type TokenBucket struct {
	mu          sync.Mutex
	tokens      int
	capacity    int
	refillEvery time.Duration
	lastFilled  time.Time
	lastSeen    time.Time
}

func NewTokenBucket(capacity int, refillEvery time.Duration) *TokenBucket {
	now := time.Now()
	return &TokenBucket{tokens: capacity, capacity: capacity, refillEvery: refillEvery, lastFilled: now, lastSeen: now}
}

// Allow takes one token if there is one.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.lastSeen = now

	if b.tokens == b.capacity {
		b.lastFilled = now
	} else if credited := int(now.Sub(b.lastFilled) / b.refillEvery); credited > 0 {
		b.tokens = min(b.capacity, b.tokens+credited)
		b.lastFilled = b.lastFilled.Add(time.Duration(credited) * b.refillEvery)
	}

	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}

func (b *TokenBucket) idleSince() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.lastSeen
}

// RateLimit gives each client IP its own token bucket and answers 429 once it is empty.
func RateLimit(capacity int, refillEvery time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]*TokenBucket{}

	// Drop buckets of clients that have been quiet for a while so the map doesn't grow forever.
	go func() {
		for range time.Tick(time.Minute) {
			mu.Lock()
			for ip, b := range buckets {
				if time.Since(b.idleSince()) > 10*time.Minute {
					delete(buckets, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		bucket, ok := buckets[ip]
		if !ok {
			bucket = NewTokenBucket(capacity, refillEvery)
			buckets[ip] = bucket
		}
		mu.Unlock()

		if !bucket.Allow() {
			log.Printf("[RATE LIMIT] ip=%s method=%s path=%s", ip, c.Request.Method, c.Request.URL.Path)
			c.Header("Retry-After", strconv.Itoa(int(refillEvery.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
