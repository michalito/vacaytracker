package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/middleware"
	"vacaytracker-api/internal/repository"
	"vacaytracker-api/internal/service"
)

// VacationHandler handles vacation request endpoints
type VacationHandler struct {
	vacationService *service.VacationService
	vacationRepo    repository.VacationRepository
	userRepo        repository.UserRepository
	emailService    *service.EmailService
}

// NewVacationHandler creates a new VacationHandler
func NewVacationHandler(
	vacationService *service.VacationService,
	vacationRepo repository.VacationRepository,
	userRepo repository.UserRepository,
	emailService *service.EmailService,
) *VacationHandler {
	return &VacationHandler{
		vacationService: vacationService,
		vacationRepo:    vacationRepo,
		userRepo:        userRepo,
		emailService:    emailService,
	}
}

// Create handles POST /api/vacation/request
// Creates a new vacation request
func (h *VacationHandler) Create(c *gin.Context) {
	var req dto.CreateVacationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    dto.ErrValidation,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	vacation, err := h.vacationService.Create(c.Request.Context(), userID, req)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to create vacation request",
			})
		}
		return
	}

	// Send email notifications (non-blocking)
	// Use background context since the request context is cancelled after the response is sent
	go h.sendVacationRequestEmails(context.Background(), userID, vacation)

	c.JSON(http.StatusCreated, dto.ToVacationRequestResponse(vacation))
}

// sendVacationRequestEmails sends emails when a vacation request is created
func (h *VacationHandler) sendVacationRequestEmails(ctx context.Context, userID string, vacation *domain.VacationRequest) {
	// Get the user who submitted the request
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: failed to get user for email notification: %v", err)
		return
	}
	if user == nil {
		return
	}

	// Send confirmation email to the user
	h.emailService.SendRequestSubmitted(user, vacation)

	// Send notification to all admins
	admins, err := h.userRepo.GetByRole(ctx, domain.RoleAdmin)
	if err != nil {
		log.Printf("ERROR: failed to get admins for email notification: %v", err)
		return
	}
	if len(admins) == 0 {
		return
	}

	h.emailService.SendAdminNewRequest(admins, user, vacation)
}

// List handles GET /api/vacation/requests
// Lists vacation requests for the current user
func (h *VacationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	// Parse query parameters
	var status *domain.VacationStatus
	if s := c.Query("status"); s != "" {
		vs := domain.VacationStatus(s)
		if vs != domain.StatusPending && vs != domain.StatusApproved && vs != domain.StatusRejected {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrValidation,
				Message: "Invalid status. Must be pending, approved, or rejected",
			})
			return
		}
		status = &vs
	}

	var year *int
	if y := c.Query("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil || parsed < 2000 || parsed > 2100 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrValidation,
				Message: "Invalid year",
			})
			return
		}
		year = &parsed
	}

	requests, err := h.vacationService.ListByUser(c.Request.Context(), userID, status, year)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to list vacation requests",
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

// Get handles GET /api/vacation/requests/:id
// Gets a single vacation request
func (h *VacationHandler) Get(c *gin.Context) {
	requestID := c.Param("id")
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	request, err := h.vacationService.GetByID(c.Request.Context(), requestID)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to get vacation request",
			})
		}
		return
	}

	// Check if user has access (own request or admin)
	if request.UserID != userID && userRole != domain.RoleAdmin {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Code:    dto.ErrForbidden,
			Message: "You can only view your own requests",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToVacationRequestResponse(request))
}

// Cancel handles DELETE /api/vacation/requests/:id
// Cancels a pending vacation request
func (h *VacationHandler) Cancel(c *gin.Context) {
	requestID := c.Param("id")
	userID := middleware.GetUserID(c)

	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	err := h.vacationService.Cancel(c.Request.Context(), requestID, userID)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to cancel vacation request",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Vacation request cancelled successfully",
	})
}

// Team handles GET /api/vacation/team
// Gets team vacation calendar for a given month/year
func (h *VacationHandler) Team(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    dto.ErrAuthTokenMissing,
			Message: "Authentication required",
		})
		return
	}

	// Parse query parameters (default to current month)
	now := time.Now()
	month := now.Month()
	year := now.Year()

	if m := c.Query("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil || parsed < 1 || parsed > 12 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrValidation,
				Message: "Invalid month. Must be 1-12",
			})
			return
		}
		month = time.Month(parsed)
	}

	if y := c.Query("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil || parsed < 2000 || parsed > 2100 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    dto.ErrValidation,
				Message: "Invalid year",
			})
			return
		}
		year = parsed
	}

	vacations, err := h.vacationService.ListTeam(c.Request.Context(), int(month), year)
	if err != nil {
		if appErr, ok := err.(*dto.AppError); ok {
			c.JSON(appErr.HTTPStatus, appErr.ToResponse())
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    dto.ErrInternal,
				Message: "Failed to get team vacations",
			})
		}
		return
	}

	// Convert to response DTOs
	items := make([]*dto.TeamVacationItem, len(vacations))
	for i, v := range vacations {
		items[i] = &dto.TeamVacationItem{
			ID:        v.ID,
			UserID:    v.UserID,
			UserName:  v.UserName,
			StartDate: v.StartDate,
			EndDate:   v.EndDate,
			TotalDays: v.TotalDays,
		}
	}

	c.JSON(http.StatusOK, dto.TeamVacationResponse{
		Vacations: items,
		Month:     int(month),
		Year:      year,
	})
}
