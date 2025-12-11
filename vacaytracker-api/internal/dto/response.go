package dto

import (
	"vacaytracker-api/internal/domain"
)

// ============================================
// Authentication Responses
// ============================================

// LoginResponse represents the login response
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID               string                  `json:"id"`
	Email            string                  `json:"email"`
	Name             string                  `json:"name"`
	Role             string                  `json:"role"`
	VacationBalance  int                     `json:"vacationBalance"`
	StartDate        *string                 `json:"startDate,omitempty"`
	EmailPreferences domain.EmailPreferences `json:"emailPreferences"`
	CreatedAt        string                  `json:"createdAt"`
	UpdatedAt        string                  `json:"updatedAt"`
}

// ToUserResponse converts a domain User to UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		Name:             user.Name,
		Role:             string(user.Role),
		VacationBalance:  user.VacationBalance,
		StartDate:        user.StartDate,
		EmailPreferences: user.EmailPreferences,
		CreatedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ============================================
// User List Response
// ============================================

// UserListResponse represents a paginated list of users
type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Pagination *PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ============================================
// Vacation Responses
// ============================================

// VacationRequestResponse represents a vacation request in API responses
type VacationRequestResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userId"`
	UserName        string  `json:"userName,omitempty"`
	UserEmail       string  `json:"userEmail,omitempty"`
	StartDate       string  `json:"startDate"`
	EndDate         string  `json:"endDate"`
	TotalDays       int     `json:"totalDays"`
	Reason          *string `json:"reason,omitempty"`
	Status          string  `json:"status"`
	ReviewedBy      *string `json:"reviewedBy,omitempty"`
	ReviewedAt      *string `json:"reviewedAt,omitempty"`
	RejectionReason *string `json:"rejectionReason,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// ToVacationRequestResponse converts a domain VacationRequest to response
func ToVacationRequestResponse(req *domain.VacationRequest) *VacationRequestResponse {
	resp := &VacationRequestResponse{
		ID:              req.ID,
		UserID:          req.UserID,
		UserName:        req.UserName,
		UserEmail:       req.UserEmail,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		TotalDays:       req.TotalDays,
		Reason:          req.Reason,
		Status:          string(req.Status),
		ReviewedBy:      req.ReviewedBy,
		RejectionReason: req.RejectionReason,
		CreatedAt:       req.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       req.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if req.ReviewedAt != nil {
		formatted := req.ReviewedAt.Format("2006-01-02T15:04:05Z")
		resp.ReviewedAt = &formatted
	}

	return resp
}

// VacationListResponse represents a list of vacation requests
type VacationListResponse struct {
	Requests []*VacationRequestResponse `json:"requests"`
	Total    int                        `json:"total"`
}

// TeamVacationResponse represents team vacation data for calendar
type TeamVacationResponse struct {
	Vacations []*TeamVacationItem `json:"vacations"`
	Month     int                 `json:"month"`
	Year      int                 `json:"year"`
}

// TeamVacationItem represents a single team vacation entry
type TeamVacationItem struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	TotalDays int    `json:"totalDays"`
}

// ============================================
// Settings Response
// ============================================

// SettingsResponse represents application settings
type SettingsResponse struct {
	ID                  string                   `json:"id"`
	WeekendPolicy       domain.WeekendPolicy     `json:"weekendPolicy"`
	Newsletter          domain.NewsletterConfig  `json:"newsletter"`
	DefaultVacationDays int                      `json:"defaultVacationDays"`
	VacationResetMonth  int                      `json:"vacationResetMonth"`
	UpdatedAt           string                   `json:"updatedAt"`
}

// ToSettingsResponse converts domain Settings to response
func ToSettingsResponse(settings *domain.Settings) *SettingsResponse {
	return &SettingsResponse{
		ID:                  settings.ID,
		WeekendPolicy:       settings.WeekendPolicy,
		Newsletter:          settings.Newsletter,
		DefaultVacationDays: settings.DefaultVacationDays,
		VacationResetMonth:  settings.VacationResetMonth,
		UpdatedAt:           settings.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ============================================
// Newsletter Responses
// ============================================

// NewsletterSendResponse represents the result of sending a newsletter
type NewsletterSendResponse struct {
	Success        bool   `json:"success"`
	RecipientCount int    `json:"recipientCount"`
	Message        string `json:"message"`
}

// NewsletterPreviewResponse represents a preview of the newsletter
type NewsletterPreviewResponse struct {
	Subject        string   `json:"subject"`
	HTMLBody       string   `json:"htmlBody"`
	TextBody       string   `json:"textBody"`
	Recipients     []string `json:"recipients"`
	RecipientCount int      `json:"recipientCount"`
}

// TestEmailResponse represents the result of sending a test email
type TestEmailResponse struct {
	Success  bool   `json:"success"`
	Template string `json:"template"`
	SentTo   string `json:"sentTo"`
	Message  string `json:"message"`
}

// EmailPreviewResponse represents a preview of an email template
type EmailPreviewResponse struct {
	Template string `json:"template"`
	Subject  string `json:"subject"`
	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`
}

// ============================================
// Balance Reset Response
// ============================================

// ResetBalancesResponse represents the result of resetting vacation balances
type ResetBalancesResponse struct {
	Success      bool   `json:"success"`
	UsersUpdated int    `json:"usersUpdated"`
	NewBalance   int    `json:"newBalance"`
	Message      string `json:"message"`
}

// ============================================
// Generic Responses
// ============================================

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// SuccessResponse represents a success response with optional data
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
