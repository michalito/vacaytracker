package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/repository/sqlite"
)

// Vacation service errors
var (
	ErrInsufficientBalance  = errors.New("insufficient vacation balance")
	ErrInvalidDateRange     = errors.New("end date must be after start date")
	ErrDateInPast           = errors.New("start date cannot be in the past")
	ErrRequestNotFound      = errors.New("vacation request not found")
	ErrCannotCancelApproved = errors.New("cannot cancel approved request")
	ErrCannotCancelRejected = errors.New("cannot cancel rejected request")
	ErrAlreadyProcessed     = errors.New("request already processed")
	ErrAccessDenied         = errors.New("access denied")
)

// VacationService handles vacation request business logic
type VacationService struct {
	vacationRepo *sqlite.VacationRepository
	userRepo     *sqlite.UserRepository
	settingsRepo *sqlite.SettingsRepository
}

// NewVacationService creates a new VacationService
func NewVacationService(
	vacationRepo *sqlite.VacationRepository,
	userRepo *sqlite.UserRepository,
	settingsRepo *sqlite.SettingsRepository,
) *VacationService {
	return &VacationService{
		vacationRepo: vacationRepo,
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
	}
}

// Create creates a new vacation request
func (s *VacationService) Create(ctx context.Context, userID string, req dto.CreateVacationRequest) (*domain.VacationRequest, error) {
	// Parse dates (DD/MM/YYYY -> time.Time)
	startDate, err := parseDDMMYYYY(req.StartDate)
	if err != nil {
		return nil, dto.ErrValidationError(fmt.Sprintf("invalid start date format: %v", err))
	}

	endDate, err := parseDDMMYYYY(req.EndDate)
	if err != nil {
		return nil, dto.ErrValidationError(fmt.Sprintf("invalid end date format: %v", err))
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, dto.ErrValidationError("end date must be after or equal to start date")
	}

	// Check if start date is in the past
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if startDate.Before(today) {
		return nil, dto.ErrValidationError("start date cannot be in the past")
	}

	// Get settings for business day calculation
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get settings")
	}

	// Calculate business days
	totalDays := calculateBusinessDays(startDate, endDate, settings.WeekendPolicy)
	if totalDays == 0 {
		return nil, dto.ErrValidationError("selected dates result in zero vacation days")
	}

	// Get user and check balance
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return nil, dto.ErrNotFoundError("user")
	}

	if user.VacationBalance < totalDays {
		return nil, dto.ErrInsufficientBalanceError(user.VacationBalance, totalDays)
	}

	// Format dates for storage
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	// Check for overlapping requests
	hasOverlap, err := s.vacationRepo.HasOverlap(ctx, userID, startDateStr, endDateStr)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to check for overlapping requests")
	}
	if hasOverlap {
		return nil, dto.ErrOverlappingRequestError()
	}

	// Create request - auto-approve for admins
	status := domain.StatusPending
	if user.IsAdmin() {
		status = domain.StatusApproved
	}

	vacation := &domain.VacationRequest{
		ID:        uuid.New().String(),
		UserID:    userID,
		StartDate: startDateStr,
		EndDate:   endDateStr,
		TotalDays: totalDays,
		Status:    status,
	}

	if req.Reason != "" {
		vacation.Reason = &req.Reason
	}

	// For admins, create request and deduct balance atomically
	if user.IsAdmin() {
		newBalance := user.VacationBalance - totalDays
		if newBalance < 0 {
			newBalance = 0
		}

		db := s.vacationRepo.GetDB()
		err = db.Transaction(func(tx *sql.Tx) error {
			if err := s.vacationRepo.CreateTx(ctx, tx, vacation); err != nil {
				return err
			}
			if err := s.userRepo.UpdateVacationBalanceTx(ctx, tx, userID, newBalance); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return nil, dto.ErrInternalErrorWithMessage("failed to create vacation request")
		}
	} else {
		if err := s.vacationRepo.Create(ctx, vacation); err != nil {
			return nil, dto.ErrInternalErrorWithMessage("failed to create vacation request")
		}
	}

	// Fetch the created request with user info
	return s.vacationRepo.GetByID(ctx, vacation.ID)
}

// Cancel cancels a pending vacation request
func (s *VacationService) Cancel(ctx context.Context, requestID, userID string) error {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return dto.ErrInternalErrorWithMessage("failed to get vacation request")
	}
	if request == nil {
		return dto.ErrNotFoundError("vacation request")
	}

	// Check ownership
	if request.UserID != userID {
		return dto.ErrForbiddenError("you can only cancel your own requests")
	}

	// Check status
	if request.IsApproved() {
		return dto.ErrForbiddenError("cannot cancel approved request")
	}
	if request.IsRejected() {
		return dto.ErrForbiddenError("cannot cancel rejected request")
	}

	return s.vacationRepo.Delete(ctx, requestID)
}

// Approve approves a pending request and deducts balance atomically using a transaction
func (s *VacationService) Approve(ctx context.Context, requestID, adminID string) (*domain.VacationRequest, error) {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get vacation request")
	}
	if request == nil {
		return nil, dto.ErrNotFoundError("vacation request")
	}

	if !request.IsPending() {
		return nil, dto.ErrConflictError("request has already been processed")
	}

	// Get user to check balance
	user, err := s.userRepo.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return nil, dto.ErrNotFoundError("user")
	}

	// Check if user still has enough balance
	if user.VacationBalance < request.TotalDays {
		return nil, dto.ErrInsufficientBalanceError(user.VacationBalance, request.TotalDays)
	}

	// Calculate new balance
	newBalance := user.VacationBalance - request.TotalDays
	if newBalance < 0 {
		newBalance = 0
	}

	// Execute status update and balance deduction atomically in a transaction
	db := s.vacationRepo.GetDB()
	err = db.Transaction(func(tx *sql.Tx) error {
		// Update status
		if err := s.vacationRepo.UpdateStatusTx(ctx, tx, requestID, domain.StatusApproved, adminID, nil); err != nil {
			return err
		}

		// Deduct vacation balance
		if err := s.userRepo.UpdateVacationBalanceTx(ctx, tx, request.UserID, newBalance); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to approve request")
	}

	// Fetch updated request
	return s.vacationRepo.GetByID(ctx, requestID)
}

// Reject rejects a pending request
func (s *VacationService) Reject(ctx context.Context, requestID, adminID string, reason *string) (*domain.VacationRequest, error) {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get vacation request")
	}
	if request == nil {
		return nil, dto.ErrNotFoundError("vacation request")
	}

	if !request.IsPending() {
		return nil, dto.ErrConflictError("request has already been processed")
	}

	if err := s.vacationRepo.UpdateStatus(ctx, requestID, domain.StatusRejected, adminID, reason); err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to reject request")
	}

	return s.vacationRepo.GetByID(ctx, requestID)
}

// GetByID retrieves a vacation request by ID
func (s *VacationService) GetByID(ctx context.Context, requestID string) (*domain.VacationRequest, error) {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get vacation request")
	}
	if request == nil {
		return nil, dto.ErrNotFoundError("vacation request")
	}
	return request, nil
}

// ListByUser retrieves vacation requests for a user
func (s *VacationService) ListByUser(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
	requests, err := s.vacationRepo.ListByUser(ctx, userID, status, year)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to list vacation requests")
	}
	return requests, nil
}

// ListPending retrieves all pending vacation requests (for admin)
func (s *VacationService) ListPending(ctx context.Context) ([]*domain.VacationRequest, error) {
	requests, err := s.vacationRepo.ListPending(ctx)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to list pending requests")
	}
	return requests, nil
}

// ListTeam retrieves team vacations for a given month/year
func (s *VacationService) ListTeam(ctx context.Context, month, year int) ([]*domain.TeamVacation, error) {
	if month < 1 || month > 12 {
		return nil, dto.ErrValidationError("month must be between 1 and 12")
	}
	if year < 2000 || year > 2100 {
		return nil, dto.ErrValidationError("invalid year")
	}

	vacations, err := s.vacationRepo.ListTeam(ctx, month, year)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to list team vacations")
	}
	return vacations, nil
}

// parseDDMMYYYY parses DD/MM/YYYY format to time.Time
func parseDDMMYYYY(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "/")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid date format, expected DD/MM/YYYY")
	}

	// Validate parts
	if len(parts[0]) < 1 || len(parts[0]) > 2 {
		return time.Time{}, fmt.Errorf("invalid day")
	}
	if len(parts[1]) < 1 || len(parts[1]) > 2 {
		return time.Time{}, fmt.Errorf("invalid month")
	}
	if len(parts[2]) != 4 {
		return time.Time{}, fmt.Errorf("invalid year")
	}

	// Pad day and month with zeros if needed
	day := parts[0]
	if len(day) == 1 {
		day = "0" + day
	}
	month := parts[1]
	if len(month) == 1 {
		month = "0" + month
	}

	// Rearrange to YYYY-MM-DD
	isoDate := fmt.Sprintf("%s-%s-%s", parts[2], month, day)
	return time.Parse("2006-01-02", isoDate)
}

// calculateBusinessDays counts business days between two dates
func calculateBusinessDays(start, end time.Time, policy domain.WeekendPolicy) int {
	if !policy.ExcludeWeekends {
		// Include all days
		return int(end.Sub(start).Hours()/24) + 1
	}

	count := 0
	current := start

	// Create a map of excluded weekdays for faster lookup
	excluded := make(map[time.Weekday]bool)
	for _, day := range policy.ExcludedDays {
		excluded[time.Weekday(day)] = true
	}

	for !current.After(end) {
		if !excluded[current.Weekday()] {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}

	return count
}
