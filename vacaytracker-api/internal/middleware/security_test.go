package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── SecurityHeaders Tests ───

func TestSecurityHeaders_AllPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rec.Header().Get("X-XSS-Protection"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "includeSubDomains")
	assert.NotEmpty(t, rec.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	assert.Contains(t, rec.Header().Get("Permissions-Policy"), "geolocation=()")
	assert.Contains(t, rec.Header().Get("Permissions-Policy"), "microphone=()")
	assert.Contains(t, rec.Header().Get("Permissions-Policy"), "camera=()")
}

func TestSecurityHeaders_CacheControl_LoginPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeaders())
	router.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "no-store, no-cache, must-revalidate, private", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
}

func TestSecurityHeaders_CacheControl_MePath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/api/auth/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "no-store, no-cache, must-revalidate, private", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
}

func TestSecurityHeaders_CacheControl_NotSetForOtherRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/api/vacation/requests", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/vacation/requests", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get("Cache-Control"))
	assert.Empty(t, rec.Header().Get("Pragma"))
}

// ─── ProductionSecurityHeaders Tests ───

func TestProductionSecurityHeaders_StricterCSP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	csp := rec.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "frame-ancestors 'none'")
	assert.Contains(t, csp, "base-uri 'self'")
	assert.Contains(t, csp, "form-action 'self'")
}

func TestProductionSecurityHeaders_HSTSIncludesPreload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	hsts := rec.Header().Get("Strict-Transport-Security")
	assert.Contains(t, hsts, "preload")
	assert.Contains(t, hsts, "includeSubDomains")
	assert.Contains(t, hsts, "max-age=31536000")
}

func TestProductionSecurityHeaders_AllBaseHeadersPresent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", rec.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", rec.Header().Get("Referrer-Policy"))
	assert.Contains(t, rec.Header().Get("Permissions-Policy"), "geolocation=()")
}

func TestProductionSecurityHeaders_CacheControl_LoginPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, "no-store, no-cache, must-revalidate, private", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
}

func TestProductionSecurityHeaders_NoCacheControlForOtherRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.GET("/api/vacation/requests", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/vacation/requests", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Cache-Control"))
	assert.Empty(t, rec.Header().Get("Pragma"))
}

func TestSecurityHeaders_CSPIncludesExpectedDirectives(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	csp := rec.Header().Get("Content-Security-Policy")
	assert.Contains(t, csp, "default-src 'self'")
	assert.Contains(t, csp, "script-src 'self'")
	assert.Contains(t, csp, "connect-src 'self'")
}

func TestProductionSecurityHeaders_CSPDoesNotAllowUnsafeInline(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ProductionSecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	csp := rec.Header().Get("Content-Security-Policy")
	assert.NotContains(t, csp, "'unsafe-inline'")
}
