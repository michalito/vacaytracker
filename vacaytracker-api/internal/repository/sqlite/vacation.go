package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vacaytracker-api/internal/domain"
)

// VacationRepository handles vacation request database operations
type VacationRepository struct {
	db *DB
}

// NewVacationRepository creates a new VacationRepository
func NewVacationRepository(db *DB) *VacationRepository {
	return &VacationRepository{db: db}
}

// GetDB returns the underlying database connection for transaction support
func (r *VacationRepository) GetDB() *DB {
	return r.db
}

// Create creates a new vacation request
func (r *VacationRepository) Create(ctx context.Context, req *domain.VacationRequest) error {
	query := `
		INSERT INTO vacation_requests (id, user_id, start_date, end_date, total_days, reason, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		req.ID,
		req.UserID,
		req.StartDate,
		req.EndDate,
		req.TotalDays,
		req.Reason,
		req.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create vacation request: %w", err)
	}
	return nil
}

// CreateTx creates a new vacation request within a transaction
func (r *VacationRepository) CreateTx(ctx context.Context, tx *sql.Tx, req *domain.VacationRequest) error {
	query := `
		INSERT INTO vacation_requests (id, user_id, start_date, end_date, total_days, reason, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tx.ExecContext(ctx, query,
		req.ID,
		req.UserID,
		req.StartDate,
		req.EndDate,
		req.TotalDays,
		req.Reason,
		req.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to create vacation request: %w", err)
	}
	return nil
}

// GetByID retrieves a vacation request by ID with user info
func (r *VacationRepository) GetByID(ctx context.Context, id string) (*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.id = ?
	`
	return r.scanRequest(r.db.QueryRowContext(ctx, query, id))
}

// ListByUser retrieves vacation requests for a specific user
func (r *VacationRepository) ListByUser(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.user_id = ?
	`
	args := []interface{}{userID}

	if status != nil {
		query += " AND vr.status = ?"
		args = append(args, *status)
	}

	if year != nil {
		query += " AND strftime('%Y', vr.start_date) = ?"
		args = append(args, fmt.Sprintf("%d", *year))
	}

	query += " ORDER BY vr.created_at DESC"

	return r.queryRequests(ctx, query, args...)
}

// ListPending retrieves all pending vacation requests
func (r *VacationRepository) ListPending(ctx context.Context) ([]*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.status = 'pending'
		ORDER BY vr.created_at ASC
	`
	return r.queryRequests(ctx, query)
}

// ListTeam retrieves approved vacations for team calendar view
func (r *VacationRepository) ListTeam(ctx context.Context, month, year int) ([]*domain.TeamVacation, error) {
	// Get start and end of month
	startOfMonth := fmt.Sprintf("%d-%02d-01", year, month)
	endOfMonth := fmt.Sprintf("%d-%02d-31", year, month)

	query := `
		SELECT vr.id, vr.user_id, u.name, vr.start_date, vr.end_date, vr.total_days
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.status = 'approved'
		AND (
			(vr.start_date >= ? AND vr.start_date <= ?)
			OR (vr.end_date >= ? AND vr.end_date <= ?)
			OR (vr.start_date <= ? AND vr.end_date >= ?)
		)
		ORDER BY vr.start_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query,
		startOfMonth, endOfMonth,
		startOfMonth, endOfMonth,
		startOfMonth, endOfMonth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list team vacations: %w", err)
	}
	defer rows.Close()

	var vacations []*domain.TeamVacation
	for rows.Next() {
		var v domain.TeamVacation
		if err := rows.Scan(&v.ID, &v.UserID, &v.UserName, &v.StartDate, &v.EndDate, &v.TotalDays); err != nil {
			return nil, fmt.Errorf("failed to scan team vacation: %w", err)
		}
		vacations = append(vacations, &v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating team vacations: %w", err)
	}

	return vacations, nil
}

// UpdateStatus updates the status of a vacation request
func (r *VacationRepository) UpdateStatus(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE vacation_requests
		SET status = ?, reviewed_by = ?, reviewed_at = ?, rejection_reason = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, status, reviewedBy, now, rejectionReason, id)
	if err != nil {
		return fmt.Errorf("failed to update vacation status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vacation request not found")
	}
	return nil
}

// UpdateStatusTx updates the status of a vacation request within a transaction
func (r *VacationRepository) UpdateStatusTx(ctx context.Context, tx *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE vacation_requests
		SET status = ?, reviewed_by = ?, reviewed_at = ?, rejection_reason = ?
		WHERE id = ?
	`
	result, err := tx.ExecContext(ctx, query, status, reviewedBy, now, rejectionReason, id)
	if err != nil {
		return fmt.Errorf("failed to update vacation status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vacation request not found")
	}
	return nil
}

// Delete deletes a vacation request
func (r *VacationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM vacation_requests WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete vacation request: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vacation request not found")
	}
	return nil
}

// MonthlyStats holds aggregated vacation request statistics for a specific month
type MonthlyStats struct {
	TotalSubmitted int
	TotalApproved  int
	TotalRejected  int
	TotalPending   int
	TotalDaysUsed  int
}

// GetMonthlyStats returns aggregated statistics for vacation requests in a specific month
func (r *VacationRepository) GetMonthlyStats(ctx context.Context, year, month int) (*MonthlyStats, error) {
	yearStr := fmt.Sprintf("%d", year)
	monthStr := fmt.Sprintf("%02d", month)

	query := `
		SELECT
			COUNT(*) as total,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0) as approved,
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(CASE WHEN status = 'approved' THEN total_days ELSE 0 END), 0) as days_used
		FROM vacation_requests
		WHERE strftime('%Y', created_at) = ? AND strftime('%m', created_at) = ?
	`

	var stats MonthlyStats
	err := r.db.QueryRowContext(ctx, query, yearStr, monthStr).Scan(
		&stats.TotalSubmitted,
		&stats.TotalApproved,
		&stats.TotalRejected,
		&stats.TotalPending,
		&stats.TotalDaysUsed,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly stats: %w", err)
	}

	return &stats, nil
}

// HasOverlap checks if a user has any pending or approved vacation requests that overlap with the given date range
func (r *VacationRepository) HasOverlap(ctx context.Context, userID, startDate, endDate string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM vacation_requests
		WHERE user_id = ?
		AND status IN ('pending', 'approved')
		AND (
			(start_date <= ? AND end_date >= ?)
			OR (start_date <= ? AND end_date >= ?)
			OR (start_date >= ? AND end_date <= ?)
		)
	`
	var count int
	err := r.db.QueryRowContext(ctx, query,
		userID,
		endDate, startDate,
		startDate, startDate,
		startDate, endDate,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for overlapping requests: %w", err)
	}
	return count > 0, nil
}

// scanRequest scans a single row into a VacationRequest
func (r *VacationRepository) scanRequest(row *sql.Row) (*domain.VacationRequest, error) {
	var req domain.VacationRequest
	var reason, reviewedBy, rejectionReason sql.NullString
	var reviewedAt sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&req.ID,
		&req.UserID,
		&req.UserName,
		&req.UserEmail,
		&req.StartDate,
		&req.EndDate,
		&req.TotalDays,
		&reason,
		&req.Status,
		&reviewedBy,
		&reviewedAt,
		&rejectionReason,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan vacation request: %w", err)
	}

	if reason.Valid {
		req.Reason = &reason.String
	}
	if reviewedBy.Valid {
		req.ReviewedBy = &reviewedBy.String
	}
	if reviewedAt.Valid {
		t, _ := time.Parse(time.RFC3339, reviewedAt.String)
		req.ReviewedAt = &t
	}
	if rejectionReason.Valid {
		req.RejectionReason = &rejectionReason.String
	}
	req.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	req.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &req, nil
}

// queryRequests executes a query and returns multiple VacationRequests
func (r *VacationRepository) queryRequests(ctx context.Context, query string, args ...interface{}) ([]*domain.VacationRequest, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query vacation requests: %w", err)
	}
	defer rows.Close()

	var requests []*domain.VacationRequest
	for rows.Next() {
		var req domain.VacationRequest
		var reason, reviewedBy, rejectionReason sql.NullString
		var reviewedAt sql.NullString
		var createdAt, updatedAt string

		err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.UserName,
			&req.UserEmail,
			&req.StartDate,
			&req.EndDate,
			&req.TotalDays,
			&reason,
			&req.Status,
			&reviewedBy,
			&reviewedAt,
			&rejectionReason,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vacation request row: %w", err)
		}

		if reason.Valid {
			req.Reason = &reason.String
		}
		if reviewedBy.Valid {
			req.ReviewedBy = &reviewedBy.String
		}
		if reviewedAt.Valid {
			t, _ := time.Parse(time.RFC3339, reviewedAt.String)
			req.ReviewedAt = &t
		}
		if rejectionReason.Valid {
			req.RejectionReason = &rejectionReason.String
		}
		req.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		req.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		requests = append(requests, &req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating vacation requests: %w", err)
	}

	return requests, nil
}
