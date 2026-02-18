package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newTestEmployee returns a domain.User with the employee role and the given balance.
func newTestEmployee(id string, balance int) *domain.User {
	return &domain.User{
		ID:              id,
		Email:           id + "@example.com",
		Name:            "Test Employee",
		Role:            domain.RoleEmployee,
		VacationBalance: balance,
	}
}

// newTestAdmin returns a domain.User with the admin role and the given balance.
func newTestAdmin(id string, balance int) *domain.User {
	return &domain.User{
		ID:              id,
		Email:           id + "@example.com",
		Name:            "Test Admin",
		Role:            domain.RoleAdmin,
		VacationBalance: balance,
	}
}

// newPendingRequest returns a domain.VacationRequest in pending status.
func newPendingRequest(id, userID string, totalDays int) *domain.VacationRequest {
	return &domain.VacationRequest{
		ID:        id,
		UserID:    userID,
		UserName:  "Test User",
		UserEmail: userID + "@example.com",
		StartDate: "2027-06-16",
		EndDate:   "2027-06-20",
		TotalDays: totalDays,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// newApprovedRequest returns a domain.VacationRequest in approved status.
func newApprovedRequest(id, userID string, totalDays int) *domain.VacationRequest {
	r := newPendingRequest(id, userID, totalDays)
	r.Status = domain.StatusApproved
	reviewedBy := "admin-1"
	reviewedAt := time.Now()
	r.ReviewedBy = &reviewedBy
	r.ReviewedAt = &reviewedAt
	return r
}

// newRejectedRequest returns a domain.VacationRequest in rejected status.
func newRejectedRequest(id, userID string, totalDays int) *domain.VacationRequest {
	r := newPendingRequest(id, userID, totalDays)
	r.Status = domain.StatusRejected
	reviewedBy := "admin-1"
	reviewedAt := time.Now()
	reason := "not enough coverage"
	r.ReviewedBy = &reviewedBy
	r.ReviewedAt = &reviewedAt
	r.RejectionReason = &reason
	return r
}

// assertVacationAppError verifies the error is an *dto.AppError with the expected code.
func assertVacationAppError(t *testing.T, err error, expectedCode string) {
	t.Helper()
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr, "expected an *dto.AppError")
	assert.Equal(t, expectedCode, appErr.Code)
}

// newServiceBundle wires up the mock repositories and returns the service plus
// references to each mock so tests can configure per-test behaviour.
type serviceDeps struct {
	svc          *service.VacationService
	vacationRepo *testutil.MockVacationRepository
	userRepo     *testutil.MockUserRepository
	settingsRepo *testutil.MockSettingsRepository
	transactor   *testutil.MockTransactor
}

func newServiceBundle() *serviceDeps {
	vr := &testutil.MockVacationRepository{}
	ur := &testutil.MockUserRepository{}
	sr := &testutil.MockSettingsRepository{}
	tx := &testutil.MockTransactor{}
	svc := service.NewVacationService(vr, ur, sr, tx)
	return &serviceDeps{
		svc:          svc,
		vacationRepo: vr,
		userRepo:     ur,
		settingsRepo: sr,
		transactor:   tx,
	}
}

// =========================================================================
// Create
// =========================================================================

func TestCreate_EmployeeCreatesPendingRequest(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}
	var createdReq *domain.VacationRequest
	d.vacationRepo.CreateFn = func(_ context.Context, req *domain.VacationRequest) error {
		createdReq = req
		return nil
	}
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if createdReq != nil && createdReq.ID == id {
			return createdReq, nil
		}
		return nil, nil
	}

	// 14/06/2027 is Monday, 18/06/2027 is Friday => 5 business days
	result, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "14/06/2027",
		EndDate:   "18/06/2027",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.StatusPending, result.Status)
	assert.Equal(t, 5, result.TotalDays)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, "2027-06-14", result.StartDate)
	assert.Equal(t, "2027-06-18", result.EndDate)
	assert.Nil(t, result.Reason)
}

func TestCreate_EmployeeWithReason(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}
	var createdReq *domain.VacationRequest
	d.vacationRepo.CreateFn = func(_ context.Context, req *domain.VacationRequest) error {
		createdReq = req
		return nil
	}
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if createdReq != nil && createdReq.ID == id {
			return createdReq, nil
		}
		return nil, nil
	}

	result, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "14/06/2027",
		EndDate:   "18/06/2027",
		Reason:    "Family vacation",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Reason)
	assert.Equal(t, "Family vacation", *result.Reason)
}

func TestCreate_AdminAutoApproves(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	adminID := "admin-1"
	admin := newTestAdmin(adminID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == adminID {
			return admin, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}

	var createdReq *domain.VacationRequest
	var balanceUpdated bool
	d.vacationRepo.CreateTxFn = func(_ context.Context, _ *sql.Tx, req *domain.VacationRequest) error {
		createdReq = req
		return nil
	}
	d.userRepo.UpdateVacationBalanceTxFn = func(_ context.Context, _ *sql.Tx, id string, balance int) error {
		assert.Equal(t, adminID, id)
		assert.Equal(t, 15, balance) // 20 - 5
		balanceUpdated = true
		return nil
	}
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if createdReq != nil && createdReq.ID == id {
			return createdReq, nil
		}
		return nil, nil
	}

	result, err := d.svc.Create(ctx, adminID, dto.CreateVacationRequest{
		StartDate: "14/06/2027",
		EndDate:   "18/06/2027",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.StatusApproved, result.Status)
	assert.True(t, balanceUpdated, "balance should have been deducted in transaction")
}

func TestCreate_InvalidStartDateFormat(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "2027-06-16", // wrong format: ISO instead of DD/MM/YYYY
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
}

func TestCreate_InvalidEndDateFormat(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "invalid-date",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
}

func TestCreate_EndBeforeStart(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "20/06/2027",
		EndDate:   "16/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "end date must be after or equal to start date")
}

func TestCreate_StartInPast(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "01/01/2020",
		EndDate:   "05/01/2020",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "start date cannot be in the past")
}

func TestCreate_ZeroBusinessDays_WeekendOnly(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	// 19/06/2027 is Saturday, 20/06/2027 is Sunday
	// With default settings (exclude weekends), this yields 0 business days.
	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "19/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "zero vacation days")
}

func TestCreate_InsufficientBalance(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 2) // only 2 days available

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}

	// 14/06/2027 Mon - 18/06/2027 Fri => 5 business days, but balance is 2
	_, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "14/06/2027",
		EndDate:   "18/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInsufficientBalance)
}

func TestCreate_OverlappingRequest(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return true, nil
	}

	_, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrOverlappingRequest)
}

func TestCreate_UserNotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	// userRepo.GetByID returns nil by default (user not found)
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}

	_, err := d.svc.Create(ctx, "nonexistent", dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestCreate_SettingsRepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.settingsRepo.GetFn = func(_ context.Context) (*domain.Settings, error) {
		return nil, errors.New("db connection lost")
	}

	_, err := d.svc.Create(ctx, "emp-1", dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

func TestCreate_SingleDay(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}
	var createdReq *domain.VacationRequest
	d.vacationRepo.CreateFn = func(_ context.Context, req *domain.VacationRequest) error {
		createdReq = req
		return nil
	}
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if createdReq != nil && createdReq.ID == id {
			return createdReq, nil
		}
		return nil, nil
	}

	// 16/06/2027 is a Wednesday => 1 business day
	result, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "16/06/2027",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.TotalDays)
}

func TestCreate_CreateRepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}
	d.vacationRepo.CreateFn = func(_ context.Context, _ *domain.VacationRequest) error {
		return errors.New("write failed")
	}

	_, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

func TestCreate_AdminTransactionError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	adminID := "admin-1"
	admin := newTestAdmin(adminID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == adminID {
			return admin, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, nil
	}
	d.transactor.TransactionFn = func(_ func(tx *sql.Tx) error) error {
		return errors.New("transaction failed")
	}

	_, err := d.svc.Create(ctx, adminID, dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

func TestCreate_HasOverlapRepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	employee := newTestEmployee(userID, 20)

	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return employee, nil
		}
		return nil, nil
	}
	d.vacationRepo.HasOverlapFn = func(_ context.Context, _, _, _ string) (bool, error) {
		return false, errors.New("db error")
	}

	_, err := d.svc.Create(ctx, userID, dto.CreateVacationRequest{
		StartDate: "16/06/2027",
		EndDate:   "20/06/2027",
	})

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// Cancel
// =========================================================================

func TestCancel_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newPendingRequest(requestID, userID, 5), nil
		}
		return nil, nil
	}
	deleted := false
	d.vacationRepo.DeleteFn = func(_ context.Context, id string) error {
		assert.Equal(t, requestID, id)
		deleted = true
		return nil
	}

	err := d.svc.Cancel(ctx, requestID, userID)

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestCancel_NotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	// GetByID returns nil by default
	err := d.svc.Cancel(ctx, "nonexistent", "emp-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestCancel_NotOwner(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newPendingRequest(requestID, "emp-1", 5), nil
		}
		return nil, nil
	}

	err := d.svc.Cancel(ctx, requestID, "emp-2") // different user

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrForbidden)
	assert.Contains(t, err.Error(), "your own requests")
}

func TestCancel_AlreadyApproved(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newApprovedRequest(requestID, userID, 5), nil
		}
		return nil, nil
	}

	err := d.svc.Cancel(ctx, requestID, userID)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrForbidden)
	assert.Contains(t, err.Error(), "approved")
}

func TestCancel_AlreadyRejected(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newRejectedRequest(requestID, userID, 5), nil
		}
		return nil, nil
	}

	err := d.svc.Cancel(ctx, requestID, userID)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrForbidden)
	assert.Contains(t, err.Error(), "rejected")
}

func TestCancel_GetByIDRepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.GetByIDFn = func(_ context.Context, _ string) (*domain.VacationRequest, error) {
		return nil, errors.New("db failure")
	}

	err := d.svc.Cancel(ctx, "req-1", "emp-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// Approve
// =========================================================================

func TestApprove_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	adminID := "admin-1"
	userID := "emp-1"
	requestID := "req-1"
	totalDays := 5
	initialBalance := 20

	pendingReq := newPendingRequest(requestID, userID, totalDays)
	user := newTestEmployee(userID, initialBalance)

	callCount := 0
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			callCount++
			if callCount == 1 {
				return pendingReq, nil
			}
			// After approval, return an approved copy
			approved := *pendingReq
			approved.Status = domain.StatusApproved
			return &approved, nil
		}
		return nil, nil
	}
	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return user, nil
		}
		return nil, nil
	}

	var statusUpdated, balanceDeducted bool
	d.vacationRepo.UpdateStatusTxFn = func(_ context.Context, _ *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, reason *string) error {
		assert.Equal(t, requestID, id)
		assert.Equal(t, domain.StatusApproved, status)
		assert.Equal(t, adminID, reviewedBy)
		assert.Nil(t, reason)
		statusUpdated = true
		return nil
	}
	d.userRepo.UpdateVacationBalanceTxFn = func(_ context.Context, _ *sql.Tx, id string, balance int) error {
		assert.Equal(t, userID, id)
		assert.Equal(t, initialBalance-totalDays, balance) // 20 - 5 = 15
		balanceDeducted = true
		return nil
	}

	result, err := d.svc.Approve(ctx, requestID, adminID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.StatusApproved, result.Status)
	assert.True(t, statusUpdated)
	assert.True(t, balanceDeducted)
}

func TestApprove_NotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Approve(ctx, "nonexistent", "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestApprove_AlreadyProcessed_Approved(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newApprovedRequest(requestID, "emp-1", 5), nil
		}
		return nil, nil
	}

	_, err := d.svc.Approve(ctx, requestID, "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrAlreadyExists) // ErrConflictError uses ErrAlreadyExists code
	assert.Contains(t, err.Error(), "already been processed")
}

func TestApprove_AlreadyProcessed_Rejected(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newRejectedRequest(requestID, "emp-1", 5), nil
		}
		return nil, nil
	}

	_, err := d.svc.Approve(ctx, requestID, "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrAlreadyExists)
}

func TestApprove_InsufficientBalance(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	requestID := "req-1"

	pendingReq := newPendingRequest(requestID, userID, 10)
	user := newTestEmployee(userID, 3) // only 3 days left

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return pendingReq, nil
		}
		return nil, nil
	}
	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return user, nil
		}
		return nil, nil
	}

	_, err := d.svc.Approve(ctx, requestID, "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInsufficientBalance)

	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 10, appErr.Details["requested"])
	assert.Equal(t, 3, appErr.Details["available"])
}

func TestApprove_UserNotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newPendingRequest(requestID, "ghost-user", 5), nil
		}
		return nil, nil
	}
	// userRepo.GetByID returns nil by default => user not found

	_, err := d.svc.Approve(ctx, requestID, "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestApprove_TransactionError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newPendingRequest(requestID, userID, 5), nil
		}
		return nil, nil
	}
	d.userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == userID {
			return newTestEmployee(userID, 20), nil
		}
		return nil, nil
	}
	d.transactor.TransactionFn = func(_ func(tx *sql.Tx) error) error {
		return errors.New("transaction failed")
	}

	_, err := d.svc.Approve(ctx, requestID, "admin-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// Reject
// =========================================================================

func TestReject_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	adminID := "admin-1"
	requestID := "req-1"
	userID := "emp-1"

	callCount := 0
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			callCount++
			if callCount == 1 {
				return newPendingRequest(requestID, userID, 5), nil
			}
			return newRejectedRequest(requestID, userID, 5), nil
		}
		return nil, nil
	}

	var statusUpdated bool
	d.vacationRepo.UpdateStatusFn = func(_ context.Context, id string, status domain.VacationStatus, reviewedBy string, reason *string) error {
		assert.Equal(t, requestID, id)
		assert.Equal(t, domain.StatusRejected, status)
		assert.Equal(t, adminID, reviewedBy)
		assert.Nil(t, reason)
		statusUpdated = true
		return nil
	}

	result, err := d.svc.Reject(ctx, requestID, adminID, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, domain.StatusRejected, result.Status)
	assert.True(t, statusUpdated)
}

func TestReject_WithReason(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	adminID := "admin-1"
	requestID := "req-1"
	userID := "emp-1"
	reason := "not enough team coverage"

	callCount := 0
	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			callCount++
			if callCount == 1 {
				return newPendingRequest(requestID, userID, 5), nil
			}
			rejected := newRejectedRequest(requestID, userID, 5)
			rejected.RejectionReason = &reason
			return rejected, nil
		}
		return nil, nil
	}

	d.vacationRepo.UpdateStatusFn = func(_ context.Context, _ string, _ domain.VacationStatus, _ string, r *string) error {
		require.NotNil(t, r)
		assert.Equal(t, reason, *r)
		return nil
	}

	result, err := d.svc.Reject(ctx, requestID, adminID, &reason)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.RejectionReason)
	assert.Equal(t, reason, *result.RejectionReason)
}

func TestReject_NotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.Reject(ctx, "nonexistent", "admin-1", nil)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestReject_AlreadyProcessed(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newApprovedRequest(requestID, "emp-1", 5), nil
		}
		return nil, nil
	}

	_, err := d.svc.Reject(ctx, requestID, "admin-1", nil)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrAlreadyExists)
	assert.Contains(t, err.Error(), "already been processed")
}

func TestReject_UpdateStatusError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return newPendingRequest(requestID, "emp-1", 5), nil
		}
		return nil, nil
	}
	d.vacationRepo.UpdateStatusFn = func(_ context.Context, _ string, _ domain.VacationStatus, _ string, _ *string) error {
		return errors.New("update failed")
	}

	_, err := d.svc.Reject(ctx, requestID, "admin-1", nil)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// GetByID
// =========================================================================

func TestVacationGetByID_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	requestID := "req-1"
	expected := newPendingRequest(requestID, "emp-1", 5)

	d.vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == requestID {
			return expected, nil
		}
		return nil, nil
	}

	result, err := d.svc.GetByID(ctx, requestID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, requestID, result.ID)
	assert.Equal(t, expected.UserID, result.UserID)
	assert.Equal(t, expected.TotalDays, result.TotalDays)
}

func TestVacationGetByID_NotFound(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.GetByID(ctx, "nonexistent")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrNotFound)
}

func TestVacationGetByID_RepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.GetByIDFn = func(_ context.Context, _ string) (*domain.VacationRequest, error) {
		return nil, errors.New("db failure")
	}

	_, err := d.svc.GetByID(ctx, "req-1")

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// ListByUser
// =========================================================================

func TestListByUser_NoFilters(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	expected := []*domain.VacationRequest{
		newPendingRequest("req-1", userID, 5),
		newApprovedRequest("req-2", userID, 3),
	}

	d.vacationRepo.ListByUserFn = func(_ context.Context, uid string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, userID, uid)
		assert.Nil(t, status)
		assert.Nil(t, year)
		return expected, nil
	}

	results, err := d.svc.ListByUser(ctx, userID, nil, nil)

	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestListByUser_FilterByStatus(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	status := domain.StatusPending

	d.vacationRepo.ListByUserFn = func(_ context.Context, uid string, s *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, userID, uid)
		require.NotNil(t, s)
		assert.Equal(t, domain.StatusPending, *s)
		assert.Nil(t, year)
		return []*domain.VacationRequest{newPendingRequest("req-1", userID, 5)}, nil
	}

	results, err := d.svc.ListByUser(ctx, userID, &status, nil)

	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestListByUser_FilterByYear(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	year := 2027

	d.vacationRepo.ListByUserFn = func(_ context.Context, uid string, status *domain.VacationStatus, y *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, userID, uid)
		assert.Nil(t, status)
		require.NotNil(t, y)
		assert.Equal(t, 2027, *y)
		return []*domain.VacationRequest{newPendingRequest("req-1", userID, 5)}, nil
	}

	results, err := d.svc.ListByUser(ctx, userID, nil, &year)

	require.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestListByUser_FilterByStatusAndYear(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	userID := "emp-1"
	status := domain.StatusApproved
	year := 2027

	d.vacationRepo.ListByUserFn = func(_ context.Context, uid string, s *domain.VacationStatus, y *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, userID, uid)
		require.NotNil(t, s)
		assert.Equal(t, domain.StatusApproved, *s)
		require.NotNil(t, y)
		assert.Equal(t, 2027, *y)
		return []*domain.VacationRequest{}, nil
	}

	results, err := d.svc.ListByUser(ctx, userID, &status, &year)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListByUser_RepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListByUserFn = func(_ context.Context, _ string, _ *domain.VacationStatus, _ *int) ([]*domain.VacationRequest, error) {
		return nil, errors.New("db error")
	}

	_, err := d.svc.ListByUser(ctx, "emp-1", nil, nil)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// ListPending
// =========================================================================

func TestListPending_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	expected := []*domain.VacationRequest{
		newPendingRequest("req-1", "emp-1", 5),
		newPendingRequest("req-2", "emp-2", 3),
	}

	d.vacationRepo.ListPendingFn = func(_ context.Context) ([]*domain.VacationRequest, error) {
		return expected, nil
	}

	results, err := d.svc.ListPending(ctx)

	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, domain.StatusPending, results[0].Status)
	assert.Equal(t, domain.StatusPending, results[1].Status)
}

func TestListPending_Empty(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListPendingFn = func(_ context.Context) ([]*domain.VacationRequest, error) {
		return []*domain.VacationRequest{}, nil
	}

	results, err := d.svc.ListPending(ctx)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListPending_RepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListPendingFn = func(_ context.Context) ([]*domain.VacationRequest, error) {
		return nil, errors.New("db error")
	}

	_, err := d.svc.ListPending(ctx)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}

// =========================================================================
// ListTeam
// =========================================================================

func TestListTeam_Success(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()
	expected := []*domain.TeamVacation{
		{
			ID:        "req-1",
			UserID:    "emp-1",
			UserName:  "Alice",
			StartDate: "2027-06-16",
			EndDate:   "2027-06-20",
			TotalDays: 5,
		},
		{
			ID:        "req-2",
			UserID:    "emp-2",
			UserName:  "Bob",
			StartDate: "2027-06-23",
			EndDate:   "2027-06-25",
			TotalDays: 3,
		},
	}

	d.vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 6, month)
		assert.Equal(t, 2027, year)
		return expected, nil
	}

	results, err := d.svc.ListTeam(ctx, 6, 2027)

	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Alice", results[0].UserName)
	assert.Equal(t, "Bob", results[1].UserName)
}

func TestListTeam_InvalidMonth_Zero(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.ListTeam(ctx, 0, 2027)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "month must be between 1 and 12")
}

func TestListTeam_InvalidMonth_Thirteen(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.ListTeam(ctx, 13, 2027)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "month must be between 1 and 12")
}

func TestListTeam_InvalidMonth_Negative(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.ListTeam(ctx, -1, 2027)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
}

func TestListTeam_InvalidYear_TooLow(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.ListTeam(ctx, 6, 1999)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "invalid year")
}

func TestListTeam_InvalidYear_TooHigh(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	_, err := d.svc.ListTeam(ctx, 6, 2101)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrValidation)
	assert.Contains(t, err.Error(), "invalid year")
}

func TestListTeam_BoundaryMonth_One(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 1, month)
		return []*domain.TeamVacation{}, nil
	}

	results, err := d.svc.ListTeam(ctx, 1, 2027)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListTeam_BoundaryMonth_Twelve(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 12, month)
		return []*domain.TeamVacation{}, nil
	}

	results, err := d.svc.ListTeam(ctx, 12, 2027)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListTeam_BoundaryYear_2000(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 2000, year)
		return []*domain.TeamVacation{}, nil
	}

	results, err := d.svc.ListTeam(ctx, 6, 2000)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListTeam_BoundaryYear_2100(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 2100, year)
		return []*domain.TeamVacation{}, nil
	}

	results, err := d.svc.ListTeam(ctx, 6, 2100)

	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestListTeam_RepoError(t *testing.T) {
	d := newServiceBundle()
	ctx := context.Background()

	d.vacationRepo.ListTeamFn = func(_ context.Context, _, _ int) ([]*domain.TeamVacation, error) {
		return nil, errors.New("db error")
	}

	_, err := d.svc.ListTeam(ctx, 6, 2027)

	require.Error(t, err)
	assertVacationAppError(t, err, dto.ErrInternal)
}
