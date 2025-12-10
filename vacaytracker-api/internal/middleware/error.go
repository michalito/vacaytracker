package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/dto"
)

// ErrorMiddleware handles panics and converts them to proper error responses
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// Return internal server error
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:    dto.ErrInternal,
					Message: "An internal error occurred",
				})
			}
		}()

		c.Next()

		// Handle any errors set during request processing
		if len(c.Errors) > 0 {
			// Get the last error
			err := c.Errors.Last()

			// Check if it's an AppError
			if appErr, ok := err.Err.(*dto.AppError); ok {
				c.JSON(appErr.HTTPStatus, appErr.ToResponse())
				return
			}

			// Log unknown errors
			log.Printf("Unhandled error: %v", err.Err)

			// Return generic error
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "An internal error occurred",
			})
		}
	}
}

// RequestLoggerMiddleware logs request information
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log request details (already handled by gin.Logger())
		// This middleware can be extended for custom logging
	}
}
