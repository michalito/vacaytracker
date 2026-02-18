package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vacaytracker-api/internal/domain"
	"vacaytracker-api/internal/dto"
	"vacaytracker-api/internal/service"
	"vacaytracker-api/internal/testutil"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func intPtr(v int) *int { return &v }

func newUserService(repo *testutil.MockUserRepository) *service.UserService {
	authSvc := service.NewAuthService(&testutil.MockUserRepository{}, "test-secret-key-for-jwt-signing")
	return service.NewUserService(repo, authSvc)
}

func existingUser() *domain.User {
	startDate := "2024-01-15"
	return &domain.User{
		ID:              "user-1",
		Email:           "alice@example.com",
		PasswordHash:    "$2a$10$fakehash",
		Name:            "Alice",
		Role:            domain.RoleEmployee,
		VacationBalance: 25,
		StartDate:       &startDate,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}
}

func existingAdmin() *domain.User {
	return &domain.User{
		ID:              "admin-1",
		Email:           "admin@example.com",
		PasswordHash:    "$2a$10$fakehash",
		Name:            "Admin",
		Role:            domain.RoleAdmin,
		VacationBalance: 25,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}
}

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

func TestCreate_Success_DefaultBalance(t *testing.T) {
	var createdUser *domain.User
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, email string) (bool, error) {
			assert.Equal(t, "new@example.com", email)
			return false, nil
		},
		CreateFn: func(_ context.Context, user *domain.User) error {
			createdUser = user
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:    "new@example.com",
		Password: "securepassword",
		Name:     "New User",
		Role:     "employee",
	})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "New User", user.Name)
	assert.Equal(t, domain.RoleEmployee, user.Role)
	assert.Equal(t, 25, user.VacationBalance) // default
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.PasswordHash)
	assert.Nil(t, user.StartDate)
	// Ensure the same object was passed to repo.Create
	assert.Equal(t, createdUser, user)
}

func TestCreate_Success_CustomBalance(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, _ *domain.User) error {
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:           "bob@example.com",
		Password:        "securepassword",
		Name:            "Bob",
		Role:            "admin",
		VacationBalance: intPtr(30),
		StartDate:       "2024-06-01",
	})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 30, user.VacationBalance)
	assert.Equal(t, domain.RoleAdmin, user.Role)
	require.NotNil(t, user.StartDate)
	assert.Equal(t, "2024-06-01", *user.StartDate)
}

func TestCreate_Success_ZeroBalance(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, _ *domain.User) error {
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:           "zero@example.com",
		Password:        "securepassword",
		Name:            "Zero Balance",
		Role:            "employee",
		VacationBalance: intPtr(0),
	})

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 0, user.VacationBalance)
}

func TestCreate_DuplicateEmail(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:    "existing@example.com",
		Password: "securepassword",
		Name:     "Dup User",
		Role:     "employee",
	})

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrAlreadyExists, appErr.Code)
	assert.Equal(t, 409, appErr.HTTPStatus)
}

func TestCreate_EmailExistsCheckError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:    "test@example.com",
		Password: "securepassword",
		Name:     "Test",
		Role:     "employee",
	})

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestCreate_PasswordTooShort(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:    "short@example.com",
		Password: "abc",
		Name:     "Short Pass",
		Role:     "employee",
	})

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code) // hash fails -> internal error
}

func TestCreate_RepoCreateError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		EmailExistsFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
		},
		CreateFn: func(_ context.Context, _ *domain.User) error {
			return errors.New("db insert failed")
		},
	}

	svc := newUserService(repo)
	user, err := svc.Create(context.Background(), dto.CreateUserRequest{
		Email:    "fail@example.com",
		Password: "securepassword",
		Name:     "Fail User",
		Role:     "employee",
	})

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_Success_ChangeName(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			assert.Equal(t, "user-1", id)
			// Return a copy so the test can compare before/after
			u := *original
			return &u, nil
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			assert.Equal(t, "Updated Name", user.Name)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Name: "Updated Name",
	}, "other-admin-id")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Updated Name", user.Name)
	assert.Equal(t, original.Email, user.Email) // unchanged
}

func TestUpdate_Success_ChangeEmail_Unique(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		EmailExistsExcludingFn: func(_ context.Context, email, excludeID string) (bool, error) {
			assert.Equal(t, "newemail@example.com", email)
			assert.Equal(t, "user-1", excludeID)
			return false, nil
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			assert.Equal(t, "newemail@example.com", user.Email)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Email: "newemail@example.com",
	}, "admin-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "newemail@example.com", user.Email)
}

func TestUpdate_EmailConflict(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		EmailExistsExcludingFn: func(_ context.Context, _ string, _ string) (bool, error) {
			return true, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Email: "taken@example.com",
	}, "admin-1")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrAlreadyExists, appErr.Code)
	assert.Equal(t, 409, appErr.HTTPStatus)
}

func TestUpdate_CannotModifyOwnRole(t *testing.T) {
	admin := existingAdmin()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *admin
			return &u, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "admin-1", dto.UpdateUserRequest{
		Role: "employee",
	}, "admin-1") // same user

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrForbidden, appErr.Code)
	assert.Equal(t, 403, appErr.HTTPStatus)
	assert.Contains(t, appErr.Message, "cannot modify your own role")
}

func TestUpdate_CannotDemoteLastAdmin(t *testing.T) {
	admin := existingAdmin()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *admin
			return &u, nil
		},
		CountByRoleFn: func(_ context.Context, role domain.Role) (int, error) {
			assert.Equal(t, domain.RoleAdmin, role)
			return 1, nil // only one admin
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "admin-1", dto.UpdateUserRequest{
		Role: "employee",
	}, "other-admin-id") // different user doing the update

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrForbidden, appErr.Code)
	assert.Contains(t, appErr.Message, "cannot demote the last admin")
}

func TestUpdate_AllowDemoteWhenMultipleAdmins(t *testing.T) {
	admin := existingAdmin()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *admin
			return &u, nil
		},
		CountByRoleFn: func(_ context.Context, _ domain.Role) (int, error) {
			return 3, nil // multiple admins
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			assert.Equal(t, domain.RoleEmployee, user.Role)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "admin-1", dto.UpdateUserRequest{
		Role: "employee",
	}, "other-admin-id")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, domain.RoleEmployee, user.Role)
}

func TestUpdate_UserNotFound(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, nil // not found
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "nonexistent", dto.UpdateUserRequest{
		Name: "X",
	}, "admin-1")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrNotFound, appErr.Code)
	assert.Equal(t, 404, appErr.HTTPStatus)
}

func TestUpdate_GetByIDError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, errors.New("db read error")
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Name: "X",
	}, "admin-1")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestUpdate_Success_ChangeBalance(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			assert.Equal(t, 42, user.VacationBalance)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		VacationBalance: intPtr(42),
	}, "admin-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 42, user.VacationBalance)
}

func TestUpdate_Success_ChangeStartDate(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			require.NotNil(t, user.StartDate)
			assert.Equal(t, "2025-03-01", *user.StartDate)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		StartDate: "2025-03-01",
	}, "admin-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.StartDate)
	assert.Equal(t, "2025-03-01", *user.StartDate)
}

func TestUpdate_SameEmailNoConflictCheck(t *testing.T) {
	// When the email is the same as existing, EmailExistsExcluding should NOT be called.
	original := existingUser()
	emailCheckCalled := false
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		EmailExistsExcludingFn: func(_ context.Context, _ string, _ string) (bool, error) {
			emailCheckCalled = true
			return false, nil
		},
		UpdateFn: func(_ context.Context, _ *domain.User) error {
			return nil
		},
	}

	svc := newUserService(repo)
	_, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Email: "alice@example.com", // same as existing
		Name:  "Alice Updated",
	}, "admin-1")

	require.NoError(t, err)
	assert.False(t, emailCheckCalled, "should not check email uniqueness when email is unchanged")
}

func TestUpdate_PromoteEmployeeToAdmin(t *testing.T) {
	emp := existingUser() // employee
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *emp
			return &u, nil
		},
		UpdateFn: func(_ context.Context, user *domain.User) error {
			assert.Equal(t, domain.RoleAdmin, user.Role)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Role: "admin",
	}, "admin-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, domain.RoleAdmin, user.Role)
}

func TestUpdate_RepoUpdateError(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		UpdateFn: func(_ context.Context, _ *domain.User) error {
			return errors.New("db update failed")
		},
	}

	svc := newUserService(repo)
	user, err := svc.Update(context.Background(), "user-1", dto.UpdateUserRequest{
		Name: "Fail",
	}, "admin-1")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_Success(t *testing.T) {
	emp := existingUser()
	deleteCalled := false
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *emp
			return &u, nil
		},
		DeleteFn: func(_ context.Context, id string) error {
			assert.Equal(t, "user-1", id)
			deleteCalled = true
			return nil
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "user-1", "admin-1")

	require.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestDelete_CannotDeleteSelf(t *testing.T) {
	repo := &testutil.MockUserRepository{}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "admin-1", "admin-1") // same ID

	require.Error(t, err)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrForbidden, appErr.Code)
	assert.Contains(t, appErr.Message, "cannot delete your own account")
}

func TestDelete_CannotDeleteLastAdmin(t *testing.T) {
	admin := existingAdmin()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *admin
			return &u, nil
		},
		CountByRoleFn: func(_ context.Context, role domain.Role) (int, error) {
			assert.Equal(t, domain.RoleAdmin, role)
			return 1, nil
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "admin-1", "other-admin-id")

	require.Error(t, err)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrForbidden, appErr.Code)
	assert.Contains(t, appErr.Message, "cannot delete the last admin")
}

func TestDelete_AdminWithMultipleAdmins(t *testing.T) {
	admin := existingAdmin()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *admin
			return &u, nil
		},
		CountByRoleFn: func(_ context.Context, _ domain.Role) (int, error) {
			return 3, nil
		},
		DeleteFn: func(_ context.Context, _ string) error {
			return nil
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "admin-1", "other-admin-id")

	require.NoError(t, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, nil
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "nonexistent", "admin-1")

	require.Error(t, err)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrNotFound, appErr.Code)
	assert.Equal(t, 404, appErr.HTTPStatus)
}

func TestDelete_GetByIDError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "user-1", "admin-1")

	require.Error(t, err)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestDelete_RepoDeleteError(t *testing.T) {
	emp := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *emp
			return &u, nil
		},
		DeleteFn: func(_ context.Context, _ string) error {
			return errors.New("db delete failed")
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "user-1", "admin-1")

	require.Error(t, err)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestDelete_EmployeeSkipsAdminCount(t *testing.T) {
	// Deleting an employee should NOT call CountByRole
	emp := existingUser()
	countCalled := false
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *emp
			return &u, nil
		},
		CountByRoleFn: func(_ context.Context, _ domain.Role) (int, error) {
			countCalled = true
			return 0, nil
		},
		DeleteFn: func(_ context.Context, _ string) error {
			return nil
		},
	}

	svc := newUserService(repo)
	err := svc.Delete(context.Background(), "user-1", "admin-1")

	require.NoError(t, err)
	assert.False(t, countCalled, "CountByRole should not be called when deleting an employee")
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestGetByID_Success(t *testing.T) {
	expected := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			assert.Equal(t, "user-1", id)
			u := *expected
			return &u, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.GetByID(context.Background(), "user-1")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "user-1", user.ID)
	assert.Equal(t, "alice@example.com", user.Email)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.GetByID(context.Background(), "nonexistent")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrNotFound, appErr.Code)
	assert.Equal(t, 404, appErr.HTTPStatus)
}

func TestGetByID_RepoError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	user, err := svc.GetByID(context.Background(), "user-1")

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestList_Success_Defaults(t *testing.T) {
	users := []*domain.User{existingUser(), existingAdmin()}
	repo := &testutil.MockUserRepository{
		GetAllFn: func(_ context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
			assert.Nil(t, role)
			assert.Empty(t, search)
			assert.Equal(t, 20, limit)
			assert.Equal(t, 0, offset) // page 1 -> offset 0
			return users, 2, nil
		},
	}

	svc := newUserService(repo)
	result, total, err := svc.List(context.Background(), nil, "", 1, 20)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 2, total)
}

func TestList_PageNormalization(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "page less than 1 normalizes to 1",
			page:           0,
			limit:          20,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "negative page normalizes to 1",
			page:           -5,
			limit:          10,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name:           "limit over 100 normalizes to 20",
			page:           1,
			limit:          200,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "limit of 0 normalizes to 20",
			page:           1,
			limit:          0,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "negative limit normalizes to 20",
			page:           1,
			limit:          -1,
			expectedLimit:  20,
			expectedOffset: 0,
		},
		{
			name:           "valid page 2 with limit 10",
			page:           2,
			limit:          10,
			expectedLimit:  10,
			expectedOffset: 10,
		},
		{
			name:           "valid page 3 with limit 25",
			page:           3,
			limit:          25,
			expectedLimit:  25,
			expectedOffset: 50,
		},
		{
			name:           "limit exactly 100 is accepted",
			page:           1,
			limit:          100,
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name:           "limit exactly 101 normalizes to 20",
			page:           1,
			limit:          101,
			expectedLimit:  20,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testutil.MockUserRepository{
				GetAllFn: func(_ context.Context, _ *domain.Role, _ string, limit, offset int) ([]*domain.User, int, error) {
					assert.Equal(t, tt.expectedLimit, limit, "limit mismatch")
					assert.Equal(t, tt.expectedOffset, offset, "offset mismatch")
					return nil, 0, nil
				},
			}

			svc := newUserService(repo)
			_, _, err := svc.List(context.Background(), nil, "", tt.page, tt.limit)
			require.NoError(t, err)
		})
	}
}

func TestList_WithRoleFilter(t *testing.T) {
	adminRole := domain.RoleAdmin
	repo := &testutil.MockUserRepository{
		GetAllFn: func(_ context.Context, role *domain.Role, _ string, _ int, _ int) ([]*domain.User, int, error) {
			require.NotNil(t, role)
			assert.Equal(t, domain.RoleAdmin, *role)
			return []*domain.User{existingAdmin()}, 1, nil
		},
	}

	svc := newUserService(repo)
	result, total, err := svc.List(context.Background(), &adminRole, "", 1, 20)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestList_WithSearch(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetAllFn: func(_ context.Context, _ *domain.Role, search string, _ int, _ int) ([]*domain.User, int, error) {
			assert.Equal(t, "alice", search)
			return []*domain.User{existingUser()}, 1, nil
		},
	}

	svc := newUserService(repo)
	result, total, err := svc.List(context.Background(), nil, "alice", 1, 20)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 1, total)
}

func TestList_RepoError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetAllFn: func(_ context.Context, _ *domain.Role, _ string, _ int, _ int) ([]*domain.User, int, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	users, total, err := svc.List(context.Background(), nil, "", 1, 20)

	require.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, 0, total)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestList_EmptyResult(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetAllFn: func(_ context.Context, _ *domain.Role, _ string, _ int, _ int) ([]*domain.User, int, error) {
			return []*domain.User{}, 0, nil
		},
	}

	svc := newUserService(repo)
	result, total, err := svc.List(context.Background(), nil, "", 1, 20)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, 0, total)
}

// ---------------------------------------------------------------------------
// UpdateBalance
// ---------------------------------------------------------------------------

func TestUpdateBalance_Success(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, id string) (*domain.User, error) {
			assert.Equal(t, "user-1", id)
			u := *original
			return &u, nil
		},
		UpdateVacationBalanceFn: func(_ context.Context, id string, balance int) error {
			assert.Equal(t, "user-1", id)
			assert.Equal(t, 30, balance)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "user-1", 30)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 30, user.VacationBalance)
	assert.Equal(t, "user-1", user.ID)
}

func TestUpdateBalance_Success_SetToZero(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		UpdateVacationBalanceFn: func(_ context.Context, _ string, balance int) error {
			assert.Equal(t, 0, balance)
			return nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "user-1", 0)

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, 0, user.VacationBalance)
}

func TestUpdateBalance_UserNotFound(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "nonexistent", 10)

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrNotFound, appErr.Code)
	assert.Equal(t, 404, appErr.HTTPStatus)
}

func TestUpdateBalance_NegativeBalance(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "user-1", -5)

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrValidation, appErr.Code)
	assert.Equal(t, 400, appErr.HTTPStatus)
	assert.Contains(t, appErr.Message, "negative")
}

func TestUpdateBalance_GetByIDError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "user-1", 10)

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestUpdateBalance_RepoUpdateError(t *testing.T) {
	original := existingUser()
	repo := &testutil.MockUserRepository{
		GetByIDFn: func(_ context.Context, _ string) (*domain.User, error) {
			u := *original
			return &u, nil
		},
		UpdateVacationBalanceFn: func(_ context.Context, _ string, _ int) error {
			return errors.New("db update failed")
		},
	}

	svc := newUserService(repo)
	user, err := svc.UpdateBalance(context.Background(), "user-1", 10)

	require.Error(t, err)
	assert.Nil(t, user)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

// ---------------------------------------------------------------------------
// ResetAllBalances
// ---------------------------------------------------------------------------

func TestResetAllBalances_Success(t *testing.T) {
	repo := &testutil.MockUserRepository{
		UpdateAllBalancesFn: func(_ context.Context, balance int) (int64, error) {
			assert.Equal(t, 25, balance)
			return 10, nil
		},
	}

	svc := newUserService(repo)
	count, err := svc.ResetAllBalances(context.Background(), 25)

	require.NoError(t, err)
	assert.Equal(t, 10, count)
}

func TestResetAllBalances_Success_ZeroDays(t *testing.T) {
	repo := &testutil.MockUserRepository{
		UpdateAllBalancesFn: func(_ context.Context, balance int) (int64, error) {
			assert.Equal(t, 0, balance)
			return 5, nil
		},
	}

	svc := newUserService(repo)
	count, err := svc.ResetAllBalances(context.Background(), 0)

	require.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestResetAllBalances_NegativeDays(t *testing.T) {
	repo := &testutil.MockUserRepository{}

	svc := newUserService(repo)
	count, err := svc.ResetAllBalances(context.Background(), -1)

	require.Error(t, err)
	assert.Equal(t, 0, count)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrValidation, appErr.Code)
	assert.Equal(t, 400, appErr.HTTPStatus)
	assert.Contains(t, appErr.Message, "negative")
}

func TestResetAllBalances_RepoError(t *testing.T) {
	repo := &testutil.MockUserRepository{
		UpdateAllBalancesFn: func(_ context.Context, _ int) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	svc := newUserService(repo)
	count, err := svc.ResetAllBalances(context.Background(), 25)

	require.Error(t, err)
	assert.Equal(t, 0, count)
	var appErr *dto.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, dto.ErrInternal, appErr.Code)
}

func TestResetAllBalances_NoUsersAffected(t *testing.T) {
	repo := &testutil.MockUserRepository{
		UpdateAllBalancesFn: func(_ context.Context, _ int) (int64, error) {
			return 0, nil
		},
	}

	svc := newUserService(repo)
	count, err := svc.ResetAllBalances(context.Background(), 25)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
