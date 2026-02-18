package handler_test

import (
	"context"
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
	"golang.org/x/crypto/bcrypt"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/handler"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

// jwtSecret must be at least 32 chars for the AuthService constructor.
const testJWTSecret = "test-secret-key-that-is-at-least-32-chars-long"

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newTestUser builds a domain.User with a bcrypt password hash.
// Uses bcrypt.MinCost so tests run fast.
func newTestUser(id, email, name string, role domain.Role, balance int, password string) *domain.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	now := time.Now()
	return &domain.User{
		ID:               id,
		Email:            email,
		PasswordHash:     string(hash),
		Name:             name,
		Role:             role,
		VacationBalance:  balance,
		EmailPreferences: domain.DefaultEmailPreferences(),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// authContextMiddleware returns a gin.HandlerFunc that sets auth context
// values so the handler under test can read them via middleware.GetUserID etc.
func authContextMiddleware(userID, email, name string, role domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		testutil.SetAuthContext(c, userID, email, name, role)
		c.Next()
	}
}

// ===================================================================
// Login tests
// ===================================================================

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "password123")

	mockRepo := &testutil.MockUserRepository{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			if email == "test@example.com" {
				return user, nil
			}
			return nil, nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.LoginResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Token)
	require.NotNil(t, resp.User)
	assert.Equal(t, "user-1", resp.User.ID)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "Test User", resp.User.Name)
	assert.Equal(t, "employee", resp.User.Role)
	assert.Equal(t, 25, resp.User.VacationBalance)
}

func TestLogin_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

func TestLogin_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	// Email is required by binding:"required,email"
	body := `{"password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Email")
}

func TestLogin_MissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	// Password is required by binding:"required,min=6,max=72"
	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
	assert.Contains(t, resp.Message, "Password")
}

func TestLogin_WrongEmail_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil // user not found
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	body := `{"email":"nobody@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInvalidCredentials, resp.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "correctPassword")

	mockRepo := &testutil.MockUserRepository{
		GetByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return user, nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.POST("/api/auth/login", h.Login)

	body := `{"email":"test@example.com","password":"wrongPassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInvalidCredentials, resp.Code)
}

// ===================================================================
// Me tests
// ===================================================================

func TestMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "password123")

	mockRepo := &testutil.MockUserRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-1" {
				return user, nil
			}
			return nil, nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.GET("/api/auth/me",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.Me,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user-1", resp.ID)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "employee", resp.Role)
	assert.Equal(t, 25, resp.VacationBalance)
}

func TestMe_NoAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	// No auth context middleware — userID will be empty.
	router.GET("/api/auth/me", h.Me)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

func TestMe_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, nil // user not found
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.GET("/api/auth/me",
		authContextMiddleware("deleted-user", "deleted@example.com", "Deleted", domain.RoleEmployee),
		h.Me,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrUserNotFound, resp.Code)
}

// ===================================================================
// ChangePassword tests
// ===================================================================

func TestChangePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "oldPassword1")

	mockRepo := &testutil.MockUserRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-1" {
				return user, nil
			}
			return nil, nil
		},
		UpdatePasswordFn: func(ctx context.Context, id, passwordHash string) error {
			return nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/password",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.ChangePassword,
	)

	body := `{"currentPassword":"oldPassword1","newPassword":"newPassword1"}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.MessageResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Password changed successfully", resp.Message)
}

func TestChangePassword_NoAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/password", h.ChangePassword)

	body := `{"currentPassword":"old","newPassword":"newPassword1"}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

func TestChangePassword_InvalidBody_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/password",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.ChangePassword,
	)

	// Missing both currentPassword and newPassword
	body := `{}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "realPassword")

	mockRepo := &testutil.MockUserRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			if id == "user-1" {
				return user, nil
			}
			return nil, nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/password",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.ChangePassword,
	)

	body := `{"currentPassword":"wrongPassword","newPassword":"newPassword1"}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrInvalidCredentials, resp.Code)
}

// ===================================================================
// UpdateEmailPreferences tests
// ===================================================================

func TestUpdateEmailPreferences_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	user := newTestUser("user-1", "test@example.com", "Test User", domain.RoleEmployee, 25, "password123")

	// After update, GetByID returns the user with modified preferences.
	updatedUser := *user // shallow copy
	updatedUser.EmailPreferences = domain.EmailPreferences{
		VacationUpdates:   false,
		WeeklyDigest:      true,
		TeamNotifications: true,
	}

	callCount := 0
	mockRepo := &testutil.MockUserRepository{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			callCount++
			if id == "user-1" {
				// First call returns original; second call (after update) returns updated.
				if callCount <= 1 {
					return user, nil
				}
				return &updatedUser, nil
			}
			return nil, nil
		},
		UpdateEmailPreferencesFn: func(ctx context.Context, id string, prefs domain.EmailPreferences) error {
			return nil
		},
	}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/email-preferences",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.UpdateEmailPreferences,
	)

	body := `{"vacationUpdates":false,"weeklyDigest":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/email-preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// The response shape is {"emailPreferences": {...}}
	var resp struct {
		EmailPreferences domain.EmailPreferences `json:"emailPreferences"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.EmailPreferences.VacationUpdates)
	assert.True(t, resp.EmailPreferences.WeeklyDigest)
	assert.True(t, resp.EmailPreferences.TeamNotifications)
}

func TestUpdateEmailPreferences_NoAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/email-preferences", h.UpdateEmailPreferences)

	body := `{"vacationUpdates":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/auth/email-preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrAuthTokenMissing, resp.Code)
}

func TestUpdateEmailPreferences_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &testutil.MockUserRepository{}
	authService := service.NewAuthService(mockRepo, testJWTSecret)
	h := handler.NewAuthHandler(authService)

	router := gin.New()
	router.PUT("/api/auth/email-preferences",
		authContextMiddleware("user-1", "test@example.com", "Test User", domain.RoleEmployee),
		h.UpdateEmailPreferences,
	)

	// Send malformed JSON
	req := httptest.NewRequest(http.MethodPut, "/api/auth/email-preferences", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp dto.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, dto.ErrValidation, resp.Code)
}

// ===================================================================
// GetPublic settings tests
// ===================================================================

func TestGetPublicSettings_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSettings := &testutil.MockSettingsRepository{
		GetFn: func(ctx context.Context) (*domain.Settings, error) {
			s := domain.DefaultSettings()
			s.DefaultVacationDays = 30
			s.VacationResetMonth = 4
			return &s, nil
		},
	}
	h := handler.NewSettingsHandler(mockSettings)

	router := gin.New()
	router.GET("/api/settings/public", h.GetPublic)

	req := httptest.NewRequest(http.MethodGet, "/api/settings/public", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		DefaultVacationDays int `json:"defaultVacationDays"`
		VacationResetMonth  int `json:"vacationResetMonth"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 30, resp.DefaultVacationDays)
	assert.Equal(t, 4, resp.VacationResetMonth)
}

func TestGetPublicSettings_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSettings := &testutil.MockSettingsRepository{
		GetFn: func(ctx context.Context) (*domain.Settings, error) {
			return nil, fmt.Errorf("database connection lost")
		},
	}
	h := handler.NewSettingsHandler(mockSettings)

	router := gin.New()
	router.GET("/api/settings/public", h.GetPublic)

	req := httptest.NewRequest(http.MethodGet, "/api/settings/public", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "Failed to get settings")
}
