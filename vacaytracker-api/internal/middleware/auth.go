package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/service"
)

// Context keys for storing user information
const (
	ContextKeyUserID = "userID"
	ContextKeyEmail  = "email"
	ContextKeyName   = "name"
	ContextKeyRole   = "role"
	ContextKeyClaims = "claims"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondWithError(c, dto.ErrTokenMissingError())
			return
		}

		// Check for Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondWithError(c, dto.ErrTokenInvalidError())
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			if appErr, ok := err.(*dto.AppError); ok {
				respondWithError(c, appErr)
			} else {
				respondWithError(c, dto.ErrTokenInvalidError())
			}
			return
		}

		// Store user info in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyName, claims.Name)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// AdminMiddleware ensures the user has admin role
// Must be used after AuthMiddleware
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists {
			respondWithError(c, dto.ErrTokenMissingError())
			return
		}

		if role.(domain.Role) != domain.RoleAdmin {
			respondWithError(c, dto.ErrAdminRequiredError())
			return
		}

		c.Next()
	}
}

// EmployeeMiddleware ensures the user has employee role
// Must be used after AuthMiddleware
func EmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists {
			respondWithError(c, dto.ErrTokenMissingError())
			return
		}

		if role.(domain.Role) != domain.RoleEmployee {
			respondWithError(c, dto.ErrForbiddenError("This action is only available to employees"))
			return
		}

		c.Next()
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get(ContextKeyUserID)
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// GetUserEmail retrieves the user email from the context
func GetUserEmail(c *gin.Context) string {
	email, _ := c.Get(ContextKeyEmail)
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetUserRole retrieves the user role from the context
func GetUserRole(c *gin.Context) domain.Role {
	role, _ := c.Get(ContextKeyRole)
	if role == nil {
		return ""
	}
	return role.(domain.Role)
}

// GetClaims retrieves the full JWT claims from the context
func GetClaims(c *gin.Context) *service.JWTClaims {
	claims, _ := c.Get(ContextKeyClaims)
	if claims == nil {
		return nil
	}
	return claims.(*service.JWTClaims)
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *gin.Context) bool {
	return GetUserRole(c) == domain.RoleAdmin
}

// IsEmployee checks if the current user is an employee
func IsEmployee(c *gin.Context) bool {
	return GetUserRole(c) == domain.RoleEmployee
}

// respondWithError sends an error response and aborts the request
func respondWithError(c *gin.Context, err *dto.AppError) {
	c.AbortWithStatusJSON(err.HTTPStatus, err.ToResponse())
}
