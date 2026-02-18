package sqlite_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/repository/sqlite"
	"vacaytracker-api/internal/testutil"
)

// findMigrationsDir locates the migrations directory by walking up from this
// test file's location to the project root.
func findMigrationsDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get caller info")
	// Walk up from internal/repository/sqlite/ to project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	return filepath.Join(projectRoot, "migrations")
}

// =============================================================================
// Settings Repository Tests
// =============================================================================

func TestSettingsGet_DefaultsWhenEmpty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)
	require.NotNil(t, settings)

	// The migration inserts a default settings row, so Get should return
	// values matching the schema defaults (same as domain.DefaultSettings).
	assert.Equal(t, "settings", settings.ID)
	assert.Equal(t, 25, settings.DefaultVacationDays)
	assert.Equal(t, 1, settings.VacationResetMonth)
	assert.True(t, settings.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{0, 6}, settings.WeekendPolicy.ExcludedDays)
	assert.False(t, settings.Newsletter.Enabled)
	assert.Equal(t, "monthly", settings.Newsletter.Frequency)
	assert.Equal(t, 1, settings.Newsletter.DayOfMonth)
	assert.Nil(t, settings.Newsletter.LastSentAt)
}

func TestSettingsUpdateAndGet_Roundtrip(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// Build custom settings
	updated := &domain.Settings{
		ID: "settings",
		WeekendPolicy: domain.WeekendPolicy{
			ExcludeWeekends: false,
			ExcludedDays:    []int{},
		},
		Newsletter: domain.NewsletterConfig{
			Enabled:    true,
			Frequency:  "weekly",
			DayOfMonth: 15,
			LastSentAt: nil,
		},
		DefaultVacationDays: 30,
		VacationResetMonth:  6,
	}

	err := repo.Update(ctx, updated)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, "settings", got.ID)
	assert.Equal(t, 30, got.DefaultVacationDays)
	assert.Equal(t, 6, got.VacationResetMonth)
	assert.False(t, got.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{}, got.WeekendPolicy.ExcludedDays)
	assert.True(t, got.Newsletter.Enabled)
	assert.Equal(t, "weekly", got.Newsletter.Frequency)
	assert.Equal(t, 15, got.Newsletter.DayOfMonth)
	assert.Nil(t, got.Newsletter.LastSentAt)
}

func TestSettingsUpdate_CustomWeekendPolicy(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// Start from current settings
	settings, err := repo.Get(ctx)
	require.NoError(t, err)

	// Disable weekend exclusion
	settings.WeekendPolicy = domain.WeekendPolicy{
		ExcludeWeekends: false,
		ExcludedDays:    []int{},
	}

	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.False(t, got.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{}, got.WeekendPolicy.ExcludedDays)
}

func TestSettingsUpdate_CustomWeekendPolicy_FridaySaturday(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)

	// Some regions use Friday/Saturday as weekend
	settings.WeekendPolicy = domain.WeekendPolicy{
		ExcludeWeekends: true,
		ExcludedDays:    []int{5, 6}, // Friday = 5, Saturday = 6
	}

	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.True(t, got.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{5, 6}, got.WeekendPolicy.ExcludedDays)
}

func TestSettingsUpdate_NewsletterConfig(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)

	settings.Newsletter = domain.NewsletterConfig{
		Enabled:    true,
		Frequency:  "weekly",
		DayOfMonth: 10,
		LastSentAt: nil,
	}

	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.True(t, got.Newsletter.Enabled)
	assert.Equal(t, "weekly", got.Newsletter.Frequency)
	assert.Equal(t, 10, got.Newsletter.DayOfMonth)
	assert.Nil(t, got.Newsletter.LastSentAt)
}

func TestSettingsUpdate_DefaultVacationDays(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)

	settings.DefaultVacationDays = 30

	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.Equal(t, 30, got.DefaultVacationDays)
}

func TestSettingsUpdate_VacationResetMonth(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)

	settings.VacationResetMonth = 6 // June

	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.Equal(t, 6, got.VacationResetMonth)
}

func TestSettingsUpdateLastNewsletterSent(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// Verify LastSentAt is initially nil
	before, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Nil(t, before.Newsletter.LastSentAt)

	// Set the newsletter sent timestamp
	sentAt := time.Date(2026, 2, 15, 10, 30, 0, 0, time.UTC)
	err = repo.UpdateLastNewsletterSent(ctx, sentAt)
	require.NoError(t, err)

	// Verify the timestamp was persisted
	got, err := repo.Get(ctx)
	require.NoError(t, err)
	require.NotNil(t, got.Newsletter.LastSentAt)
	assert.True(t, got.Newsletter.LastSentAt.Equal(sentAt),
		"expected LastSentAt %v, got %v", sentAt, *got.Newsletter.LastSentAt)
}

func TestSettingsUpdateLastNewsletterSent_PreservesOtherFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// First set custom newsletter config
	settings, err := repo.Get(ctx)
	require.NoError(t, err)
	settings.Newsletter = domain.NewsletterConfig{
		Enabled:    true,
		Frequency:  "weekly",
		DayOfMonth: 20,
		LastSentAt: nil,
	}
	settings.DefaultVacationDays = 28
	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	// Now update only LastSentAt
	sentAt := time.Date(2026, 1, 20, 8, 0, 0, 0, time.UTC)
	err = repo.UpdateLastNewsletterSent(ctx, sentAt)
	require.NoError(t, err)

	// Verify all fields are preserved
	got, err := repo.Get(ctx)
	require.NoError(t, err)

	assert.True(t, got.Newsletter.Enabled)
	assert.Equal(t, "weekly", got.Newsletter.Frequency)
	assert.Equal(t, 20, got.Newsletter.DayOfMonth)
	require.NotNil(t, got.Newsletter.LastSentAt)
	assert.True(t, got.Newsletter.LastSentAt.Equal(sentAt))
	assert.Equal(t, 28, got.DefaultVacationDays)
}

func TestSettingsUpdate_IsUpsert(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// The migration already inserts a default row.
	// Delete it so we can test the INSERT path of the upsert.
	_, err := db.ExecContext(ctx, "DELETE FROM settings WHERE id = 'settings'")
	require.NoError(t, err)

	// Verify no row exists
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM settings").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// First Update should INSERT
	first := &domain.Settings{
		ID:                  "settings",
		WeekendPolicy:       domain.DefaultWeekendPolicy(),
		Newsletter:          domain.DefaultNewsletterConfig(),
		DefaultVacationDays: 20,
		VacationResetMonth:  3,
	}
	err = repo.Update(ctx, first)
	require.NoError(t, err)

	got, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, 20, got.DefaultVacationDays)
	assert.Equal(t, 3, got.VacationResetMonth)

	// Second Update should UPDATE (upsert conflict path)
	second := &domain.Settings{
		ID:                  "settings",
		WeekendPolicy:       domain.DefaultWeekendPolicy(),
		Newsletter:          domain.DefaultNewsletterConfig(),
		DefaultVacationDays: 35,
		VacationResetMonth:  9,
	}
	err = repo.Update(ctx, second)
	require.NoError(t, err)

	got2, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, 35, got2.DefaultVacationDays)
	assert.Equal(t, 9, got2.VacationResetMonth)

	// Verify there is still exactly one row
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM settings").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestSettingsMultipleUpdates_PreserveData(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := sqlite.NewSettingsRepository(db)
	ctx := context.Background()

	// Set custom vacation days
	settings, err := repo.Get(ctx)
	require.NoError(t, err)
	settings.DefaultVacationDays = 30
	err = repo.Update(ctx, settings)
	require.NoError(t, err)

	// Now update only the weekend policy
	settings2, err := repo.Get(ctx)
	require.NoError(t, err)
	settings2.WeekendPolicy = domain.WeekendPolicy{
		ExcludeWeekends: false,
		ExcludedDays:    []int{},
	}
	err = repo.Update(ctx, settings2)
	require.NoError(t, err)

	// Verify vacation days were preserved
	got, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, 30, got.DefaultVacationDays)
	assert.False(t, got.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{}, got.WeekendPolicy.ExcludedDays)

	// Now update only newsletter config
	got.Newsletter = domain.NewsletterConfig{
		Enabled:    true,
		Frequency:  "monthly",
		DayOfMonth: 5,
		LastSentAt: nil,
	}
	err = repo.Update(ctx, got)
	require.NoError(t, err)

	// Verify all previous changes are still intact
	final, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, 30, final.DefaultVacationDays)
	assert.False(t, final.WeekendPolicy.ExcludeWeekends)
	assert.Equal(t, []int{}, final.WeekendPolicy.ExcludedDays)
	assert.True(t, final.Newsletter.Enabled)
	assert.Equal(t, "monthly", final.Newsletter.Frequency)
	assert.Equal(t, 5, final.Newsletter.DayOfMonth)
}

// =============================================================================
// Migration Tests
// =============================================================================

func TestMigration_FreshRunSuccessful(t *testing.T) {
	// SetupTestDB creates a fresh database and runs migrations.
	// If this does not panic or fail, migrations ran successfully.
	db := testutil.SetupTestDB(t)
	require.NotNil(t, db)
}

func TestMigration_IdempotentRerun(t *testing.T) {
	// SetupTestDB already runs migrations once.
	db := testutil.SetupTestDB(t)

	// Run migrations a second time; already-applied migrations should be skipped.
	migrationsDir := findMigrationsDir(t)
	err := db.RunMigrations(migrationsDir)
	require.NoError(t, err, "running migrations a second time should not error")

	// Verify the database is still functional
	repo := sqlite.NewSettingsRepository(db)
	settings, err := repo.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 25, settings.DefaultVacationDays)
}

func TestMigration_TablesExist(t *testing.T) {
	db := testutil.SetupTestDB(t)

	expectedTables := []string{
		"users",
		"vacation_requests",
		"settings",
		"schema_migrations",
	}

	for _, table := range expectedTables {
		t.Run(table, func(t *testing.T) {
			var name string
			err := db.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
				table,
			).Scan(&name)
			require.NoError(t, err, "table %q should exist", table)
			assert.Equal(t, table, name)
		})
	}
}

func TestMigration_IndexesExist(t *testing.T) {
	db := testutil.SetupTestDB(t)

	expectedIndexes := []string{
		"idx_vacation_requests_user_id",
		"idx_vacation_requests_status",
		"idx_users_role",
		"idx_users_email",
	}

	for _, idx := range expectedIndexes {
		t.Run(idx, func(t *testing.T) {
			var name string
			err := db.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='index' AND name=?",
				idx,
			).Scan(&name)
			require.NoError(t, err, "index %q should exist", idx)
			assert.Equal(t, idx, name)
		})
	}
}

func TestMigration_TriggersExist(t *testing.T) {
	db := testutil.SetupTestDB(t)

	expectedTriggers := []string{
		"users_updated_at",
		"vacation_requests_updated_at",
		"settings_updated_at",
	}

	for _, trigger := range expectedTriggers {
		t.Run(trigger, func(t *testing.T) {
			var name string
			err := db.QueryRow(
				"SELECT name FROM sqlite_master WHERE type='trigger' AND name=?",
				trigger,
			).Scan(&name)
			require.NoError(t, err, "trigger %q should exist", trigger)
			assert.Equal(t, trigger, name)
		})
	}
}

func TestMigration_DefaultSettingsRowInserted(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// The migration includes INSERT OR IGNORE for default settings
	var id string
	var defaultDays int
	var resetMonth int
	err := db.QueryRow(
		"SELECT id, default_vacation_days, vacation_reset_month FROM settings WHERE id = 'settings'",
	).Scan(&id, &defaultDays, &resetMonth)
	require.NoError(t, err)

	assert.Equal(t, "settings", id)
	assert.Equal(t, 25, defaultDays)
	assert.Equal(t, 1, resetMonth)
}

func TestMigration_SchemaMigrationsTracked(t *testing.T) {
	db := testutil.SetupTestDB(t)

	// After running migrations, schema_migrations should have at least one entry
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1, "at least one migration should be recorded")

	// Verify the first migration version is recorded
	var version string
	err = db.QueryRow("SELECT version FROM schema_migrations WHERE version = '001'").Scan(&version)
	require.NoError(t, err)
	assert.Equal(t, "001", version)
}
