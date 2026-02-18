package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/config"
	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/handler"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

// ---------------------------------------------------------------------------
// Test setup helpers
// ---------------------------------------------------------------------------

type adminTestDeps struct {
	userRepo     *testutil.MockUserRepository
	vacRepo      *testutil.MockVacationRepository
	settingsRepo *testutil.MockSettingsRepository
	transactor   *testutil.MockTransactor
	handler      *handler.AdminHandler
	router       *gin.Engine
}

func setupAdminTest(t *testing.T) *adminTestDeps {
	t.Helper()
	gin.SetMode(gin.TestMode)

	userRepo := &testutil.MockUserRepository{}
	vacRepo := &testutil.MockVacationRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	cfg := &config.Config{
		JWTSecret: "test-secret-key-that-is-at-least-32-chars",
		AppURL:    "http://localhost:3000",
	}

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, authService)
	vacationService := service.NewVacationService(vacRepo, userRepo, settingsRepo, transactor)
	emailService := service.NewEmailService(cfg)
	newsletterService := service.NewNewsletterService(cfg, userRepo, vacRepo, settingsRepo, emailService)

	h := handler.NewAdminHandler(cfg, userService, userRepo, vacationService, vacRepo, settingsRepo, emailService, newsletterService)

	r := gin.New()
	admin := r.Group("/api/admin")
	admin.Use(func(c *gin.Context) {
		testutil.SetAuthContext(c, "admin-1", "admin@test.com", "Admin", domain.RoleAdmin)
		c.Next()
	})
	{
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.GET("/users/:id", h.GetUser)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.DELETE("/users/:id", h.DeleteUser)
		admin.PUT("/users/:id/balance", h.UpdateBalance)
		admin.POST("/users/reset-balances", h.ResetBalances)
		admin.GET("/vacation/pending", h.ListPending)
		admin.PUT("/vacation/:id/review", h.Review)
		admin.GET("/settings", h.GetSettings)
		admin.PUT("/settings", h.UpdateSettings)
	}

	return &adminTestDeps{
		userRepo:     userRepo,
		vacRepo:      vacRepo,
		settingsRepo: settingsRepo,
		transactor:   transactor,
		handler:      h,
		router:       r,
	}
}

// sampleUser builds a domain.User for testing.
func sampleUser(id, email, name string, role domain.Role, balance int) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:               id,
		Email:            email,
		Name:             name,
		Role:             role,
		VacationBalance:  balance,
		PasswordHash:     "$2a$04$fakehashfortest",
		EmailPreferences: domain.DefaultEmailPreferences(),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// sampleVacation builds a domain.VacationRequest for testing.
func sampleVacation(id, userID string, status domain.VacationStatus, totalDays int) *domain.VacationRequest {
	now := time.Now()
	return &domain.VacationRequest{
		ID:        id,
		UserID:    userID,
		UserName:  "Test Employee",
		UserEmail: "employee@test.com",
		StartDate: "2026-03-01",
		EndDate:   "2026-03-05",
		TotalDays: totalDays,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ===================================================================
// ListUsers tests
// ===================================================================

func TestAdminListUsers_Success(t *testing.T) {
	deps := setupAdminTest(t)

	users := []*domain.User{
		sampleUser("u1", "alice@test.com", "Alice", domain.RoleEmployee, 20),
		sampleUser("u2", "bob@test.com", "Bob", domain.RoleAdmin, 25),
	}

	deps.userRepo.GetAllFn = func(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
		return users, 2, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Users, 2)
	require.NotNil(t, resp.Pagination)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 20, resp.Pagination.Limit)
	assert.Equal(t, 2, resp.Pagination.Total)
	assert.Equal(t, 1, resp.Pagination.TotalPages)
	assert.Equal(t, "u1", resp.Users[0].ID)
	assert.Equal(t, "u2", resp.Users[1].ID)
}

func TestAdminListUsers_WithRoleFilter(t *testing.T) {
	deps := setupAdminTest(t)

	var capturedRole *domain.Role
	deps.userRepo.GetAllFn = func(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
		capturedRole = role
		return []*domain.User{sampleUser("u1", "alice@test.com", "Alice", domain.RoleEmployee, 20)}, 1, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users?role=employee", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	require.NotNil(t, capturedRole)
	assert.Equal(t, domain.RoleEmployee, *capturedRole)

	var resp dto.UserListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Users, 1)
}

func TestAdminListUsers_InvalidRole(t *testing.T) {
	deps := setupAdminTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users?role=invalid", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Invalid role")
}

// ===================================================================
// CreateUser tests
// ===================================================================

func TestAdminCreateUser_Success(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.EmailExistsFn = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	deps.userRepo.CreateFn = func(ctx context.Context, user *domain.User) error {
		return nil
	}

	body := `{"email":"new@test.com","password":"securePass1","name":"New User","role":"employee"}`
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "new@test.com", resp.Email)
	assert.Equal(t, "New User", resp.Name)
	assert.Equal(t, "employee", resp.Role)
	assert.Equal(t, 25, resp.VacationBalance) // default
}

func TestAdminCreateUser_InvalidBody(t *testing.T) {
	deps := setupAdminTest(t)

	// Missing required fields (email, password, name, role)
	body := `{"email":"notanemail"}`
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

func TestAdminCreateUser_DuplicateEmail(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.EmailExistsFn = func(ctx context.Context, email string) (bool, error) {
		return true, nil
	}

	body := `{"email":"exists@test.com","password":"securePass1","name":"Dup User","role":"employee"}`
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrAlreadyExists, resp.Code)
}

// ===================================================================
// GetUser tests
// ===================================================================

func TestAdminGetUser_Success(t *testing.T) {
	deps := setupAdminTest(t)

	user := sampleUser("user-42", "jane@test.com", "Jane Doe", domain.RoleEmployee, 18)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-42" {
			return user, nil
		}
		return nil, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/user-42", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user-42", resp.ID)
	assert.Equal(t, "jane@test.com", resp.Email)
	assert.Equal(t, "Jane Doe", resp.Name)
	assert.Equal(t, 18, resp.VacationBalance)
}

func TestAdminGetUser_NotFound(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		return nil, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/nonexistent", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

// ===================================================================
// UpdateUser tests
// ===================================================================

func TestAdminUpdateUser_Success(t *testing.T) {
	deps := setupAdminTest(t)

	user := sampleUser("user-42", "old@test.com", "Old Name", domain.RoleEmployee, 20)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-42" {
			return user, nil
		}
		return nil, nil
	}
	deps.userRepo.EmailExistsExcludingFn = func(ctx context.Context, email, excludeID string) (bool, error) {
		return false, nil
	}
	deps.userRepo.UpdateFn = func(ctx context.Context, u *domain.User) error {
		return nil
	}

	body := `{"name":"New Name","email":"new@test.com"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user-42", resp.ID)
	assert.Equal(t, "New Name", resp.Name)
	assert.Equal(t, "new@test.com", resp.Email)
}

func TestAdminUpdateUser_NotFound(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		return nil, nil
	}

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/nonexistent", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

func TestAdminUpdateUser_InvalidBody(t *testing.T) {
	deps := setupAdminTest(t)

	// Send invalid role value which fails binding:"omitempty,oneof=admin employee"
	body := `{"role":"superuser"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

// ===================================================================
// DeleteUser tests
// ===================================================================

func TestAdminDeleteUser_Success(t *testing.T) {
	deps := setupAdminTest(t)

	user := sampleUser("user-42", "target@test.com", "Target User", domain.RoleEmployee, 20)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-42" {
			return user, nil
		}
		return nil, nil
	}
	deps.userRepo.DeleteFn = func(ctx context.Context, id string) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/user-42", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "User deleted successfully", resp.Message)
}

func TestAdminDeleteUser_CannotDeleteSelf(t *testing.T) {
	deps := setupAdminTest(t)

	// The auth context sets userID to "admin-1", so deleting "admin-1" = self-delete
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/admin-1", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrForbidden, resp.Code)
	assert.Contains(t, resp.Message, "cannot delete your own account")
}

func TestAdminDeleteUser_NotFound(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		return nil, nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/nonexistent", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

// ===================================================================
// UpdateBalance tests
// ===================================================================

func TestAdminUpdateBalance_Success(t *testing.T) {
	deps := setupAdminTest(t)

	user := sampleUser("user-42", "emp@test.com", "Employee", domain.RoleEmployee, 20)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-42" {
			return user, nil
		}
		return nil, nil
	}
	deps.userRepo.UpdateVacationBalanceFn = func(ctx context.Context, id string, balance int) error {
		return nil
	}

	body := `{"vacationBalance":30}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-42/balance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user-42", resp.ID)
	assert.Equal(t, 30, resp.VacationBalance)
}

func TestAdminUpdateBalance_UserNotFound(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		return nil, nil
	}

	body := `{"vacationBalance":30}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/nonexistent/balance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

func TestAdminUpdateBalance_InvalidBody(t *testing.T) {
	deps := setupAdminTest(t)

	// Missing required vacationBalance field
	body := `{}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/user-42/balance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

// ===================================================================
// ListPending tests
// ===================================================================

func TestAdminListPending_Success(t *testing.T) {
	deps := setupAdminTest(t)

	pending := []*domain.VacationRequest{
		sampleVacation("vac-1", "user-10", domain.StatusPending, 3),
		sampleVacation("vac-2", "user-20", domain.StatusPending, 5),
	}

	deps.vacRepo.ListPendingFn = func(ctx context.Context) ([]*domain.VacationRequest, error) {
		return pending, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/vacation/pending", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Requests, 2)
	assert.Equal(t, "vac-1", resp.Requests[0].ID)
	assert.Equal(t, "pending", resp.Requests[0].Status)
	assert.Equal(t, "vac-2", resp.Requests[1].ID)
}

// ===================================================================
// Review tests
// ===================================================================

func TestAdminReview_ApproveSuccess(t *testing.T) {
	deps := setupAdminTest(t)

	vacation := sampleVacation("vac-1", "user-10", domain.StatusPending, 3)
	user := sampleUser("user-10", "emp@test.com", "Employee", domain.RoleEmployee, 20)

	approvedVacation := *vacation
	approvedVacation.Status = domain.StatusApproved

	callCount := 0
	deps.vacRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-1" {
			callCount++
			// First call: pending (for Approve check); subsequent calls: approved (after update)
			if callCount <= 1 {
				return vacation, nil
			}
			return &approvedVacation, nil
		}
		return nil, nil
	}

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-10" {
			return user, nil
		}
		return nil, nil
	}

	deps.vacRepo.UpdateStatusTxFn = func(ctx context.Context, tx *sql.Tx, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
		assert.Equal(t, "vac-1", id)
		assert.Equal(t, domain.StatusApproved, status)
		assert.Equal(t, "admin-1", reviewedBy)
		assert.Nil(t, rejectionReason)
		return nil
	}

	deps.userRepo.UpdateVacationBalanceTxFn = func(ctx context.Context, tx *sql.Tx, id string, balance int) error {
		assert.Equal(t, "user-10", id)
		assert.Equal(t, 17, balance) // 20 - 3
		return nil
	}

	body := `{"status":"approved"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/vac-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationRequestResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "vac-1", resp.ID)
	assert.Equal(t, "approved", resp.Status)
}

func TestAdminReview_RejectSuccess(t *testing.T) {
	deps := setupAdminTest(t)

	vacation := sampleVacation("vac-2", "user-10", domain.StatusPending, 5)

	rejectedVacation := *vacation
	rejectedVacation.Status = domain.StatusRejected
	reason := "Project deadline conflict"
	rejectedVacation.RejectionReason = &reason

	callCount := 0
	deps.vacRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-2" {
			callCount++
			if callCount <= 1 {
				return vacation, nil
			}
			return &rejectedVacation, nil
		}
		return nil, nil
	}

	deps.vacRepo.UpdateStatusFn = func(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
		assert.Equal(t, "vac-2", id)
		assert.Equal(t, domain.StatusRejected, status)
		assert.Equal(t, "admin-1", reviewedBy)
		require.NotNil(t, rejectionReason)
		assert.Equal(t, "Project deadline conflict", *rejectionReason)
		return nil
	}

	// For reject, the user repo GetByID is called in the async sendReviewEmail goroutine.
	// We set it up to avoid nil pointer issues but it does not affect the HTTP response.
	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		return sampleUser(id, "emp@test.com", "Employee", domain.RoleEmployee, 20), nil
	}

	body := `{"status":"rejected","reason":"Project deadline conflict"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/vac-2/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationRequestResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "vac-2", resp.ID)
	assert.Equal(t, "rejected", resp.Status)
}

func TestAdminReview_InvalidBody(t *testing.T) {
	deps := setupAdminTest(t)

	// Status must be "approved" or "rejected"
	body := `{"status":"cancelled"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/vac-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

// ===================================================================
// GetSettings tests
// ===================================================================

func TestAdminGetSettings_Success(t *testing.T) {
	deps := setupAdminTest(t)

	settings := domain.DefaultSettings()
	settings.DefaultVacationDays = 30
	settings.VacationResetMonth = 4

	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		return &settings, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/settings", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SettingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 30, resp.DefaultVacationDays)
	assert.Equal(t, 4, resp.VacationResetMonth)
	assert.True(t, resp.WeekendPolicy.ExcludeWeekends)
}

// ===================================================================
// UpdateSettings tests
// ===================================================================

func TestAdminUpdateSettings_Success(t *testing.T) {
	deps := setupAdminTest(t)

	settings := domain.DefaultSettings()

	callCount := 0
	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		callCount++
		if callCount == 1 {
			// First call: return current settings
			return &settings, nil
		}
		// Second call: return updated settings (after update)
		updated := settings
		updated.DefaultVacationDays = 30
		updated.VacationResetMonth = 6
		return &updated, nil
	}

	var updatedSettings *domain.Settings
	deps.settingsRepo.UpdateFn = func(ctx context.Context, s *domain.Settings) error {
		updatedSettings = s
		return nil
	}

	body := `{"defaultVacationDays":30,"vacationResetMonth":6}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify Update was called with the right values
	require.NotNil(t, updatedSettings)
	assert.Equal(t, 30, updatedSettings.DefaultVacationDays)
	assert.Equal(t, 6, updatedSettings.VacationResetMonth)

	var resp dto.SettingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 30, resp.DefaultVacationDays)
	assert.Equal(t, 6, resp.VacationResetMonth)
}

func TestAdminUpdateSettings_InvalidBody(t *testing.T) {
	deps := setupAdminTest(t)

	// vacationResetMonth must be 1-12
	body := `{"vacationResetMonth":13}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

// ===================================================================
// ResetBalances tests
// ===================================================================

func TestAdminResetBalances_Success(t *testing.T) {
	deps := setupAdminTest(t)

	settings := domain.DefaultSettings()
	settings.DefaultVacationDays = 25

	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		return &settings, nil
	}

	deps.userRepo.UpdateAllBalancesFn = func(ctx context.Context, balance int) (int64, error) {
		assert.Equal(t, 25, balance)
		return 10, nil
	}

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/reset-balances", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.ResetBalancesResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp.Success)
	assert.Equal(t, 10, resp.UsersUpdated)
	assert.Equal(t, 25, resp.NewBalance)
	assert.Contains(t, resp.Message, "Reset vacation balance to 25 days for 10 employees")
}

// ===================================================================
// Additional edge-case tests
// ===================================================================

func TestAdminListUsers_WithPagination(t *testing.T) {
	deps := setupAdminTest(t)

	var capturedLimit, capturedOffset int
	deps.userRepo.GetAllFn = func(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
		capturedLimit = limit
		capturedOffset = offset
		return []*domain.User{sampleUser("u1", "a@test.com", "A", domain.RoleEmployee, 20)}, 50, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users?page=3&limit=10", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 10, capturedLimit)
	assert.Equal(t, 20, capturedOffset) // (page-1) * limit = (3-1) * 10

	var resp dto.UserListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 3, resp.Pagination.Page)
	assert.Equal(t, 10, resp.Pagination.Limit)
	assert.Equal(t, 50, resp.Pagination.Total)
	assert.Equal(t, 5, resp.Pagination.TotalPages)
}

func TestAdminReview_VacationNotFound(t *testing.T) {
	deps := setupAdminTest(t)

	deps.vacRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.VacationRequest, error) {
		return nil, nil
	}

	body := `{"status":"approved"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/nonexistent/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

func TestAdminReview_AlreadyProcessed(t *testing.T) {
	deps := setupAdminTest(t)

	// Vacation is already approved, cannot be reviewed again
	vacation := sampleVacation("vac-1", "user-10", domain.StatusApproved, 3)

	deps.vacRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-1" {
			return vacation, nil
		}
		return nil, nil
	}

	body := `{"status":"rejected","reason":"too late"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/vac-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrAlreadyExists, resp.Code)
}

func TestAdminGetSettings_RepoError(t *testing.T) {
	deps := setupAdminTest(t)

	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		return nil, fmt.Errorf("database connection lost")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/settings", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInternal, resp.Code)
}

func TestAdminUpdateSettings_WeekendPolicy(t *testing.T) {
	deps := setupAdminTest(t)

	settings := domain.DefaultSettings()

	callCount := 0
	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		callCount++
		if callCount == 1 {
			return &settings, nil
		}
		updated := settings
		updated.WeekendPolicy.ExcludeWeekends = false
		return &updated, nil
	}
	deps.settingsRepo.UpdateFn = func(ctx context.Context, s *domain.Settings) error {
		return nil
	}

	body := `{"weekendPolicy":{"excludeWeekends":false}}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.SettingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.WeekendPolicy.ExcludeWeekends)
}

func TestAdminDeleteUser_LastAdmin(t *testing.T) {
	deps := setupAdminTest(t)

	// Target is a different admin, but they are the last admin.
	user := sampleUser("other-admin", "otheradmin@test.com", "Other Admin", domain.RoleAdmin, 25)

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "other-admin" {
			return user, nil
		}
		return nil, nil
	}
	deps.userRepo.CountByRoleFn = func(ctx context.Context, role domain.Role) (int, error) {
		if role == domain.RoleAdmin {
			return 1, nil // only one admin left
		}
		return 0, nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/other-admin", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrForbidden, resp.Code)
	assert.Contains(t, resp.Message, "last admin")
}

func TestAdminCreateUser_WithCustomBalance(t *testing.T) {
	deps := setupAdminTest(t)

	deps.userRepo.EmailExistsFn = func(ctx context.Context, email string) (bool, error) {
		return false, nil
	}
	deps.userRepo.CreateFn = func(ctx context.Context, user *domain.User) error {
		return nil
	}

	body := `{"email":"custom@test.com","password":"securePass1","name":"Custom Balance","role":"employee","vacationBalance":15}`
	req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 15, resp.VacationBalance)
}

func TestAdminListPending_Empty(t *testing.T) {
	deps := setupAdminTest(t)

	deps.vacRepo.ListPendingFn = func(ctx context.Context) ([]*domain.VacationRequest, error) {
		return []*domain.VacationRequest{}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/vacation/pending", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Total)
	assert.Len(t, resp.Requests, 0)
}

func TestAdminReview_ApproveInsufficientBalance(t *testing.T) {
	deps := setupAdminTest(t)

	vacation := sampleVacation("vac-1", "user-10", domain.StatusPending, 10)
	// User only has 5 days left, but request is for 10
	user := sampleUser("user-10", "emp@test.com", "Employee", domain.RoleEmployee, 5)

	deps.vacRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-1" {
			return vacation, nil
		}
		return nil, nil
	}

	deps.userRepo.GetByIDFn = func(ctx context.Context, id string) (*domain.User, error) {
		if id == "user-10" {
			return user, nil
		}
		return nil, nil
	}

	body := `{"status":"approved"}`
	req := httptest.NewRequest(http.MethodPut, "/api/admin/vacation/vac-1/review", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInsufficientBalance, resp.Code)
}

func TestAdminResetBalances_SettingsRepoError(t *testing.T) {
	deps := setupAdminTest(t)

	deps.settingsRepo.GetFn = func(ctx context.Context) (*domain.Settings, error) {
		return nil, fmt.Errorf("settings database error")
	}

	req := httptest.NewRequest(http.MethodPost, "/api/admin/users/reset-balances", nil)
	w := httptest.NewRecorder()
	deps.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInternal, resp.Code)
}
