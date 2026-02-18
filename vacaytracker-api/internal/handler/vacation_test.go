package handler_test

import (
	"context"
	"encoding/json"
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

// setupVacationRouter creates a Gin router with the vacation handler routes and
// a middleware that injects the given user identity into the Gin context.
// The role is stored as domain.Role so that middleware.GetUserRole can retrieve it.
func setupVacationRouter(h *handler.VacationHandler, userID, email, name string, role domain.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	authMiddleware := func(c *gin.Context) {
		// Store userID, email, name as strings (matching middleware.GetUserID, etc.)
		c.Set("userID", userID)
		c.Set("email", email)
		c.Set("name", name)
		// Store role as domain.Role so middleware.GetUserRole type assertion succeeds.
		c.Set("role", role)
		c.Next()
	}

	r.POST("/api/vacation/request", authMiddleware, h.Create)
	r.GET("/api/vacation/requests", authMiddleware, h.List)
	r.GET("/api/vacation/requests/:id", authMiddleware, h.Get)
	r.DELETE("/api/vacation/requests/:id", authMiddleware, h.Cancel)
	r.GET("/api/vacation/team", authMiddleware, h.Team)

	return r
}

// setupVacationRouterNoAuth creates a Gin router with no auth middleware (context values not set).
func setupVacationRouterNoAuth(h *handler.VacationHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/api/vacation/request", h.Create)
	r.GET("/api/vacation/requests", h.List)
	r.GET("/api/vacation/requests/:id", h.Get)
	r.DELETE("/api/vacation/requests/:id", h.Cancel)
	r.GET("/api/vacation/team", h.Team)

	return r
}

// newTestEmailService creates an EmailService with no Resend API key so emails
// are silently skipped (no external calls).
func newTestEmailService() *service.EmailService {
	cfg := &config.Config{AppURL: "http://localhost:3000"}
	return service.NewEmailService(cfg)
}

// futureMonday returns the next Monday that is at least daysAhead days from now.
func futureMonday(daysAhead int) time.Time {
	d := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, daysAhead)
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, 1)
	}
	return d
}

// ============================================
// Create Tests
// ============================================

func TestCreate_Success_Employee(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	// Pick two future weekdays (Mon-Tue) to ensure non-zero business days
	monday := futureMonday(30)
	tuesday := monday.AddDate(0, 0, 1)
	startDateStr := monday.Format("02/01/2006")
	endDateStr := tuesday.Format("02/01/2006")
	startDateISO := monday.Format("2006-01-02")
	endDateISO := tuesday.Format("2006-01-02")

	// Service needs: userRepo.GetByID, settingsRepo.Get (default), vacationRepo.HasOverlap,
	// vacationRepo.Create, vacationRepo.GetByID
	userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		if id == "user-1" {
			return &domain.User{
				ID:              "user-1",
				Email:           "employee@test.com",
				Name:            "Test Employee",
				Role:            domain.RoleEmployee,
				VacationBalance: 20,
			}, nil
		}
		return nil, nil
	}

	vacationRepo.HasOverlapFn = func(_ context.Context, userID, start, end string) (bool, error) {
		return false, nil
	}

	var createdVacation *domain.VacationRequest
	vacationRepo.CreateFn = func(_ context.Context, req *domain.VacationRequest) error {
		createdVacation = req
		return nil
	}

	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if createdVacation != nil && createdVacation.ID == id {
			return &domain.VacationRequest{
				ID:        createdVacation.ID,
				UserID:    "user-1",
				StartDate: startDateISO,
				EndDate:   endDateISO,
				TotalDays: 2,
				Status:    domain.StatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, nil
	}

	// Also mock GetByRole for the email notification goroutine (handler reads admins).
	userRepo.GetByRoleFn = func(_ context.Context, role domain.Role) ([]*domain.User, error) {
		return nil, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	body := `{"startDate":"` + startDateStr + `","endDate":"` + endDateStr + `","reason":"Family trip"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/vacation/request", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.VacationRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, "user-1", resp.UserID)
	assert.Equal(t, startDateISO, resp.StartDate)
	assert.Equal(t, endDateISO, resp.EndDate)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, 2, resp.TotalDays)

	// Allow goroutine to finish before test cleanup
	time.Sleep(50 * time.Millisecond)
}

func TestCreate_InvalidJSON(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodPost, "/api/vacation/request", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

func TestCreate_NoAuthContext(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouterNoAuth(h)

	body := `{"startDate":"15/06/2027","endDate":"20/06/2027"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/vacation/request", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

func TestCreate_InsufficientBalance(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	// User with only 1 day balance requesting 2+ business days
	monday := futureMonday(30)
	tuesday := monday.AddDate(0, 0, 1)
	startDateStr := monday.Format("02/01/2006")
	endDateStr := tuesday.Format("02/01/2006")

	userRepo.GetByIDFn = func(_ context.Context, id string) (*domain.User, error) {
		return &domain.User{
			ID:              "user-1",
			Email:           "employee@test.com",
			Name:            "Test Employee",
			Role:            domain.RoleEmployee,
			VacationBalance: 1, // only 1 day, requesting 2
		}, nil
	}

	vacationRepo.HasOverlapFn = func(_ context.Context, userID, start, end string) (bool, error) {
		return false, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	body := `{"startDate":"` + startDateStr + `","endDate":"` + endDateStr + `"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/vacation/request", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrInsufficientBalance, resp.Code)
}

func TestCreate_InvalidDates(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	// Test with badly formatted date (not DD/MM/YYYY)
	body := `{"startDate":"2027-06-15","endDate":"2027-06-20"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/vacation/request", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "invalid start date format")
}

// ============================================
// List Tests
// ============================================

func TestList_Success_NoFilters(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.ListByUserFn = func(_ context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, "user-1", userID)
		assert.Nil(t, status)
		assert.Nil(t, year)
		return []*domain.VacationRequest{
			{
				ID:        "vac-1",
				UserID:    "user-1",
				StartDate: "2027-06-15",
				EndDate:   "2027-06-20",
				TotalDays: 5,
				Status:    domain.StatusPending,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "vac-2",
				UserID:    "user-1",
				StartDate: "2027-07-01",
				EndDate:   "2027-07-05",
				TotalDays: 4,
				Status:    domain.StatusApproved,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Requests, 2)
	assert.Equal(t, "vac-1", resp.Requests[0].ID)
	assert.Equal(t, "vac-2", resp.Requests[1].ID)
}

func TestList_WithStatusFilter(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.ListByUserFn = func(_ context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
		assert.Equal(t, "user-1", userID)
		require.NotNil(t, status)
		assert.Equal(t, domain.StatusApproved, *status)
		return []*domain.VacationRequest{
			{
				ID:        "vac-2",
				UserID:    "user-1",
				StartDate: "2027-07-01",
				EndDate:   "2027-07-05",
				TotalDays: 4,
				Status:    domain.StatusApproved,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests?status=approved", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationListResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Total)
	assert.Equal(t, "approved", resp.Requests[0].Status)
}

func TestList_InvalidStatus(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests?status=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Invalid status")
}

func TestList_NoAuthContext(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouterNoAuth(h)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

// ============================================
// Get Tests
// ============================================

func TestGet_Success_OwnRequest(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-1" {
			return &domain.VacationRequest{
				ID:        "vac-1",
				UserID:    "user-1",
				UserName:  "Test Employee",
				StartDate: "2027-06-15",
				EndDate:   "2027-06-20",
				TotalDays: 5,
				Status:    domain.StatusPending,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}
		return nil, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests/vac-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "vac-1", resp.ID)
	assert.Equal(t, "user-1", resp.UserID)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, 5, resp.TotalDays)
}

func TestGet_NotFound(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	// GetByID returns nil (not found) -- the service wraps this in an AppError
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		return nil, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

func TestGet_OtherUserRequest_NotAdmin(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		return &domain.VacationRequest{
			ID:        "vac-1",
			UserID:    "other-user", // belongs to a different user
			StartDate: "2027-06-15",
			EndDate:   "2027-06-20",
			TotalDays: 5,
			Status:    domain.StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	// Logged in as user-1 (employee), trying to view other-user's request
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests/vac-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrForbidden, resp.Code)
}

func TestGet_AdminCanViewAnyRequest(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		return &domain.VacationRequest{
			ID:        "vac-1",
			UserID:    "other-user", // belongs to a different user
			UserName:  "Other Employee",
			StartDate: "2027-06-15",
			EndDate:   "2027-06-20",
			TotalDays: 5,
			Status:    domain.StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	// Logged in as admin
	router := setupVacationRouter(h, "admin-1", "admin@test.com", "Admin User", domain.RoleAdmin)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/requests/vac-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.VacationRequestResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "vac-1", resp.ID)
	assert.Equal(t, "other-user", resp.UserID)
}

// ============================================
// Cancel Tests
// ============================================

func TestCancel_Success(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		if id == "vac-1" {
			return &domain.VacationRequest{
				ID:        "vac-1",
				UserID:    "user-1",
				StartDate: "2027-06-15",
				EndDate:   "2027-06-20",
				TotalDays: 5,
				Status:    domain.StatusPending, // must be pending to cancel
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}
		return nil, nil
	}

	deleteCalled := false
	vacationRepo.DeleteFn = func(_ context.Context, id string) error {
		assert.Equal(t, "vac-1", id)
		deleteCalled = true
		return nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodDelete, "/api/vacation/requests/vac-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, deleteCalled)

	var resp dto.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp.Message, "cancelled successfully")
}

func TestCancel_NotFound(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	// GetByID returns nil -- service returns NotFound AppError
	vacationRepo.GetByIDFn = func(_ context.Context, id string) (*domain.VacationRequest, error) {
		return nil, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodDelete, "/api/vacation/requests/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrNotFound, resp.Code)
}

func TestCancel_NoAuthContext(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouterNoAuth(h)

	req, _ := http.NewRequest(http.MethodDelete, "/api/vacation/requests/vac-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

// ============================================
// Team Tests
// ============================================

func TestTeam_Success_DefaultMonthYear(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	now := time.Now()
	expectedMonth := int(now.Month())
	expectedYear := now.Year()

	vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, expectedMonth, month)
		assert.Equal(t, expectedYear, year)
		return []*domain.TeamVacation{
			{
				ID:        "vac-1",
				UserID:    "user-2",
				UserName:  "Team Member",
				StartDate: "2027-06-15",
				EndDate:   "2027-06-20",
				TotalDays: 5,
			},
		}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/team", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TeamVacationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, expectedMonth, resp.Month)
	assert.Equal(t, expectedYear, resp.Year)
	assert.Len(t, resp.Vacations, 1)
	assert.Equal(t, "vac-1", resp.Vacations[0].ID)
	assert.Equal(t, "Team Member", resp.Vacations[0].UserName)
}

func TestTeam_Success_ExplicitMonthYear(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationRepo.ListTeamFn = func(_ context.Context, month, year int) ([]*domain.TeamVacation, error) {
		assert.Equal(t, 8, month)
		assert.Equal(t, 2027, year)
		return []*domain.TeamVacation{}, nil
	}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/team?month=8&year=2027", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.TeamVacationResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 8, resp.Month)
	assert.Equal(t, 2027, resp.Year)
	assert.Empty(t, resp.Vacations)
}

func TestTeam_InvalidMonth(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/team?month=13", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Invalid month")
}

func TestTeam_InvalidYear(t *testing.T) {
	vacationRepo := &testutil.MockVacationRepository{}
	userRepo := &testutil.MockUserRepository{}
	settingsRepo := &testutil.MockSettingsRepository{}
	transactor := &testutil.MockTransactor{}

	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo, transactor)
	emailService := newTestEmailService()

	h := handler.NewVacationHandler(vacationService, vacationRepo, userRepo, emailService)
	router := setupVacationRouter(h, "user-1", "employee@test.com", "Test Employee", domain.RoleEmployee)

	req, _ := http.NewRequest(http.MethodGet, "/api/vacation/team?year=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Invalid year")
}
