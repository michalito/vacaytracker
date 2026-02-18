package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

const testJWTSecret = "test-secret-key-that-is-long-enough"

// newTestAuthService creates an AuthService with a mock repo and the default test secret.
func newTestAuthService(repo *testutil.MockUserRepository) *service.AuthService {
	return service.NewAuthService(repo, testJWTSecret)
}

// testUser returns a sample domain.User for testing.
func testUser() *domain.User {
	return &domain.User{
		ID:              "usr_test001",
		Email:           "employee@example.com",
		PasswordHash:    "", // populated as needed per test
		Name:            "Test Employee",
		Role:            domain.RoleEmployee,
		VacationBalance: 25,
		EmailPreferences: domain.EmailPreferences{
			VacationUpdates:   true,
			WeeklyDigest:      false,
			TeamNotifications: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// boolPtr is a helper to get a pointer to a bool value.
func boolPtr(b bool) *bool {
	return &b
}

// assertAppError checks that err is an *dto.AppError with the expected code.
func assertAppError(t *testing.T, err error, expectedCode string) {
	t.Helper()
	require.Error(t, err)
	var appErr *dto.AppError
	require.True(t, errors.As(err, &appErr), "expected *dto.AppError, got %T: %v", err, err)
	assert.Equal(t, expectedCode, appErr.Code)
}

// --------------------------------------------------------------------------
// HashPassword
// --------------------------------------------------------------------------

func TestHashPassword(t *testing.T) {
	svc := newTestAuthService(&testutil.MockUserRepository{})

	t.Run("success", func(t *testing.T) {
		hash, err := svc.HashPassword("validPassword123")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		// bcrypt hashes always start with "$2"
		assert.True(t, strings.HasPrefix(hash, "$2"), "expected bcrypt hash prefix")
	})

	t.Run("minimum valid length (6 chars)", func(t *testing.T) {
		hash, err := svc.HashPassword("abcdef")
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("maximum valid length (72 chars)", func(t *testing.T) {
		password := strings.Repeat("a", 72)
		hash, err := svc.HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("too short", func(t *testing.T) {
		hash, err := svc.HashPassword("abc")
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "at least 6 characters")
	})

	t.Run("empty password", func(t *testing.T) {
		hash, err := svc.HashPassword("")
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "at least 6 characters")
	})

	t.Run("too long (>72 chars)", func(t *testing.T) {
		password := strings.Repeat("x", 73)
		hash, err := svc.HashPassword(password)
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "cannot exceed 72 characters")
	})

	t.Run("different hashes for same password", func(t *testing.T) {
		hash1, err1 := svc.HashPassword("samePassword")
		require.NoError(t, err1)
		hash2, err2 := svc.HashPassword("samePassword")
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "bcrypt should produce different hashes due to random salt")
	})
}

// --------------------------------------------------------------------------
// VerifyPassword
// --------------------------------------------------------------------------

func TestVerifyPassword(t *testing.T) {
	svc := newTestAuthService(&testutil.MockUserRepository{})

	hash, err := svc.HashPassword("correctPassword")
	require.NoError(t, err)

	t.Run("correct password", func(t *testing.T) {
		assert.True(t, svc.VerifyPassword("correctPassword", hash))
	})

	t.Run("wrong password", func(t *testing.T) {
		assert.False(t, svc.VerifyPassword("wrongPassword", hash))
	})

	t.Run("empty password against hash", func(t *testing.T) {
		assert.False(t, svc.VerifyPassword("", hash))
	})

	t.Run("password against invalid hash", func(t *testing.T) {
		assert.False(t, svc.VerifyPassword("correctPassword", "not-a-bcrypt-hash"))
	})
}

// --------------------------------------------------------------------------
// GenerateToken
// --------------------------------------------------------------------------

func TestGenerateToken(t *testing.T) {
	svc := newTestAuthService(&testutil.MockUserRepository{})

	t.Run("generates valid JWT string", func(t *testing.T) {
		user := testUser()
		tokenStr, err := svc.GenerateToken(user)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenStr)

		// JWT has 3 dot-separated parts
		parts := strings.Split(tokenStr, ".")
		assert.Len(t, parts, 3, "JWT should have 3 segments")
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		user := testUser()
		tokenStr, err := svc.GenerateToken(user)
		require.NoError(t, err)

		// Validate the token we just created
		claims, err := svc.ValidateToken(tokenStr)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Name, claims.Name)
		assert.Equal(t, user.Role, claims.Role)
		assert.Equal(t, "vacaytracker", claims.Issuer)
		// Note: RegisteredClaims.Subject is shadowed by the custom UserID field
		// (both have json tag "sub"), so Subject may be empty. The UserID field
		// is the authoritative source for the subject claim.
	})

	t.Run("token for admin user", func(t *testing.T) {
		user := &domain.User{
			ID:    "usr_admin001",
			Email: "admin@example.com",
			Name:  "Captain Admin",
			Role:  domain.RoleAdmin,
		}
		tokenStr, err := svc.GenerateToken(user)
		require.NoError(t, err)

		claims, err := svc.ValidateToken(tokenStr)
		require.NoError(t, err)
		assert.Equal(t, domain.RoleAdmin, claims.Role)
	})
}

// --------------------------------------------------------------------------
// ValidateToken
// --------------------------------------------------------------------------

func TestValidateToken(t *testing.T) {
	svc := newTestAuthService(&testutil.MockUserRepository{})

	t.Run("valid token returns correct claims", func(t *testing.T) {
		user := testUser()
		tokenStr, err := svc.GenerateToken(user)
		require.NoError(t, err)

		claims, err := svc.ValidateToken(tokenStr)
		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Name, claims.Name)
		assert.Equal(t, user.Role, claims.Role)
	})

	t.Run("expired token returns token expired error", func(t *testing.T) {
		user := testUser()

		// Manually craft an expired token using the same secret
		now := time.Now()
		claims := service.JWTClaims{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Role:   user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // expired 1 hour ago
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				Issuer:    "vacaytracker",
				Subject:   user.ID,
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		expiredToken, err := token.SignedString([]byte(testJWTSecret))
		require.NoError(t, err)

		result, err := svc.ValidateToken(expiredToken)
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenExpired)
	})

	t.Run("malformed token returns token invalid error", func(t *testing.T) {
		result, err := svc.ValidateToken("not.a.valid.jwt")
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenInvalid)
	})

	t.Run("empty token returns token invalid error", func(t *testing.T) {
		result, err := svc.ValidateToken("")
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenInvalid)
	})

	t.Run("random garbage returns token invalid error", func(t *testing.T) {
		result, err := svc.ValidateToken("aaaaaaa.bbbbbbb.ccccccc")
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenInvalid)
	})

	t.Run("wrong signing key returns token invalid error", func(t *testing.T) {
		// Generate a token with a different secret
		otherSvc := service.NewAuthService(&testutil.MockUserRepository{}, "completely-different-secret-key!!")
		user := testUser()
		tokenStr, err := otherSvc.GenerateToken(user)
		require.NoError(t, err)

		// Try to validate with the original service (different secret)
		result, err := svc.ValidateToken(tokenStr)
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenInvalid)
	})

	t.Run("token signed with different algorithm returns invalid", func(t *testing.T) {
		// Create a token with "none" algorithm (unsigned) -- the parser should reject it
		claims := service.JWTClaims{
			UserID: "usr_test001",
			Email:  "test@example.com",
			Name:   "Test",
			Role:   domain.RoleEmployee,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "vacaytracker",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
		unsignedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)

		result, err := svc.ValidateToken(unsignedToken)
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrAuthTokenInvalid)
	})
}

// --------------------------------------------------------------------------
// Login
// --------------------------------------------------------------------------

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		password := "securePassword123"
		hash, err := svc.HashPassword(password)
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
				if email == user.Email {
					return user, nil
				}
				return nil, nil
			},
		}
		svc = newTestAuthService(repo)

		token, returnedUser, err := svc.Login(ctx, user.Email, password)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		require.NotNil(t, returnedUser)
		assert.Equal(t, user.ID, returnedUser.ID)
		assert.Equal(t, user.Email, returnedUser.Email)

		// Verify the returned token is valid
		claims, err := svc.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
	})

	t.Run("wrong email - user not found", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
				return nil, nil // user not found
			},
		}
		svc := newTestAuthService(repo)

		token, user, err := svc.Login(ctx, "nonexistent@example.com", "anypassword")
		assert.Empty(t, token)
		assert.Nil(t, user)
		assertAppError(t, err, dto.ErrInvalidCredentials)
	})

	t.Run("wrong password", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		hash, err := svc.HashPassword("correctPassword")
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
				return user, nil
			},
		}
		svc = newTestAuthService(repo)

		token, returnedUser, err := svc.Login(ctx, user.Email, "wrongPassword")
		assert.Empty(t, token)
		assert.Nil(t, returnedUser)
		assertAppError(t, err, dto.ErrInvalidCredentials)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
				return nil, errors.New("database connection lost")
			},
		}
		svc := newTestAuthService(repo)

		token, user, err := svc.Login(ctx, "test@example.com", "password")
		assert.Empty(t, token)
		assert.Nil(t, user)
		assertAppError(t, err, dto.ErrInvalidCredentials)
	})
}

// --------------------------------------------------------------------------
// GetUserByID
// --------------------------------------------------------------------------

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		user := testUser()
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				if id == user.ID {
					return user, nil
				}
				return nil, nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, user.ID, result.ID)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Name, result.Name)
	})

	t.Run("user not found - returns nil from repo", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.GetUserByID(ctx, "usr_nonexistent")
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrUserNotFound)
	})

	t.Run("user not found - repo returns error", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, errors.New("database error")
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.GetUserByID(ctx, "usr_test001")
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrUserNotFound)
	})
}

// --------------------------------------------------------------------------
// ChangePassword
// --------------------------------------------------------------------------

func TestChangePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		currentPassword := "oldPassword123"
		hash, err := svc.HashPassword(currentPassword)
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		var updatedHash string
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				if id == user.ID {
					return user, nil
				}
				return nil, nil
			},
			UpdatePasswordFn: func(_ context.Context, id, passwordHash string) error {
				updatedHash = passwordHash
				return nil
			},
		}
		svc = newTestAuthService(repo)

		err = svc.ChangePassword(ctx, user.ID, currentPassword, "newPassword456")
		require.NoError(t, err)
		assert.NotEmpty(t, updatedHash)
		assert.NotEqual(t, hash, updatedHash, "new hash should differ from old hash")

		// The new hash should verify against the new password
		assert.True(t, svc.VerifyPassword("newPassword456", updatedHash))
	})

	t.Run("wrong current password", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		hash, err := svc.HashPassword("correctOldPassword")
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return user, nil
			},
		}
		svc = newTestAuthService(repo)

		err = svc.ChangePassword(ctx, user.ID, "wrongOldPassword", "newPassword456")
		assertAppError(t, err, dto.ErrInvalidCredentials)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, nil
			},
		}
		svc := newTestAuthService(repo)

		err := svc.ChangePassword(ctx, "usr_nonexistent", "anyPassword", "newPassword456")
		assertAppError(t, err, dto.ErrUserNotFound)
	})

	t.Run("new password too short", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		currentPassword := "oldPassword123"
		hash, err := svc.HashPassword(currentPassword)
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return user, nil
			},
		}
		svc = newTestAuthService(repo)

		err = svc.ChangePassword(ctx, user.ID, currentPassword, "short")
		// HashPassword will fail, which triggers ErrInternalError
		assertAppError(t, err, dto.ErrInternal)
	})

	t.Run("new password too long", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		currentPassword := "oldPassword123"
		hash, err := svc.HashPassword(currentPassword)
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return user, nil
			},
		}
		svc = newTestAuthService(repo)

		err = svc.ChangePassword(ctx, user.ID, currentPassword, strings.Repeat("a", 73))
		assertAppError(t, err, dto.ErrInternal)
	})

	t.Run("repo UpdatePassword error", func(t *testing.T) {
		svc := newTestAuthService(&testutil.MockUserRepository{})
		currentPassword := "oldPassword123"
		hash, err := svc.HashPassword(currentPassword)
		require.NoError(t, err)

		user := testUser()
		user.PasswordHash = hash

		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return user, nil
			},
			UpdatePasswordFn: func(_ context.Context, id, passwordHash string) error {
				return errors.New("write failed")
			},
		}
		svc = newTestAuthService(repo)

		err = svc.ChangePassword(ctx, user.ID, currentPassword, "newValidPassword")
		assertAppError(t, err, dto.ErrInternal)
	})
}

// --------------------------------------------------------------------------
// UpdateEmailPreferences
// --------------------------------------------------------------------------

func TestUpdateEmailPreferences(t *testing.T) {
	ctx := context.Background()

	t.Run("update single field - VacationUpdates", func(t *testing.T) {
		user := testUser()
		user.EmailPreferences = domain.EmailPreferences{
			VacationUpdates:   true,
			WeeklyDigest:      false,
			TeamNotifications: true,
		}

		var savedPrefs domain.EmailPreferences
		updatedUser := *user // copy for the second GetByID call
		updatedUser.EmailPreferences.VacationUpdates = false

		callCount := 0
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				callCount++
				if callCount == 1 {
					return user, nil
				}
				return &updatedUser, nil
			},
			UpdateEmailPreferencesFn: func(_ context.Context, id string, prefs domain.EmailPreferences) error {
				savedPrefs = prefs
				return nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, user.ID, &dto.UpdateEmailPreferencesRequest{
			VacationUpdates: boolPtr(false),
		})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Check that only VacationUpdates was changed in saved prefs
		assert.False(t, savedPrefs.VacationUpdates)
		assert.False(t, savedPrefs.WeeklyDigest)    // unchanged
		assert.True(t, savedPrefs.TeamNotifications) // unchanged
	})

	t.Run("update multiple fields", func(t *testing.T) {
		user := testUser()
		user.EmailPreferences = domain.EmailPreferences{
			VacationUpdates:   true,
			WeeklyDigest:      false,
			TeamNotifications: true,
		}

		var savedPrefs domain.EmailPreferences
		updatedUser := *user
		updatedUser.EmailPreferences = domain.EmailPreferences{
			VacationUpdates:   false,
			WeeklyDigest:      true,
			TeamNotifications: false,
		}

		callCount := 0
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				callCount++
				if callCount == 1 {
					return user, nil
				}
				return &updatedUser, nil
			},
			UpdateEmailPreferencesFn: func(_ context.Context, id string, prefs domain.EmailPreferences) error {
				savedPrefs = prefs
				return nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, user.ID, &dto.UpdateEmailPreferencesRequest{
			VacationUpdates:   boolPtr(false),
			WeeklyDigest:      boolPtr(true),
			TeamNotifications: boolPtr(false),
		})
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.False(t, savedPrefs.VacationUpdates)
		assert.True(t, savedPrefs.WeeklyDigest)
		assert.False(t, savedPrefs.TeamNotifications)
	})

	t.Run("update no fields - keeps existing values", func(t *testing.T) {
		user := testUser()
		user.EmailPreferences = domain.EmailPreferences{
			VacationUpdates:   true,
			WeeklyDigest:      false,
			TeamNotifications: true,
		}

		var savedPrefs domain.EmailPreferences
		callCount := 0
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				callCount++
				return user, nil
			},
			UpdateEmailPreferencesFn: func(_ context.Context, id string, prefs domain.EmailPreferences) error {
				savedPrefs = prefs
				return nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, user.ID, &dto.UpdateEmailPreferencesRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)

		// All values should be preserved
		assert.True(t, savedPrefs.VacationUpdates)
		assert.False(t, savedPrefs.WeeklyDigest)
		assert.True(t, savedPrefs.TeamNotifications)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, nil
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, "usr_nonexistent", &dto.UpdateEmailPreferencesRequest{
			VacationUpdates: boolPtr(true),
		})
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrUserNotFound)
	})

	t.Run("repo GetByID error", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				return nil, errors.New("database error")
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, "usr_test001", &dto.UpdateEmailPreferencesRequest{
			VacationUpdates: boolPtr(true),
		})
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrUserNotFound)
	})

	t.Run("repo UpdateEmailPreferences error", func(t *testing.T) {
		user := testUser()
		callCount := 0
		repo := &testutil.MockUserRepository{
			GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
				callCount++
				return user, nil
			},
			UpdateEmailPreferencesFn: func(_ context.Context, id string, prefs domain.EmailPreferences) error {
				return errors.New("write failed")
			},
		}
		svc := newTestAuthService(repo)

		result, err := svc.UpdateEmailPreferences(ctx, user.ID, &dto.UpdateEmailPreferencesRequest{
			VacationUpdates: boolPtr(false),
		})
		assert.Nil(t, result)
		assertAppError(t, err, dto.ErrInternal)
	})
}

// --------------------------------------------------------------------------
// CreateInitialAdmin
// --------------------------------------------------------------------------

func TestCreateInitialAdmin(t *testing.T) {
	ctx := context.Background()

	t.Run("creates admin when none exists", func(t *testing.T) {
		var createdUser *domain.User
		repo := &testutil.MockUserRepository{
			EmailExistsFn: func(_ context.Context, email string) (bool, error) {
				return false, nil
			},
			CreateFn: func(_ context.Context, user *domain.User) error {
				createdUser = user
				return nil
			},
		}
		svc := newTestAuthService(repo)

		err := svc.CreateInitialAdmin(ctx, "admin@example.com", "adminPassword123", "Captain Admin", 30)
		require.NoError(t, err)

		require.NotNil(t, createdUser)
		assert.Equal(t, "usr_admin001", createdUser.ID)
		assert.Equal(t, "admin@example.com", createdUser.Email)
		assert.Equal(t, "Captain Admin", createdUser.Name)
		assert.Equal(t, domain.RoleAdmin, createdUser.Role)
		assert.Equal(t, 30, createdUser.VacationBalance)
		assert.NotEmpty(t, createdUser.PasswordHash)

		// Verify the stored hash matches the password
		assert.True(t, svc.VerifyPassword("adminPassword123", createdUser.PasswordHash))

		// Verify default email preferences
		defaults := domain.DefaultEmailPreferences()
		assert.Equal(t, defaults, createdUser.EmailPreferences)
	})

	t.Run("does nothing when admin already exists", func(t *testing.T) {
		createCalled := false
		repo := &testutil.MockUserRepository{
			EmailExistsFn: func(_ context.Context, email string) (bool, error) {
				return true, nil // admin already exists
			},
			CreateFn: func(_ context.Context, user *domain.User) error {
				createCalled = true
				return nil
			},
		}
		svc := newTestAuthService(repo)

		err := svc.CreateInitialAdmin(ctx, "admin@example.com", "adminPassword123", "Captain Admin", 30)
		require.NoError(t, err)
		assert.False(t, createCalled, "Create should not be called when admin already exists")
	})

	t.Run("handles EmailExists repo error", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			EmailExistsFn: func(_ context.Context, email string) (bool, error) {
				return false, errors.New("database error")
			},
		}
		svc := newTestAuthService(repo)

		err := svc.CreateInitialAdmin(ctx, "admin@example.com", "adminPassword123", "Captain Admin", 30)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check admin existence")
	})

	t.Run("handles Create repo error", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			EmailExistsFn: func(_ context.Context, email string) (bool, error) {
				return false, nil
			},
			CreateFn: func(_ context.Context, user *domain.User) error {
				return errors.New("unique constraint violated")
			},
		}
		svc := newTestAuthService(repo)

		err := svc.CreateInitialAdmin(ctx, "admin@example.com", "adminPassword123", "Captain Admin", 30)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create admin user")
	})

	t.Run("handles invalid password for admin", func(t *testing.T) {
		repo := &testutil.MockUserRepository{
			EmailExistsFn: func(_ context.Context, email string) (bool, error) {
				return false, nil
			},
		}
		svc := newTestAuthService(repo)

		err := svc.CreateInitialAdmin(ctx, "admin@example.com", "short", "Captain Admin", 30)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash admin password")
	})
}

// --------------------------------------------------------------------------
// Integration-style: Login then ValidateToken round-trip
// --------------------------------------------------------------------------

func TestLoginAndValidateTokenRoundTrip(t *testing.T) {
	ctx := context.Background()

	svc := newTestAuthService(&testutil.MockUserRepository{})
	password := "securePassword123"
	hash, err := svc.HashPassword(password)
	require.NoError(t, err)

	user := testUser()
	user.PasswordHash = hash

	repo := &testutil.MockUserRepository{
		GetByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			if email == user.Email {
				return user, nil
			}
			return nil, nil
		},
	}
	svc = newTestAuthService(repo)

	// Login to get a token
	token, returnedUser, err := svc.Login(ctx, user.Email, password)
	require.NoError(t, err)
	require.NotNil(t, returnedUser)

	// Validate the token
	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Claims match the user
	assert.Equal(t, returnedUser.ID, claims.UserID)
	assert.Equal(t, returnedUser.Email, claims.Email)
	assert.Equal(t, returnedUser.Name, claims.Name)
	assert.Equal(t, returnedUser.Role, claims.Role)

	// Token expiry is in the future
	assert.True(t, claims.ExpiresAt.After(time.Now()), "token should expire in the future")
}
