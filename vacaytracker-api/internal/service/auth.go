package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/repository"
)

// JWTClaims represents the claims stored in JWT tokens
type JWTClaims struct {
	UserID string      `json:"sub"`
	Email  string      `json:"email"`
	Name   string      `json:"name"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: 24 * time.Hour, // 24 hour token expiry
	}
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	// Validate password length (bcrypt silently truncates at 72 bytes)
	if len(password) < 6 {
		return "", fmt.Errorf("password must be at least 6 characters")
	}
	if len(password) > 72 {
		return "", fmt.Errorf("password cannot exceed 72 characters")
	}

	// Cost of 10 is a good balance between security and performance
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword compares a password with its hash
func (s *AuthService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates a JWT token for a user
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "vacaytracker",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, dto.ErrTokenExpiredError()
		}
		return nil, dto.ErrTokenInvalidError()
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, dto.ErrTokenInvalidError()
	}

	return claims, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, dto.ErrInvalidCredentialsError()
	}

	// Verify password
	if !s.VerifyPassword(password, user.PasswordHash) {
		return "", nil, dto.ErrInvalidCredentialsError()
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, dto.ErrInternalError()
	}

	return token, user, nil
}

// GetUserByID retrieves a user by their ID
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		return nil, dto.ErrUserNotFoundError()
	}
	return user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return dto.ErrUserNotFoundError()
	}

	// Verify current password
	if !s.VerifyPassword(currentPassword, user.PasswordHash) {
		return dto.ErrInvalidCredentialsError()
	}

	// Hash new password
	newHash, err := s.HashPassword(newPassword)
	if err != nil {
		return dto.ErrInternalError()
	}

	// Update password
	if err := s.userRepo.UpdatePassword(ctx, userID, newHash); err != nil {
		return dto.ErrInternalError()
	}

	return nil
}

// UpdateEmailPreferences updates a user's email notification preferences
func (s *AuthService) UpdateEmailPreferences(ctx context.Context, userID string, updates *dto.UpdateEmailPreferencesRequest) (*domain.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, dto.ErrUserNotFoundError()
	}

	// Apply updates (only update fields that are provided)
	if updates.VacationUpdates != nil {
		user.EmailPreferences.VacationUpdates = *updates.VacationUpdates
	}
	if updates.WeeklyDigest != nil {
		user.EmailPreferences.WeeklyDigest = *updates.WeeklyDigest
	}
	if updates.TeamNotifications != nil {
		user.EmailPreferences.TeamNotifications = *updates.TeamNotifications
	}

	// Save preferences
	if err := s.userRepo.UpdateEmailPreferences(ctx, userID, user.EmailPreferences); err != nil {
		return nil, dto.ErrInternalError()
	}

	// Get updated user
	return s.userRepo.GetByID(ctx, userID)
}

// CreateInitialAdmin creates the initial admin user if it doesn't exist
func (s *AuthService) CreateInitialAdmin(ctx context.Context, email, password, name string, defaultBalance int) error {
	// Check if admin already exists
	exists, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check admin existence: %w", err)
	}

	if exists {
		return nil // Admin already exists, nothing to do
	}

	// Hash password
	hash, err := s.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	// Create admin user
	admin := &domain.User{
		ID:               "usr_admin001",
		Email:            email,
		PasswordHash:     hash,
		Name:             name,
		Role:             domain.RoleAdmin,
		VacationBalance:  defaultBalance,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	if err := s.userRepo.Create(ctx, admin); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	return nil
}
