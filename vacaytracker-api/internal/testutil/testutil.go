package testutil

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/repository/sqlite"
)

// SetupTestDB creates a temp-file SQLite database with migrations applied.
// Returns the DB and a cleanup function.
func SetupTestDB(t *testing.T) *sqlite.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sqlite.New(dbPath)
	require.NoError(t, err)

	// Find migrations directory relative to project root
	migrationsDir := findMigrationsDir(t)
	err = db.RunMigrations(migrationsDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

// findMigrationsDir locates the migrations directory by walking up from this file's location.
func findMigrationsDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "failed to get caller info")

	// Walk up from internal/testutil/ to project root
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	migrationsDir := filepath.Join(projectRoot, "migrations")
	return migrationsDir
}

// SetupTestRouter creates a Gin router in test mode.
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// CreateTestUser creates a user in the database and returns it.
func CreateTestUser(t *testing.T, repo *sqlite.UserRepository, id, email, name string, role domain.Role, balance int) *domain.User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	require.NoError(t, err)

	user := &domain.User{
		ID:               id,
		Email:            email,
		PasswordHash:     string(hash),
		Name:             name,
		Role:             role,
		VacationBalance:  balance,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	return user
}

// CreateTestVacation creates a vacation request in the database.
func CreateTestVacation(t *testing.T, repo *sqlite.VacationRepository, id, userID, startDate, endDate string, totalDays int, status domain.VacationStatus) *domain.VacationRequest {
	t.Helper()

	req := &domain.VacationRequest{
		ID:        id,
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
		TotalDays: totalDays,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), req)
	require.NoError(t, err)

	return req
}

// SetAuthContext sets authentication context values on a Gin context.
func SetAuthContext(c *gin.Context, userID, email, name string, role domain.Role) {
	c.Set("userID", userID)
	c.Set("email", email)
	c.Set("name", name)
	c.Set("role", string(role))
}
