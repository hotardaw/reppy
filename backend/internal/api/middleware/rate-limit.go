/*
how this "token bucket" algo implementation works:

1. user gets a bucket that can hold tokens (represented by `tokens` & `capacity`)
2. tokens refill steadily over time (controlled by `rate`)
3. each API req requires 1 token
4. if no tokens available, user's bucket is empty and req is rejected
*/

package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-fitsync/backend/internal/api/response"
)

type RateLimiter struct {
	tokens     float64
	capacity   float64
	rate       float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) refill() {
	now := time.Now()
	timePassed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += timePassed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastRefill = now
}

func (rl *RateLimiter) allow() bool {
	rl.mu.Lock()         // lock for now
	defer rl.mu.Unlock() // make sure we unlock when done

	rl.refill() // add new tokens based on time elapsed
	if rl.tokens >= 1 {
		rl.tokens--
		return true // allow request
	}
	return false // deny request
}

// stores rate limiters for different clients in map
var limiters = struct {
	sync.RWMutex
	m map[string]*RateLimiter // key is IP address
}{m: make(map[string]*RateLimiter)} // immediate initialization

func RateLimitMiddleware(requestsPerSecond float64, burstSize float64) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// using IP address as client identifier,
			clientIP := r.RemoteAddr

			// get/create rate limiter for this client
			limiters.Lock()
			limiter, exists := limiters.m[clientIP]
			if !exists {
				limiter = NewRateLimiter(requestsPerSecond, burstSize)
				limiters.m[clientIP] = limiter
			}
			limiters.Unlock()

			if !limiter.allow() {
				response.SendError(w, "Rate limit exceeded", http.StatusTooManyRequests) // 429
				return
			}

			next(w, r)
		}
	}
}
