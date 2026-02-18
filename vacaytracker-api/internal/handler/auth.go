package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/middleware"
	"vacaytracker-api/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles POST /api/auth/login
// Authenticates user and returns JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Attempt login
	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Login failed",
			})
		}
		return
	}

	// Return token and user
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	})
}

// Me handles GET /api/auth/me
// Returns the currently authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	// Get user from database
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to get user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// ChangePassword handles PUT /api/auth/password
// Changes the current user's password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	var req dto.ChangePasswordRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Change password
	err := h.authService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to change password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Password changed successfully",
	})
}

// UpdateEmailPreferences handles PUT /api/auth/email-preferences
// Updates the current user's email notification preferences
func (h *AuthHandler) UpdateEmailPreferences(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	var req dto.UpdateEmailPreferencesRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update preferences
	user, err := h.authService.UpdateEmailPreferences(c.Request.Context(), userID, &req)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to update email preferences",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emailPreferences": user.EmailPreferences,
	})
}
