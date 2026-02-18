package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/dto"
)

// RateLimiter provides IP-based rate limiting
type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*rateLimitEntry
	limit    int
	window   time.Duration
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*rateLimitEntry),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine to remove expired entries
	go rl.cleanup()

	return rl
}

// cleanup periodically removes expired entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, entry := range rl.requests {
			if now.After(entry.resetTime) {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if the request should be allowed based on rate limiting
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.requests[ip]

	if !exists || now.After(entry.resetTime) {
		// Create new entry or reset expired entry
		rl.requests[ip] = &rateLimitEntry{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if entry.count >= rl.limit {
		return false
	}

	entry.count++
	return true
}

// RemainingRequests returns how many requests are remaining for an IP
func (rl *RateLimiter) RemainingRequests(ip string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	entry, exists := rl.requests[ip]
	if !exists || time.Now().After(entry.resetTime) {
		return rl.limit
	}

	remaining := rl.limit - entry.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Middleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Code:    dto.ErrRateLimitExceeded,
				Message: "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(rl.RemainingRequests(ip)))

		c.Next()
	}
}

// LoginRateLimiter creates a rate limiter for login attempts (5 per minute)
func LoginRateLimiter() *RateLimiter {
	return NewRateLimiter(5, time.Minute)
}

// APIRateLimiter creates a rate limiter for general API requests (100 per minute)
func APIRateLimiter() *RateLimiter {
	return NewRateLimiter(100, time.Minute)
}
