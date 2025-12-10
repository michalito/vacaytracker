package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityEvent represents a security-relevant event
type SecurityEvent struct {
	Timestamp   string `json:"timestamp"`
	EventType   string `json:"eventType"`
	IP          string `json:"ip"`
	UserAgent   string `json:"userAgent"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	StatusCode  int    `json:"statusCode"`
	UserID      string `json:"userId,omitempty"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
}

// SecurityLogger provides security event logging
type SecurityLogger struct {
	// Could be extended to write to file, database, or external service
}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{}
}

// LogEvent logs a security event
func (sl *SecurityLogger) LogEvent(event SecurityEvent) {
	log.Printf("[SECURITY] %s | Type: %s | IP: %s | Path: %s | Method: %s | Status: %d | UserID: %s | Email: %s | %s",
		event.Timestamp,
		event.EventType,
		event.IP,
		event.Path,
		event.Method,
		event.StatusCode,
		event.UserID,
		event.Email,
		event.Description,
	)
}

// LogLoginAttempt logs a login attempt (successful or failed)
func (sl *SecurityLogger) LogLoginAttempt(c *gin.Context, email string, success bool) {
	eventType := "LOGIN_SUCCESS"
	description := "User logged in successfully"
	if !success {
		eventType = "LOGIN_FAILED"
		description = "Failed login attempt"
	}

	sl.LogEvent(SecurityEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		EventType:   eventType,
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		StatusCode:  c.Writer.Status(),
		Email:       email,
		Description: description,
	})
}

// LogAdminAction logs an administrative action
func (sl *SecurityLogger) LogAdminAction(c *gin.Context, userID, action string) {
	sl.LogEvent(SecurityEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		EventType:   "ADMIN_ACTION",
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		StatusCode:  c.Writer.Status(),
		UserID:      userID,
		Description: action,
	})
}

// LogUnauthorizedAccess logs an unauthorized access attempt
func (sl *SecurityLogger) LogUnauthorizedAccess(c *gin.Context, reason string) {
	sl.LogEvent(SecurityEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		EventType:   "UNAUTHORIZED_ACCESS",
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		StatusCode:  c.Writer.Status(),
		Description: reason,
	})
}

// LogRateLimitExceeded logs when rate limit is exceeded
func (sl *SecurityLogger) LogRateLimitExceeded(c *gin.Context) {
	sl.LogEvent(SecurityEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		EventType:   "RATE_LIMIT_EXCEEDED",
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		StatusCode:  429,
		Description: "Rate limit exceeded",
	})
}

// LogSuspiciousActivity logs potentially suspicious activity
func (sl *SecurityLogger) LogSuspiciousActivity(c *gin.Context, description string) {
	sl.LogEvent(SecurityEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		EventType:   "SUSPICIOUS_ACTIVITY",
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Path:        c.Request.URL.Path,
		Method:      c.Request.Method,
		StatusCode:  c.Writer.Status(),
		Description: description,
	})
}

// SecurityLoggingMiddleware returns middleware that logs security events for certain routes
func SecurityLoggingMiddleware(logger *SecurityLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request first
		c.Next()

		// Log security-relevant responses
		status := c.Writer.Status()
		path := c.Request.URL.Path

		// Log failed authentication attempts
		if path == "/api/auth/login" && status == 401 {
			email, _ := c.Get("attempted_email")
			if emailStr, ok := email.(string); ok {
				logger.LogLoginAttempt(c, emailStr, false)
			} else {
				logger.LogLoginAttempt(c, "", false)
			}
		}

		// Log unauthorized access attempts
		if status == 401 || status == 403 {
			logger.LogUnauthorizedAccess(c, "Access denied")
		}

		// Log admin actions
		if len(path) > 11 && path[:11] == "/api/admin/" && (c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE") {
			userID, _ := c.Get("userID")
			if userIDStr, ok := userID.(string); ok {
				action := c.Request.Method + " " + path
				logger.LogAdminAction(c, userIDStr, action)
			}
		}
	}
}
