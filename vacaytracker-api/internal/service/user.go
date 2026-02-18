package service

import (
	"context"

	"github.com/google/uuid"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/repository"
)

// UserService handles user management business logic
type UserService struct {
	userRepo    repository.UserRepository
	authService *AuthService
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to check email")
	}
	if exists {
		return nil, dto.ErrConflictError("email already exists")
	}

	// Hash password
	hash, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to hash password")
	}

	// Set defaults
	balance := 25
	if req.VacationBalance != nil {
		balance = *req.VacationBalance
	}

	var startDate *string
	if req.StartDate != "" {
		startDate = &req.StartDate
	}

	user := &domain.User{
		ID:               uuid.New().String(),
		Email:            req.Email,
		PasswordHash:     hash,
		Name:             req.Name,
		Role:             domain.Role(req.Role),
		VacationBalance:  balance,
		StartDate:        startDate,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to create user")
	}

	return user, nil
}

// Update updates a user's information
func (s *UserService) Update(ctx context.Context, id string, req dto.UpdateUserRequest, currentUserID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return nil, dto.ErrNotFoundError("user")
	}

	// Check email uniqueness if changing
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.EmailExistsExcluding(ctx, req.Email, id)
		if err != nil {
			return nil, dto.ErrInternalErrorWithMessage("failed to check email")
		}
		if exists {
			return nil, dto.ErrConflictError("email already exists")
		}
		user.Email = req.Email
	}

	// Check role change restrictions
	if req.Role != "" && domain.Role(req.Role) != user.Role {
		// Cannot modify own role
		if id == currentUserID {
			return nil, dto.ErrForbiddenError("cannot modify your own role")
		}

		// Cannot demote if last admin
		if user.Role == domain.RoleAdmin && req.Role == string(domain.RoleEmployee) {
			count, err := s.userRepo.CountByRole(ctx, domain.RoleAdmin)
			if err != nil {
				return nil, dto.ErrInternalErrorWithMessage("failed to count admins")
			}
			if count <= 1 {
				return nil, dto.ErrForbiddenError("cannot demote the last admin")
			}
		}

		user.Role = domain.Role(req.Role)
	}

	// Update other fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.VacationBalance != nil {
		user.VacationBalance = *req.VacationBalance
	}
	if req.StartDate != "" {
		user.StartDate = &req.StartDate
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to update user")
	}

	return user, nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id, currentUserID string) error {
	// Cannot delete self
	if id == currentUserID {
		return dto.ErrForbiddenError("cannot delete your own account")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return dto.ErrNotFoundError("user")
	}

	// Cannot delete last admin
	if user.Role == domain.RoleAdmin {
		count, err := s.userRepo.CountByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return dto.ErrInternalErrorWithMessage("failed to count admins")
		}
		if count <= 1 {
			return dto.ErrForbiddenError("cannot delete the last admin")
		}
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return dto.ErrInternalErrorWithMessage("failed to delete user")
	}

	return nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return nil, dto.ErrNotFoundError("user")
	}
	return user, nil
}

// List lists all users with optional filtering and pagination
func (s *UserService) List(ctx context.Context, role *domain.Role, search string, page, limit int) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, total, err := s.userRepo.GetAll(ctx, role, search, limit, offset)
	if err != nil {
		return nil, 0, dto.ErrInternalErrorWithMessage("failed to list users")
	}

	return users, total, nil
}

// UpdateBalance updates a user's vacation balance
func (s *UserService) UpdateBalance(ctx context.Context, id string, balance int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to get user")
	}
	if user == nil {
		return nil, dto.ErrNotFoundError("user")
	}

	if balance < 0 {
		return nil, dto.ErrValidationError("vacation balance cannot be negative")
	}

	if err := s.userRepo.UpdateVacationBalance(ctx, id, balance); err != nil {
		return nil, dto.ErrInternalErrorWithMessage("failed to update vacation balance")
	}

	user.VacationBalance = balance
	return user, nil
}

// ResetAllBalances resets all employee vacation balances to the specified default value
func (s *UserService) ResetAllBalances(ctx context.Context, defaultDays int) (int, error) {
	if defaultDays < 0 {
		return 0, dto.ErrValidationError("default vacation days cannot be negative")
	}

	count, err := s.userRepo.UpdateAllBalances(ctx, defaultDays)
	if err != nil {
		return 0, dto.ErrInternalErrorWithMessage("failed to reset vacation balances")
	}

	return int(count), nil
}
