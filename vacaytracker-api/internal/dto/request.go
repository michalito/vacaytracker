package dto

// ============================================
// Authentication Requests
// ============================================

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// ChangePasswordRequest represents the password change request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=72"`
}

// UpdateEmailPreferencesRequest represents the email preferences update request
type UpdateEmailPreferencesRequest struct {
	VacationUpdates   *bool `json:"vacationUpdates"`
	WeeklyDigest      *bool `json:"weeklyDigest"`
	TeamNotifications *bool `json:"teamNotifications"`
}

// ============================================
// User Management Requests (Admin)
// ============================================

// CreateUserRequest represents the user creation request body
type CreateUserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=72"`
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Role            string `json:"role" binding:"required,oneof=admin employee"`
	VacationBalance *int   `json:"vacationBalance"`
	StartDate       string `json:"startDate,omitempty"`
}

// UpdateUserRequest represents the user update request body
type UpdateUserRequest struct {
	Email           string `json:"email,omitempty" binding:"omitempty,email"`
	Name            string `json:"name,omitempty" binding:"omitempty,max=100"`
	Role            string `json:"role,omitempty" binding:"omitempty,oneof=admin employee"`
	VacationBalance *int   `json:"vacationBalance,omitempty"`
	StartDate       string `json:"startDate,omitempty"`
}

// UpdateVacationBalanceRequest represents the balance update request
type UpdateVacationBalanceRequest struct {
	VacationBalance int `json:"vacationBalance" binding:"required,min=0"`
}

// ============================================
// Vacation Requests
// ============================================

// CreateVacationRequest represents the vacation request creation body
// Dates should be in DD/MM/YYYY format (EU format)
type CreateVacationRequest struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Reason    string `json:"reason,omitempty" binding:"max=200"`
}

// ReviewVacationRequest represents the approval/rejection request
type ReviewVacationRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
	Reason string `json:"reason,omitempty" binding:"max=200"`
}

// ============================================
// Settings Requests (Admin)
// ============================================

// UpdateSettingsRequest represents the settings update request
type UpdateSettingsRequest struct {
	WeekendPolicy       *WeekendPolicyRequest    `json:"weekendPolicy,omitempty"`
	Newsletter          *NewsletterConfigRequest `json:"newsletter,omitempty"`
	DefaultVacationDays *int                     `json:"defaultVacationDays,omitempty" binding:"omitempty,min=0,max=365"`
	VacationResetMonth  *int                     `json:"vacationResetMonth,omitempty" binding:"omitempty,min=1,max=12"`
}

// WeekendPolicyRequest represents weekend policy settings
type WeekendPolicyRequest struct {
	ExcludeWeekends *bool  `json:"excludeWeekends,omitempty"`
	ExcludedDays    *[]int `json:"excludedDays,omitempty"`
}

// NewsletterConfigRequest represents newsletter settings
type NewsletterConfigRequest struct {
	Enabled    *bool   `json:"enabled,omitempty"`
	Frequency  *string `json:"frequency,omitempty" binding:"omitempty,oneof=weekly monthly"`
	DayOfMonth *int    `json:"dayOfMonth,omitempty" binding:"omitempty,min=1,max=28"`
}

// ============================================
// Email Test Requests (Admin)
// ============================================

// TestEmailRequest represents a request to send a test email
type TestEmailRequest struct {
	Template string `json:"template" binding:"required,oneof=welcome request_submitted request_approved request_rejected admin_notification newsletter"`
}

// PreviewEmailRequest represents a request to preview an email template
type PreviewEmailRequest struct {
	Template string `json:"template" binding:"required,oneof=welcome request_submitted request_approved request_rejected admin_notification newsletter"`
}
