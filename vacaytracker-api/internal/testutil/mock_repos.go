package testutil

import (
	"context"
	"database/sql"
	"time"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/repository"
)

// MockUserRepository is a mock implementation of repository.UserRepository.
// Set function fields to customize behavior per test.
type MockUserRepository struct {
	CreateFn                func(ctx context.Context, user *domain.User) error
	GetByIDFn               func(ctx context.Context, id string) (*domain.User, error)
	GetByEmailFn            func(ctx context.Context, email string) (*domain.User, error)
	GetAllFn                func(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error)
	GetByRoleFn             func(ctx context.Context, role domain.Role) ([]*domain.User, error)
	CountByRoleFn           func(ctx context.Context, role domain.Role) (int, error)
	UpdateFn                func(ctx context.Context, user *domain.User) error
	UpdatePasswordFn        func(ctx context.Context, id, passwordHash string) error
	UpdateEmailPreferencesFn func(ctx context.Context, id string, prefs domain.EmailPreferences) error
	UpdateVacationBalanceFn  func(ctx context.Context, id string, balance int) error
	UpdateVacationBalanceTxFn func(ctx context.Context, tx *sql.Tx, id string, balance int) error
	DeleteFn                func(ctx context.Context, id string) error
	EmailExistsFn           func(ctx context.Context, email string) (bool, error)
	EmailExistsExcludingFn  func(ctx context.Context, email, excludeID string) (bool, error)
	GetNewsletterRecipientsFn func(ctx context.Context) ([]*domain.User, error)
	GetLowBalanceUsersFn    func(ctx context.Context, threshold int) ([]*domain.User, error)
	UpdateAllBalancesFn     func(ctx context.Context, balance int) (int64, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx, role, search, limit, offset)
	}
	return nil, 0, nil
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	if m.GetByRoleFn != nil {
		return m.GetByRoleFn(ctx, role)
	}
	return nil, nil
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role domain.Role) (int, error) {
	if m.CountByRoleFn != nil {
		return m.CountByRoleFn(ctx, role)
	}
	return 0, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	if m.UpdatePasswordFn != nil {
		return m.UpdatePasswordFn(ctx, id, passwordHash)
	}
	return nil
}

func (m *MockUserRepository) UpdateEmailPreferences(ctx context.Context, id string, prefs domain.EmailPreferences) error {
	if m.UpdateEmailPreferencesFn != nil {
		return m.UpdateEmailPreferencesFn(ctx, id, prefs)
	}
	return nil
}

func (m *MockUserRepository) UpdateVacationBalance(ctx context.Context, id string, balance int) error {
	if m.UpdateVacationBalanceFn != nil {
		return m.UpdateVacationBalanceFn(ctx, id, balance)
	}
	return nil
}

func (m *MockUserRepository) UpdateVacationBalanceTx(ctx context.Context, tx *sql.Tx, id string, balance int) error {
	if m.UpdateVacationBalanceTxFn != nil {
		return m.UpdateVacationBalanceTxFn(ctx, tx, id, balance)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if m.EmailExistsFn != nil {
		return m.EmailExistsFn(ctx, email)
	}
	return false, nil
}

func (m *MockUserRepository) EmailExistsExcluding(ctx context.Context, email, excludeID string) (bool, error) {
	if m.EmailExistsExcludingFn != nil {
		return m.EmailExistsExcludingFn(ctx, email, excludeID)
	}
	return false, nil
}

func (m *MockUserRepository) GetNewsletterRecipients(ctx context.Context) ([]*domain.User, error) {
	if m.GetNewsletterRecipientsFn != nil {
		return m.GetNewsletterRecipientsFn(ctx)
	}
	return nil, nil
}

func (m *MockUserRepository) GetLowBalanceUsers(ctx context.Context, threshold int) ([]*domain.User, error) {
	if m.GetLowBalanceUsersFn != nil {
		return m.GetLowBalanceUsersFn(ctx, threshold)
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateAllBalances(ctx context.Context, balance int) (int64, error) {
	if m.UpdateAllBalancesFn != nil {
		return m.UpdateAllBalancesFn(ctx, balance)
	}
	return 0, nil
}

// MockVacationRepository is a mock implementation of repository.VacationRepository.
type MockVacationRepository struct {
	CreateFn        func(ctx context.Context, req *domain.VacationRequest) error
	CreateTxFn      func(ctx context.Context, tx *sql.Tx, req *domain.VacationRequest) error
	GetByIDFn       func(ctx context.Context, id string) (*domain.VacationRequest, error)
	ListByUserFn    func(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error)
	ListPendingFn   func(ctx context.Context) ([]*domain.VacationRequest, error)
	ListTeamFn      func(ctx context.Context, month, year int) ([]*domain.TeamVacation, error)
	UpdateStatusFn  func(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error
	UpdateStatusTxFn func(ctx context.Context, tx *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error
	DeleteFn        func(ctx context.Context, id string) error
	HasOverlapFn    func(ctx context.Context, userID, startDate, endDate string) (bool, error)
	GetMonthlyStatsFn func(ctx context.Context, year, month int) (*repository.MonthlyStats, error)
}

func (m *MockVacationRepository) Create(ctx context.Context, req *domain.VacationRequest) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, req)
	}
	return nil
}

func (m *MockVacationRepository) CreateTx(ctx context.Context, tx *sql.Tx, req *domain.VacationRequest) error {
	if m.CreateTxFn != nil {
		return m.CreateTxFn(ctx, tx, req)
	}
	return nil
}

func (m *MockVacationRepository) GetByID(ctx context.Context, id string) (*domain.VacationRequest, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockVacationRepository) ListByUser(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
	if m.ListByUserFn != nil {
		return m.ListByUserFn(ctx, userID, status, year)
	}
	return nil, nil
}

func (m *MockVacationRepository) ListPending(ctx context.Context) ([]*domain.VacationRequest, error) {
	if m.ListPendingFn != nil {
		return m.ListPendingFn(ctx)
	}
	return nil, nil
}

func (m *MockVacationRepository) ListTeam(ctx context.Context, month, year int) ([]*domain.TeamVacation, error) {
	if m.ListTeamFn != nil {
		return m.ListTeamFn(ctx, month, year)
	}
	return nil, nil
}

func (m *MockVacationRepository) UpdateStatus(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, status, reviewedBy, rejectionReason)
	}
	return nil
}

func (m *MockVacationRepository) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
	if m.UpdateStatusTxFn != nil {
		return m.UpdateStatusTxFn(ctx, tx, id, status, reviewedBy, rejectionReason)
	}
	return nil
}

func (m *MockVacationRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockVacationRepository) HasOverlap(ctx context.Context, userID, startDate, endDate string) (bool, error) {
	if m.HasOverlapFn != nil {
		return m.HasOverlapFn(ctx, userID, startDate, endDate)
	}
	return false, nil
}

func (m *MockVacationRepository) GetMonthlyStats(ctx context.Context, year, month int) (*repository.MonthlyStats, error) {
	if m.GetMonthlyStatsFn != nil {
		return m.GetMonthlyStatsFn(ctx, year, month)
	}
	return &repository.MonthlyStats{}, nil
}

// MockSettingsRepository is a mock implementation of repository.SettingsRepository.
type MockSettingsRepository struct {
	GetFn                    func(ctx context.Context) (*domain.Settings, error)
	UpdateFn                 func(ctx context.Context, settings *domain.Settings) error
	UpdateLastNewsletterSentFn func(ctx context.Context, sentAt time.Time) error
}

func (m *MockSettingsRepository) Get(ctx context.Context) (*domain.Settings, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx)
	}
	defaults := domain.DefaultSettings()
	return &defaults, nil
}

func (m *MockSettingsRepository) Update(ctx context.Context, settings *domain.Settings) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, settings)
	}
	return nil
}

func (m *MockSettingsRepository) UpdateLastNewsletterSent(ctx context.Context, sentAt time.Time) error {
	if m.UpdateLastNewsletterSentFn != nil {
		return m.UpdateLastNewsletterSentFn(ctx, sentAt)
	}
	return nil
}

// MockTransactor is a mock implementation of repository.Transactor.
type MockTransactor struct {
	TransactionFn func(fn func(tx *sql.Tx) error) error
}

func (m *MockTransactor) Transaction(fn func(tx *sql.Tx) error) error {
	if m.TransactionFn != nil {
		return m.TransactionFn(fn)
	}
	// Default: execute the function with a nil tx (for simple tests)
	return fn(nil)
}
