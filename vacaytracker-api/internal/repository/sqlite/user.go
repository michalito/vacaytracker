package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vacaytracker-api/internal/domain"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// generateUserID generates a unique user ID with prefix
func generateUserID() string {
	return "usr_" + uuid.New().String()[:8]
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = generateUserID()
	}

	emailPrefsJSON, err := user.EmailPreferences.ToJSONString()
	if err != nil {
		return fmt.Errorf("failed to serialize email preferences: %w", err)
	}

	query := `
		INSERT INTO users (id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	_, err = r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		string(user.Role),
		user.VacationBalance,
		user.StartDate,
		emailPrefsJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

// GetAll retrieves all users with optional filtering and pagination
func (r *UserRepository) GetAll(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
	// Build query with filters
	baseQuery := "FROM users WHERE 1=1"
	args := []interface{}{}

	if role != nil {
		baseQuery += " AND role = ?"
		args = append(args, string(*role))
	}

	if search != "" {
		baseQuery += " AND (name LIKE ? OR email LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	selectQuery := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
	` + baseQuery + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByRole retrieves all users with a specific role
func (r *UserRepository) GetByRole(ctx context.Context, role domain.Role) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users
		WHERE role = ?
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, string(role))
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// CountByRole counts users with a specific role
func (r *UserRepository) CountByRole(ctx context.Context, role domain.Role) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = ?`

	var count int
	if err := r.db.QueryRowContext(ctx, query, string(role)).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}

	return count, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	emailPrefsJSON, err := user.EmailPreferences.ToJSONString()
	if err != nil {
		return fmt.Errorf("failed to serialize email preferences: %w", err)
	}

	query := `
		UPDATE users
		SET email = ?, name = ?, role = ?, vacation_balance = ?, start_date = ?, email_preferences = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Name,
		string(user.Role),
		user.VacationBalance,
		user.StartDate,
		emailPrefsJSON,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdatePassword updates a user's password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateEmailPreferences updates a user's email preferences
func (r *UserRepository) UpdateEmailPreferences(ctx context.Context, id string, prefs domain.EmailPreferences) error {
	prefsJSON, err := prefs.ToJSONString()
	if err != nil {
		return fmt.Errorf("failed to serialize email preferences: %w", err)
	}

	query := `UPDATE users SET email_preferences = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, prefsJSON, id)
	if err != nil {
		return fmt.Errorf("failed to update email preferences: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateVacationBalance updates a user's vacation balance
func (r *UserRepository) UpdateVacationBalance(ctx context.Context, id string, balance int) error {
	query := `UPDATE users SET vacation_balance = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, balance, id)
	if err != nil {
		return fmt.Errorf("failed to update vacation balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateVacationBalanceTx updates a user's vacation balance within a transaction
func (r *UserRepository) UpdateVacationBalanceTx(ctx context.Context, tx *sql.Tx, id string, balance int) error {
	query := `UPDATE users SET vacation_balance = ? WHERE id = ?`

	result, err := tx.ExecContext(ctx, query, balance, id)
	if err != nil {
		return fmt.Errorf("failed to update vacation balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// EmailExists checks if an email address is already in use
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// EmailExistsExcluding checks if an email is in use by another user
func (r *UserRepository) EmailExistsExcluding(ctx context.Context, email, excludeID string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND id != ?`

	var count int
	if err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

// GetNewsletterRecipients returns users who have weeklyDigest email preference enabled
func (r *UserRepository) GetNewsletterRecipients(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users
		WHERE json_extract(email_preferences, '$.weeklyDigest') = 1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query newsletter recipients: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// GetLowBalanceUsers returns users with vacation balance at or below the threshold
func (r *UserRepository) GetLowBalanceUsers(ctx context.Context, threshold int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users
		WHERE vacation_balance <= ? AND role = 'employee'
		ORDER BY vacation_balance ASC
	`

	rows, err := r.db.QueryContext(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to query low balance users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// UpdateAllBalances resets vacation balance for all employees to the specified value
func (r *UserRepository) UpdateAllBalances(ctx context.Context, balance int) (int64, error) {
	query := `UPDATE users SET vacation_balance = ? WHERE role = 'employee'`

	result, err := r.db.ExecContext(ctx, query, balance)
	if err != nil {
		return 0, fmt.Errorf("failed to update all balances: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// scanUser scans a single user row
func (r *UserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var role string
	var startDate sql.NullString
	var emailPrefsJSON string
	var createdAt, updatedAt string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&role,
		&user.VacationBalance,
		&startDate,
		&emailPrefsJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	user.Role = domain.Role(role)

	if startDate.Valid {
		user.StartDate = &startDate.String
	}

	user.EmailPreferences, _ = domain.ParseEmailPreferences(emailPrefsJSON)

	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &user, nil
}

// scanUsers scans multiple user rows
func (r *UserRepository) scanUsers(rows *sql.Rows) ([]*domain.User, error) {
	var users []*domain.User

	for rows.Next() {
		var user domain.User
		var role string
		var startDate sql.NullString
		var emailPrefsJSON string
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&role,
			&user.VacationBalance,
			&startDate,
			&emailPrefsJSON,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.Role = domain.Role(role)

		if startDate.Valid {
			user.StartDate = &startDate.String
		}

		user.EmailPreferences, _ = domain.ParseEmailPreferences(emailPrefsJSON)

		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
