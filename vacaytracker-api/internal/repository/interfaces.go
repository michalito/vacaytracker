package repository

import (
	"context"
	"database/sql"
	"time"

	"vacaytracker-api/internal/domain"
)

// Transactor provides database transaction support
type Transactor interface {
	Transaction(fn func(tx *sql.Tx) error) error
}

// UserRepository defines user data access operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error)
	GetByRole(ctx context.Context, role domain.Role) ([]*domain.User, error)
	CountByRole(ctx context.Context, role domain.Role) (int, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	UpdateEmailPreferences(ctx context.Context, id string, prefs domain.EmailPreferences) error
	UpdateVacationBalance(ctx context.Context, id string, balance int) error
	UpdateVacationBalanceTx(ctx context.Context, tx *sql.Tx, id string, balance int) error
	Delete(ctx context.Context, id string) error
	EmailExists(ctx context.Context, email string) (bool, error)
	EmailExistsExcluding(ctx context.Context, email, excludeID string) (bool, error)
	GetNewsletterRecipients(ctx context.Context) ([]*domain.User, error)
	GetLowBalanceUsers(ctx context.Context, threshold int) ([]*domain.User, error)
	UpdateAllBalances(ctx context.Context, balance int) (int64, error)
}

// VacationRepository defines vacation request data access operations
type VacationRepository interface {
	Create(ctx context.Context, req *domain.VacationRequest) error
	CreateTx(ctx context.Context, tx *sql.Tx, req *domain.VacationRequest) error
	GetByID(ctx context.Context, id string) (*domain.VacationRequest, error)
	ListByUser(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error)
	ListPending(ctx context.Context) ([]*domain.VacationRequest, error)
	ListTeam(ctx context.Context, month, year int) ([]*domain.TeamVacation, error)
	UpdateStatus(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error
	Delete(ctx context.Context, id string) error
	HasOverlap(ctx context.Context, userID, startDate, endDate string) (bool, error)
	GetMonthlyStats(ctx context.Context, year, month int) (*MonthlyStats, error)
}

// SettingsRepository defines settings data access operations
type SettingsRepository interface {
	Get(ctx context.Context) (*domain.Settings, error)
	Update(ctx context.Context, settings *domain.Settings) error
	UpdateLastNewsletterSent(ctx context.Context, sentAt time.Time) error
}

// MonthlyStats holds aggregated vacation request statistics for a specific month
type MonthlyStats struct {
	TotalSubmitted int
	TotalApproved  int
	TotalRejected  int
	TotalPending   int
	TotalDaysUsed  int
}
