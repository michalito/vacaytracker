package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/middleware"
	"vacaytracker-api/internal/repository/sqlite"
	"vacaytracker-api/internal/service"
)

// AdminHandler handles admin management endpoints
type AdminHandler struct {
	userService       *service.UserService
	userRepo          *sqlite.UserRepository
	vacationService   *service.VacationService
	vacationRepo      *sqlite.VacationRepository
	settingsRepo      *sqlite.SettingsRepository
	emailService      *service.EmailService
	newsletterService *service.NewsletterService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(
	userService *service.UserService,
	userRepo *sqlite.UserRepository,
	vacationService *service.VacationService,
	vacationRepo *sqlite.VacationRepository,
	settingsRepo *sqlite.SettingsRepository,
	emailService *service.EmailService,
	newsletterService *service.NewsletterService,
) *AdminHandler {
	return &AdminHandler{
		userService:       userService,
		userRepo:          userRepo,
		vacationService:   vacationService,
		vacationRepo:      vacationRepo,
		settingsRepo:      settingsRepo,
		emailService:      emailService,
		newsletterService: newsletterService,
	}
}

// ============================================
// User Management Endpoints
// ============================================

// ListUsers handles GET /api/admin/users
// Lists all users with optional filtering and pagination
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	var role *domain.Role
	if r := c.Query("role"); r != "" {
		if r != string(domain.RoleAdmin) && r != string(domain.RoleEmployee) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrValidation,
				Message: "Invalid role. Must be admin or employee",
			})
			return
		}
		roleVal := domain.Role(r)
		role = &roleVal
	}

	search := c.Query("search")

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	users, total, err := h.userService.List(c.Request.Context(), role, search, page, limit)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to list users",
			})
		}
		return
	}

	// Convert to response DTOs
	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = dto.ToUserResponse(user)
	}

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, dto.UserListResponse{
		Users: responses,
		Pagination: &dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// CreateUser handles POST /api/admin/users
// Creates a new user
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Store the password before creating user (it gets hashed)
	tempPassword := req.Password

	user, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to create user",
			})
		}
		return
	}

	// Send welcome email with temporary password (non-blocking)
	h.emailService.SendWelcome(user, tempPassword)

	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

// GetUser handles GET /api/admin/users/:id
// Gets a user by ID
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetByID(c.Request.Context(), userID)
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

// UpdateUser handles PUT /api/admin/users/:id
// Updates a user
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, req, currentUserID)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to update user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// DeleteUser handles DELETE /api/admin/users/:id
// Deletes a user
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)

	err := h.userService.Delete(c.Request.Context(), userID, currentUserID)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to delete user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "User deleted successfully",
	})
}

// UpdateBalance handles PUT /api/admin/users/:id/balance
// Updates a user's vacation balance
func (h *AdminHandler) UpdateBalance(c *gin.Context) {
	userID := c.Param("id")

	var req dto.UpdateVacationBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateBalance(c.Request.Context(), userID, req.VacationBalance)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to update balance",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// ============================================
// Vacation Management Endpoints
// ============================================

// ListPending handles GET /api/admin/vacation/pending
// Lists all pending vacation requests
func (h *AdminHandler) ListPending(c *gin.Context) {
	requests, err := h.vacationService.ListPending(c.Request.Context())
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to list pending requests",
			})
		}
		return
	}

	// Convert to response DTOs
	responses := make([]*dto.VacationRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = dto.ToVacationRequestResponse(req)
	}

	c.JSON(http.StatusOK, dto.VacationListResponse{
		Requests: responses,
		Total:    len(responses),
	})
}

// Review handles PUT /api/admin/vacation/:id/review
// Approves or rejects a vacation request
func (h *AdminHandler) Review(c *gin.Context) {
	requestID := c.Param("id")
	adminID := middleware.GetUserID(c)

	var req dto.ReviewVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	var vacation *domain.VacationRequest
	var err error

	switch domain.VacationStatus(req.Status) {
	case domain.StatusApproved:
		vacation, err = h.vacationService.Approve(c.Request.Context(), requestID, adminID)
	case domain.StatusRejected:
		var reason *string
		if req.Reason != "" {
			reason = &req.Reason
		}
		vacation, err = h.vacationService.Reject(c.Request.Context(), requestID, adminID, reason)
	default:
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Status must be 'approved' or 'rejected'",
		})
		return
	}

	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to review request",
			})
		}
		return
	}

	// Send email notification to the user (non-blocking)
	go h.sendReviewEmail(c.Request.Context(), vacation, req.Status, req.Reason)

	c.JSON(http.StatusOK, dto.ToVacationRequestResponse(vacation))
}

// sendReviewEmail sends an email after a vacation request is reviewed
func (h *AdminHandler) sendReviewEmail(ctx context.Context, vacation *domain.VacationRequest, status string, reason string) {
	user, err := h.userRepo.GetByID(ctx, vacation.UserID)
	if err != nil || user == nil {
		return
	}

	switch domain.VacationStatus(status) {
	case domain.StatusApproved:
		h.emailService.SendRequestApproved(user, vacation)
	case domain.StatusRejected:
		h.emailService.SendRequestRejected(user, vacation, reason)
	}
}

// ============================================
// Balance Reset Endpoint
// ============================================

// ResetBalances handles POST /api/admin/users/reset-balances
// Resets all employee vacation balances to the default value from settings
func (h *AdminHandler) ResetBalances(c *gin.Context) {
	// Get settings to determine default vacation days
	settings, err := h.settingsRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to get settings",
		})
		return
	}

	// Reset all balances
	count, err := h.userService.ResetAllBalances(c.Request.Context(), settings.DefaultVacationDays)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to reset balances",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ResetBalancesResponse{
		Success:      true,
		UsersUpdated: count,
		NewBalance:   settings.DefaultVacationDays,
		Message:      fmt.Sprintf("Reset vacation balance to %d days for %d employees", settings.DefaultVacationDays, count),
	})
}

// ============================================
// Settings Endpoints
// ============================================

// GetSettings handles GET /api/admin/settings
// Gets application settings
func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to get settings",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToSettingsResponse(settings))
}

// UpdateSettings handles PUT /api/admin/settings
// Updates application settings
func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Get current settings
	settings, err := h.settingsRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to get settings",
		})
		return
	}

	// Apply updates
	if req.WeekendPolicy != nil {
		if req.WeekendPolicy.ExcludeWeekends != nil {
			settings.WeekendPolicy.ExcludeWeekends = *req.WeekendPolicy.ExcludeWeekends
		}
		if req.WeekendPolicy.ExcludedDays != nil {
			settings.WeekendPolicy.ExcludedDays = *req.WeekendPolicy.ExcludedDays
		}
	}

	if req.Newsletter != nil {
		if req.Newsletter.Enabled != nil {
			settings.Newsletter.Enabled = *req.Newsletter.Enabled
		}
		if req.Newsletter.Frequency != nil {
			settings.Newsletter.Frequency = *req.Newsletter.Frequency
		}
		if req.Newsletter.DayOfMonth != nil {
			settings.Newsletter.DayOfMonth = *req.Newsletter.DayOfMonth
		}
	}

	if req.DefaultVacationDays != nil {
		settings.DefaultVacationDays = *req.DefaultVacationDays
	}

	if req.VacationResetMonth != nil {
		settings.VacationResetMonth = *req.VacationResetMonth
	}

	// Save settings
	if err := h.settingsRepo.Update(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to update settings",
		})
		return
	}

	// Fetch updated settings
	settings, _ = h.settingsRepo.Get(c.Request.Context())

	c.JSON(http.StatusOK, dto.ToSettingsResponse(settings))
}

// ============================================
// Newsletter Endpoints
// ============================================

// SendNewsletter handles POST /api/admin/newsletter/send
// Manually triggers newsletter sending
func (h *AdminHandler) SendNewsletter(c *gin.Context) {
	count, err := h.newsletterService.Send(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to send newsletter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.NewsletterSendResponse{
		Success:        true,
		RecipientCount: count,
		Message:        fmt.Sprintf("Newsletter sent to %d recipients", count),
	})
}

// PreviewNewsletter handles GET /api/admin/newsletter/preview
// Returns newsletter preview without sending
func (h *AdminHandler) PreviewNewsletter(c *gin.Context) {
	preview, err := h.newsletterService.GeneratePreview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    dto.ErrInternal,
			Message: "Failed to generate preview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, preview)
}
