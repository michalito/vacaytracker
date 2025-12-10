package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders returns a middleware that adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable XSS filter in browsers
		c.Header("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS (only send over HTTPS in production)
		// max-age=31536000 = 1 year
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy - restrictive defaults
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")

		// Prevent browser from sending referrer info to other sites
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Disable browser features that could be security risks
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Cache control for sensitive endpoints
		if c.Request.URL.Path == "/api/auth/login" || c.Request.URL.Path == "/api/auth/me" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}

// ProductionSecurityHeaders returns security headers optimized for production
// This includes stricter CSP and additional protections
func ProductionSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// All the standard headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Stricter CSP for production
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")

		// Cache control
		if c.Request.URL.Path == "/api/auth/login" || c.Request.URL.Path == "/api/auth/me" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
		}

		c.Next()
	}
}
