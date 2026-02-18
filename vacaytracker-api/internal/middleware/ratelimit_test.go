package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── RateLimiter.Allow Tests ───

func TestRateLimiter_Allow_FirstRequestWithinLimit(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)

	assert.True(t, rl.Allow("192.168.1.1"))
}

func TestRateLimiter_Allow_RequestsUpToLimit(t *testing.T) {
	limit := 5
	rl := NewRateLimiter(limit, time.Minute)

	for i := 0; i < limit; i++ {
		assert.True(t, rl.Allow("192.168.1.1"), "request %d should be allowed", i+1)
	}
}

func TestRateLimiter_Allow_ExceedingLimit(t *testing.T) {
	limit := 3
	rl := NewRateLimiter(limit, time.Minute)

	for i := 0; i < limit; i++ {
		rl.Allow("192.168.1.1")
	}

	assert.False(t, rl.Allow("192.168.1.1"), "request beyond limit should be denied")
}

func TestRateLimiter_Allow_DifferentIPsAreIndependent(t *testing.T) {
	limit := 2
	rl := NewRateLimiter(limit, time.Minute)

	// Exhaust limit for IP1
	for i := 0; i < limit; i++ {
		rl.Allow("10.0.0.1")
	}
	assert.False(t, rl.Allow("10.0.0.1"), "IP1 should be blocked")

	// IP2 should still be allowed
	assert.True(t, rl.Allow("10.0.0.2"), "IP2 should be allowed independently")
}

func TestRateLimiter_Allow_AfterWindowExpires(t *testing.T) {
	// Use a very short window so the test doesn't take long
	rl := NewRateLimiter(1, 50*time.Millisecond)

	assert.True(t, rl.Allow("192.168.1.1"))
	assert.False(t, rl.Allow("192.168.1.1"))

	// Wait for the window to expire
	time.Sleep(80 * time.Millisecond)

	assert.True(t, rl.Allow("192.168.1.1"), "should be allowed after window expires")
}

// ─── RateLimiter.RemainingRequests Tests ───

func TestRateLimiter_RemainingRequests_FullLimitForNewIP(t *testing.T) {
	limit := 10
	rl := NewRateLimiter(limit, time.Minute)

	assert.Equal(t, limit, rl.RemainingRequests("10.0.0.99"))
}

func TestRateLimiter_RemainingRequests_DecrementedAfterRequests(t *testing.T) {
	limit := 10
	rl := NewRateLimiter(limit, time.Minute)

	rl.Allow("10.0.0.1")
	rl.Allow("10.0.0.1")
	rl.Allow("10.0.0.1")

	assert.Equal(t, limit-3, rl.RemainingRequests("10.0.0.1"))
}

func TestRateLimiter_RemainingRequests_ZeroWhenExhausted(t *testing.T) {
	limit := 2
	rl := NewRateLimiter(limit, time.Minute)

	for i := 0; i < limit; i++ {
		rl.Allow("10.0.0.1")
	}

	assert.Equal(t, 0, rl.RemainingRequests("10.0.0.1"))
}

// ─── RateLimiter.Middleware Tests ───

func TestRateLimiterMiddleware_AllowedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(10, time.Minute)

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-RateLimit-Remaining"))
	assert.Equal(t, "10", rec.Header().Get("X-RateLimit-Limit"))
}

func TestRateLimiterMiddleware_BlockedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(1, time.Minute)

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First request is allowed
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	rec1 := httptest.NewRecorder()
	router.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// Second request is blocked
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec2.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", body["code"])
}

// ─── Factory Tests ───

func TestLoginRateLimiter_Limit(t *testing.T) {
	rl := LoginRateLimiter()

	assert.Equal(t, 5, rl.limit)
}

func TestAPIRateLimiter_Limit(t *testing.T) {
	rl := APIRateLimiter()

	assert.Equal(t, 100, rl.limit)
}

func TestRateLimiterMiddleware_HeaderValues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := NewRateLimiter(5, time.Minute)

	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Make 3 requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.5:9999"
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Fourth request — check remaining header
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.5:9999"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "5", rec.Header().Get("X-RateLimit-Limit"))
	// After 4 Allow calls, remaining = 5 - 4 = 1
	assert.Equal(t, "1", rec.Header().Get("X-RateLimit-Remaining"))
}
