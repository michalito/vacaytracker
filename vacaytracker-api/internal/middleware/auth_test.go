package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

const testJWTSecret = "test-secret-key-for-middleware-tests-32chars!"

// newTestAuthService creates an AuthService backed by a mock user repo.
func newTestAuthService() *service.AuthService {
	mockRepo := &testutil.MockUserRepository{}
	return service.NewAuthService(mockRepo, testJWTSecret)
}

// generateValidToken creates a valid JWT for the given user via the real AuthService.
func generateValidToken(t *testing.T, user *domain.User) string {
	t.Helper()
	authService := newTestAuthService()
	token, err := authService.GenerateToken(user)
	require.NoError(t, err)
	return token
}

// generateExpiredToken creates a JWT whose expiry is in the past.
func generateExpiredToken(t *testing.T) string {
	t.Helper()
	claims := service.JWTClaims{
		UserID: "usr_expired",
		Email:  "expired@example.com",
		Name:   "Expired User",
		Role:   domain.RoleEmployee,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "vacaytracker",
			Subject:   "usr_expired",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testJWTSecret))
	require.NoError(t, err)
	return signed
}

// ─── AuthMiddleware Tests ───

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_MISSING", body["code"])
}

func TestAuthMiddleware_InvalidFormat_NoBearerPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token some-token-value")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", body["code"])
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer totally-invalid-jwt-string")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_INVALID", body["code"])
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	testUser := &domain.User{
		ID:    "usr_test123",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  domain.RoleEmployee,
	}
	token := generateValidToken(t, testUser)

	var capturedUserID, capturedEmail, capturedName string
	var capturedRole domain.Role

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		capturedUserID = GetUserID(c)
		capturedEmail = GetUserEmail(c)
		nameVal, _ := c.Get(ContextKeyName)
		capturedName, _ = nameVal.(string)
		capturedRole = GetUserRole(c)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "usr_test123", capturedUserID)
	assert.Equal(t, "test@example.com", capturedEmail)
	assert.Equal(t, "Test User", capturedName)
	assert.Equal(t, domain.RoleEmployee, capturedRole)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	expiredToken := generateExpiredToken(t)

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_EXPIRED", body["code"])
}

func TestAuthMiddleware_BearerCaseInsensitive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authService := newTestAuthService()

	testUser := &domain.User{
		ID:    "usr_case",
		Email: "case@example.com",
		Name:  "Case User",
		Role:  domain.RoleAdmin,
	}
	token := generateValidToken(t, testUser)

	router := gin.New()
	router.Use(AuthMiddleware(authService))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Use lowercase "bearer"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// ─── AdminMiddleware Tests ───

func TestAdminMiddleware_NoRoleInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(AdminMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_MISSING", body["code"])
}

func TestAdminMiddleware_EmployeeRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// Simulate AuthMiddleware having set the role
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, domain.RoleEmployee)
		c.Next()
	})
	router.Use(AdminMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ADMIN_REQUIRED", body["code"])
}

func TestAdminMiddleware_AdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, domain.RoleAdmin)
		c.Next()
	})
	router.Use(AdminMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// ─── EmployeeMiddleware Tests ───

func TestEmployeeMiddleware_NoRoleInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(EmployeeMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "AUTH_TOKEN_MISSING", body["code"])
}

func TestEmployeeMiddleware_AdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, domain.RoleAdmin)
		c.Next()
	})
	router.Use(EmployeeMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", body["code"])
}

func TestEmployeeMiddleware_EmployeeRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, domain.RoleEmployee)
		c.Next()
	})
	router.Use(EmployeeMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// ─── Helper Function Tests ───

func TestGetUserID_Present(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyUserID, "usr_abc123")

	assert.Equal(t, "usr_abc123", GetUserID(c))
}

func TestGetUserID_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.Equal(t, "", GetUserID(c))
}

func TestGetUserRole_Present(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRole, domain.RoleAdmin)

	assert.Equal(t, domain.RoleAdmin, GetUserRole(c))
}

func TestGetUserRole_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.Equal(t, domain.Role(""), GetUserRole(c))
}

func TestIsAdmin_True(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRole, domain.RoleAdmin)

	assert.True(t, IsAdmin(c))
}

func TestIsAdmin_False(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRole, domain.RoleEmployee)

	assert.False(t, IsAdmin(c))
}

func TestIsEmployee_True(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRole, domain.RoleEmployee)

	assert.True(t, IsEmployee(c))
}

func TestIsEmployee_False(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyRole, domain.RoleAdmin)

	assert.False(t, IsEmployee(c))
}

func TestGetUserEmail_Present(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(ContextKeyEmail, "user@example.com")

	assert.Equal(t, "user@example.com", GetUserEmail(c))
}

func TestGetUserEmail_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.Equal(t, "", GetUserEmail(c))
}

func TestGetClaims_Present(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	expected := &service.JWTClaims{
		UserID: "usr_claims",
		Email:  "claims@example.com",
		Name:   "Claims User",
		Role:   domain.RoleEmployee,
	}
	c.Set(ContextKeyClaims, expected)

	result := GetClaims(c)
	require.NotNil(t, result)
	assert.Equal(t, "usr_claims", result.UserID)
	assert.Equal(t, "claims@example.com", result.Email)
}

func TestGetClaims_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	assert.Nil(t, GetClaims(c))
}

func TestAdminMiddleware_WrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// Set role as a plain string instead of domain.Role — type assertion should fail
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, "admin")
		c.Next()
	})
	router.Use(AdminMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var body map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "ADMIN_REQUIRED", body["code"])
}
