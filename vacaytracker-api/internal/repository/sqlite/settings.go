package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vacaytracker-api/internal/domain"
)

// SettingsRepository handles settings database operations
type SettingsRepository struct {
	db *DB
}

// NewSettingsRepository creates a new SettingsRepository
func NewSettingsRepository(db *DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// Get retrieves the application settings
func (r *SettingsRepository) Get(ctx context.Context) (*domain.Settings, error) {
	query := `
		SELECT id, weekend_policy, newsletter, default_vacation_days, vacation_reset_month, updated_at
		FROM settings
		WHERE id = 'settings'
	`

	var settings domain.Settings
	var weekendPolicyJSON, newsletterJSON string
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query).Scan(
		&settings.ID,
		&weekendPolicyJSON,
		&newsletterJSON,
		&settings.DefaultVacationDays,
		&settings.VacationResetMonth,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		// Return default settings if none exist
		defaults := domain.DefaultSettings()
		return &defaults, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	settings.WeekendPolicy, _ = domain.ParseWeekendPolicy(weekendPolicyJSON)
	settings.Newsletter, _ = domain.ParseNewsletterConfig(newsletterJSON)
	settings.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &settings, nil
}

// Update updates the application settings (upsert)
func (r *SettingsRepository) Update(ctx context.Context, settings *domain.Settings) error {
	weekendPolicyJSON, err := settings.WeekendPolicy.ToJSONString()
	if err != nil {
		return fmt.Errorf("failed to serialize weekend policy: %w", err)
	}

	newsletterJSON, err := settings.Newsletter.ToJSONString()
	if err != nil {
		return fmt.Errorf("failed to serialize newsletter config: %w", err)
	}

	query := `
		INSERT INTO settings (id, weekend_policy, newsletter, default_vacation_days, vacation_reset_month)
		VALUES ('settings', ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			weekend_policy = excluded.weekend_policy,
			newsletter = excluded.newsletter,
			default_vacation_days = excluded.default_vacation_days,
			vacation_reset_month = excluded.vacation_reset_month
	`

	_, err = r.db.ExecContext(ctx, query,
		weekendPolicyJSON,
		newsletterJSON,
		settings.DefaultVacationDays,
		settings.VacationResetMonth,
	)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}
	return nil
}

// UpdateLastNewsletterSent updates only the newsletter lastSentAt timestamp
func (r *SettingsRepository) UpdateLastNewsletterSent(ctx context.Context, sentAt time.Time) error {
	// Get current settings first
	settings, err := r.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get settings for newsletter update: %w", err)
	}

	// Update the lastSentAt field
	settings.Newsletter.LastSentAt = &sentAt

	// Save using the existing Update method
	return r.Update(ctx, settings)
}
