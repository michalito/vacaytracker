package dto

import (
	"fmt"
	"net/http"
)

// Error codes
const (
	// Authentication errors
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrAuthTokenMissing   = "AUTH_TOKEN_MISSING"
	ErrAuthTokenInvalid   = "AUTH_TOKEN_INVALID"
	ErrAuthTokenExpired   = "AUTH_TOKEN_EXPIRED"

	// Authorization errors
	ErrAdminRequired    = "ADMIN_REQUIRED"
	ErrForbidden        = "FORBIDDEN"
	ErrUnauthorized     = "UNAUTHORIZED"

	// Validation errors
	ErrValidation       = "VALIDATION_ERROR"
	ErrInvalidDateRange = "INVALID_DATE_RANGE"
	ErrDateInPast       = "DATE_IN_PAST"
	ErrInvalidInput     = "INVALID_INPUT"

	// Resource errors
	ErrUserNotFound     = "USER_NOT_FOUND"
	ErrRequestNotFound  = "REQUEST_NOT_FOUND"
	ErrSettingsNotFound = "SETTINGS_NOT_FOUND"
	ErrNotFound         = "NOT_FOUND"
	ErrAlreadyExists    = "ALREADY_EXISTS"

	// Business logic errors
	ErrInsufficientBalance   = "INSUFFICIENT_BALANCE"
	ErrCannotCancelApproved  = "CANNOT_CANCEL_APPROVED"
	ErrCannotCancelRejected  = "CANNOT_CANCEL_REJECTED"
	ErrOverlappingRequest    = "OVERLAPPING_REQUEST"
	ErrInvalidStatus         = "INVALID_STATUS"

	// Rate limiting errors
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"

	// Server errors
	ErrInternal = "INTERNAL_ERROR"
	ErrDatabase = "DATABASE_ERROR"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// AppError represents an application error with HTTP status
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ToResponse converts AppError to ErrorResponse for API output
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}

// NewAppError creates a new AppError
func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithDetails adds details to an AppError
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// ============================================
// Error Constructors
// ============================================

// ErrInvalidCredentialsError returns an invalid credentials error
func ErrInvalidCredentialsError() *AppError {
	return NewAppError(ErrInvalidCredentials, "Invalid email or password", http.StatusUnauthorized)
}

// ErrTokenMissingError returns a token missing error
func ErrTokenMissingError() *AppError {
	return NewAppError(ErrAuthTokenMissing, "Authorization token is required", http.StatusUnauthorized)
}

// ErrTokenInvalidError returns a token invalid error
func ErrTokenInvalidError() *AppError {
	return NewAppError(ErrAuthTokenInvalid, "Invalid or malformed token", http.StatusUnauthorized)
}

// ErrTokenExpiredError returns a token expired error
func ErrTokenExpiredError() *AppError {
	return NewAppError(ErrAuthTokenExpired, "Token has expired", http.StatusUnauthorized)
}

// ErrAdminRequiredError returns an admin required error
func ErrAdminRequiredError() *AppError {
	return NewAppError(ErrAdminRequired, "Admin privileges required", http.StatusForbidden)
}

// ErrForbiddenError returns a forbidden error
func ErrForbiddenError(message string) *AppError {
	return NewAppError(ErrForbidden, message, http.StatusForbidden)
}

// ErrValidationError returns a validation error
func ErrValidationError(message string) *AppError {
	return NewAppError(ErrValidation, message, http.StatusBadRequest)
}

// ErrInvalidDateRangeError returns an invalid date range error
func ErrInvalidDateRangeError() *AppError {
	return NewAppError(ErrInvalidDateRange, "End date must be after start date", http.StatusBadRequest)
}

// ErrDateInPastError returns a date in past error
func ErrDateInPastError() *AppError {
	return NewAppError(ErrDateInPast, "Start date cannot be in the past", http.StatusBadRequest)
}

// ErrUserNotFoundError returns a user not found error
func ErrUserNotFoundError() *AppError {
	return NewAppError(ErrUserNotFound, "User not found", http.StatusNotFound)
}

// ErrRequestNotFoundError returns a request not found error
func ErrRequestNotFoundError() *AppError {
	return NewAppError(ErrRequestNotFound, "Vacation request not found", http.StatusNotFound)
}

// ErrNotFoundError returns a generic not found error
func ErrNotFoundError(resource string) *AppError {
	return NewAppError(ErrNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

// ErrAlreadyExistsError returns an already exists error
func ErrAlreadyExistsError(resource string) *AppError {
	return NewAppError(ErrAlreadyExists, fmt.Sprintf("%s already exists", resource), http.StatusConflict)
}

// ErrInsufficientBalanceError returns an insufficient balance error
func ErrInsufficientBalanceError(requested, available int) *AppError {
	return NewAppError(
		ErrInsufficientBalance,
		fmt.Sprintf("Insufficient vacation balance: requested %d days, available %d days", requested, available),
		http.StatusUnprocessableEntity,
	).WithDetails(map[string]interface{}{
		"requested": requested,
		"available": available,
	})
}

// ErrCannotCancelError returns a cannot cancel error
func ErrCannotCancelError(status string) *AppError {
	return NewAppError(
		ErrInvalidStatus,
		fmt.Sprintf("Cannot cancel a request with status: %s", status),
		http.StatusBadRequest,
	)
}

// ErrOverlappingRequestError returns an overlapping request error
func ErrOverlappingRequestError() *AppError {
	return NewAppError(ErrOverlappingRequest, "Request overlaps with an existing vacation", http.StatusConflict)
}

// ErrInternalError returns an internal server error
func ErrInternalError() *AppError {
	return NewAppError(ErrInternal, "An internal error occurred", http.StatusInternalServerError)
}

// ErrDatabaseError returns a database error
func ErrDatabaseError(err error) *AppError {
	return NewAppError(ErrDatabase, "Database operation failed", http.StatusInternalServerError)
}

// ErrInternalErrorWithMessage returns an internal server error with a custom message
func ErrInternalErrorWithMessage(message string) *AppError {
	return NewAppError(ErrInternal, message, http.StatusInternalServerError)
}

// ErrConflictError returns a conflict error (e.g., duplicate resource)
func ErrConflictError(message string) *AppError {
	return NewAppError(ErrAlreadyExists, message, http.StatusConflict)
}
