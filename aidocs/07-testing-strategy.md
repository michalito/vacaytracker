# 07 - Testing Strategy

> Comprehensive testing plan for backend and frontend components

## Table of Contents

1. [Testing Overview](#testing-overview)
2. [Backend Testing](#backend-testing)
3. [Frontend Testing](#frontend-testing)
4. [Integration Testing](#integration-testing)
5. [E2E Testing](#e2e-testing)
6. [Test Data & Fixtures](#test-data--fixtures)
7. [Coverage Requirements](#coverage-requirements)

---

## Testing Overview

### Testing Pyramid

```
                    ┌─────────┐
                    │   E2E   │  ~10 tests
                    ├─────────┤
                 ┌──┴─────────┴──┐
                 │  Integration  │  ~30 tests
                 ├───────────────┤
          ┌──────┴───────────────┴──────┐
          │          Unit Tests         │  ~100 tests
          └─────────────────────────────┘
```

### Test Categories

| Category | Backend | Frontend | Purpose |
|----------|---------|----------|---------|
| Unit | Go test | Vitest | Individual functions/components |
| Integration | HTTP tests | Component tests | Module interactions |
| E2E | - | Playwright | Full user flows |

### Testing Stack

**Backend:**
- `testing` package (standard library)
- `testify/assert` for assertions
- `httptest` for handler tests
- In-memory SQLite for database tests

**Frontend:**
- Vitest for unit tests
- Testing Library for component tests
- Playwright for E2E tests

---

## Backend Testing

### Unit Tests

#### Service Tests

- [ ] **Auth Service Tests** `internal/service/auth_test.go`

```go
package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	s := &AuthService{}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"short password", "12345", false},
		{"empty password", "", false},
		{"long password", string(make([]byte, 73)), true}, // bcrypt max is 72
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := s.HashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	s := &AuthService{}
	password := "testPassword123"
	hash, _ := s.HashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"correct password", password, hash, true},
		{"wrong password", "wrongPassword", hash, false},
		{"empty password", "", hash, false},
		{"invalid hash", password, "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.VerifyPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key-minimum-32-chars!"}
	s := NewAuthService(nil, cfg)

	user := &domain.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  domain.RoleEmployee,
	}

	token, err := s.GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := s.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)
}

func TestValidateToken_Invalid(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-key-minimum-32-chars!"}
	s := NewAuthService(nil, cfg)

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"invalid format", "not-a-jwt"},
		{"wrong signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.ValidateToken(tt.token)
			assert.Error(t, err)
		})
	}
}
```

- [ ] **Vacation Service Tests** `internal/service/vacation_test.go`

```go
package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

func TestParseDDMMYYYY(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDate string
		wantErr  bool
	}{
		{"valid date", "15/01/2024", "2024-01-15", false},
		{"end of month", "31/12/2024", "2024-12-31", false},
		{"first day", "01/01/2024", "2024-01-01", false},
		{"invalid format", "2024-01-15", "", true},
		{"invalid day", "32/01/2024", "", true},
		{"invalid month", "15/13/2024", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDDMMYYYY(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantDate, result.Format("2006-01-02"))
		})
	}
}

func TestCalculateBusinessDays(t *testing.T) {
	tests := []struct {
		name            string
		start           string
		end             string
		excludeWeekends bool
		excludedDays    []int
		expected        int
	}{
		{
			name:            "weekdays only, Mon-Fri",
			start:           "2024-01-15", // Monday
			end:             "2024-01-19", // Friday
			excludeWeekends: true,
			excludedDays:    []int{0, 6},
			expected:        5,
		},
		{
			name:            "span weekend, exclude",
			start:           "2024-01-15", // Monday
			end:             "2024-01-22", // Monday
			excludeWeekends: true,
			excludedDays:    []int{0, 6},
			expected:        6, // Mon-Fri (5) + Mon (1)
		},
		{
			name:            "include weekends",
			start:           "2024-01-15", // Monday
			end:             "2024-01-21", // Sunday
			excludeWeekends: false,
			excludedDays:    []int{},
			expected:        7,
		},
		{
			name:            "single day",
			start:           "2024-01-15",
			end:             "2024-01-15",
			excludeWeekends: true,
			excludedDays:    []int{0, 6},
			expected:        1,
		},
		{
			name:            "only Saturday excluded",
			start:           "2024-01-15", // Monday
			end:             "2024-01-21", // Sunday
			excludeWeekends: true,
			excludedDays:    []int{6},
			expected:        6, // Excludes only Saturday
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := time.Parse("2006-01-02", tt.start)
			end, _ := time.Parse("2006-01-02", tt.end)
			policy := domain.WeekendPolicy{
				ExcludeWeekends: tt.excludeWeekends,
				ExcludedDays:    tt.excludedDays,
			}

			result := calculateBusinessDays(start, end, policy)
			assert.Equal(t, tt.expected, result)
		})
	}
}
```

- [ ] **User Service Tests** `internal/service/user_test.go`

```go
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string, excludeID string) (bool, error) {
	args := m.Called(ctx, email, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CountByRole(ctx context.Context, role domain.Role) (int, error) {
	args := m.Called(ctx, role)
	return args.Int(0), args.Error(1)
}

// ... implement other methods

func TestUserService_Delete_CannotDeleteSelf(t *testing.T) {
	mockRepo := new(MockUserRepository)
	s := &UserService{userRepo: mockRepo}

	err := s.Delete(context.Background(), "user-123", "user-123")
	assert.Equal(t, ErrCannotDeleteSelf, err)
}

func TestUserService_Delete_CannotDeleteLastAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetByID", mock.Anything, "admin-123").Return(&domain.User{
		ID:   "admin-123",
		Role: domain.RoleAdmin,
	}, nil)
	mockRepo.On("CountByRole", mock.Anything, domain.RoleAdmin).Return(1, nil)

	s := &UserService{userRepo: mockRepo}

	err := s.Delete(context.Background(), "admin-123", "other-user")
	assert.Equal(t, ErrCannotDeleteLastAdmin, err)
}
```

#### Repository Tests

- [ ] **User Repository Tests** `internal/repository/sqlite/user_test.go`

```go
package sqlite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

func setupTestDB(t *testing.T) *DB {
	db, err := New(":memory:")
	require.NoError(t, err)

	err = db.RunMigrations("../../../migrations")
	require.NoError(t, err)

	return db
}

func TestUserRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:               "test-user-id",
		Email:            "test@example.com",
		PasswordHash:     "hashed-password",
		Name:             "Test User",
		Role:             domain.RoleEmployee,
		VacationBalance:  25,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	// Create
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Get by ID
	found, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Name, found.Name)

	// Get by Email
	found, err = repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_EmailExists(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	ctx := context.Background()

	user := &domain.User{
		ID:           "test-user-id",
		Email:        "existing@example.com",
		PasswordHash: "hash",
		Name:         "Test",
		Role:         domain.RoleEmployee,
	}
	repo.Create(ctx, user)

	exists, err := repo.EmailExists(ctx, "existing@example.com", "")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.EmailExists(ctx, "nonexistent@example.com", "")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Exclude self
	exists, err = repo.EmailExists(ctx, "existing@example.com", user.ID)
	assert.NoError(t, err)
	assert.False(t, exists)
}
```

#### Handler Tests

- [ ] **Auth Handler Tests** `internal/handler/auth_test.go`

```go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Setup
	// ... mock services

	r := setupTestRouter()
	// h := NewAuthHandler(mockAuthService, mockUserRepo)
	// r.POST("/api/auth/login", h.Login)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
	assert.NotNil(t, response["user"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Similar test for 401 response
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	r := setupTestRouter()
	// ... setup

	body := map[string]string{
		"email": "invalid-email",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

---

## Frontend Testing

### Unit Tests (Vitest)

- [ ] **Setup Vitest** `vitest.config.ts`

```typescript
import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
  plugins: [sveltekit()],
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/tests/setup.ts'],
  },
});
```

- [ ] **Test Setup** `src/tests/setup.ts`

```typescript
import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
vi.stubGlobal('localStorage', localStorageMock);

// Mock fetch
vi.stubGlobal('fetch', vi.fn());
```

### Store Tests

- [ ] **Auth Store Tests** `src/lib/stores/auth.test.ts`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { auth } from './auth.svelte';

describe('auth store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('starts with null user', () => {
    expect(auth.user).toBeNull();
    expect(auth.isAuthenticated).toBe(false);
  });

  it('login sets user and token', async () => {
    const mockUser = {
      id: '123',
      email: 'test@example.com',
      name: 'Test User',
      role: 'employee' as const,
      vacationBalance: 25,
      emailPreferences: {
        vacationUpdates: true,
        weeklyDigest: false,
        teamNotifications: true,
      },
    };

    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ token: 'jwt-token', user: mockUser }),
    } as Response);

    await auth.login('test@example.com', 'password');

    expect(auth.user).toEqual(mockUser);
    expect(auth.isAuthenticated).toBe(true);
    expect(localStorage.setItem).toHaveBeenCalledWith('auth_token', 'jwt-token');
  });

  it('logout clears user and token', () => {
    auth.logout();

    expect(auth.user).toBeNull();
    expect(auth.isAuthenticated).toBe(false);
    expect(localStorage.removeItem).toHaveBeenCalledWith('auth_token');
  });
});
```

### Component Tests

- [ ] **Button Component Tests** `src/lib/components/ui/Button.test.ts`

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import Button from './Button.svelte';

describe('Button', () => {
  it('renders with default props', () => {
    render(Button, { props: { children: 'Click me' } });
    expect(screen.getByRole('button')).toBeInTheDocument();
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('applies variant classes', () => {
    render(Button, { props: { variant: 'primary', children: 'Primary' } });
    const button = screen.getByRole('button');
    expect(button).toHaveClass('bg-ocean-500');
  });

  it('shows loading spinner when loading', () => {
    render(Button, { props: { loading: true, children: 'Loading' } });
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
    expect(button.querySelector('svg')).toBeInTheDocument();
  });

  it('calls onclick handler', async () => {
    const onclick = vi.fn();
    render(Button, { props: { onclick, children: 'Click' } });

    await fireEvent.click(screen.getByRole('button'));
    expect(onclick).toHaveBeenCalledOnce();
  });

  it('does not call onclick when disabled', async () => {
    const onclick = vi.fn();
    render(Button, { props: { onclick, disabled: true, children: 'Click' } });

    await fireEvent.click(screen.getByRole('button'));
    expect(onclick).not.toHaveBeenCalled();
  });
});
```

- [ ] **Input Component Tests** `src/lib/components/ui/Input.test.ts`

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import Input from './Input.svelte';

describe('Input', () => {
  it('renders with label', () => {
    render(Input, { props: { label: 'Email' } });
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
  });

  it('shows error message', () => {
    render(Input, { props: { error: 'Invalid input' } });
    expect(screen.getByText('Invalid input')).toBeInTheDocument();
  });

  it('binds value', async () => {
    const { component } = render(Input, { props: { value: '' } });

    const input = screen.getByRole('textbox');
    await fireEvent.input(input, { target: { value: 'test' } });

    expect(input).toHaveValue('test');
  });

  it('shows required indicator', () => {
    render(Input, { props: { label: 'Name', required: true } });
    expect(screen.getByText('*')).toBeInTheDocument();
  });
});
```

### API Client Tests

- [ ] **API Client Tests** `src/lib/api/client.test.ts`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { request, setAuthToken, clearAuthToken, ApiException } from './client';

describe('API client', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('makes successful request', async () => {
    const mockData = { id: '123', name: 'Test' };
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: () => Promise.resolve(mockData),
    } as Response);

    const result = await request('/test');
    expect(result).toEqual(mockData);
  });

  it('includes auth token when present', async () => {
    vi.mocked(localStorage.getItem).mockReturnValue('test-token');
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({}),
    } as Response);

    await request('/test');

    expect(fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer test-token',
        }),
      })
    );
  });

  it('throws ApiException on error response', async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      status: 401,
      json: () => Promise.resolve({ code: 'UNAUTHORIZED', message: 'Invalid token' }),
    } as Response);

    await expect(request('/test')).rejects.toThrow(ApiException);
  });

  it('handles 204 No Content', async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: true,
      status: 204,
    } as Response);

    const result = await request('/test', { method: 'DELETE' });
    expect(result).toBeUndefined();
  });
});
```

---

## Integration Testing

### API Integration Tests

- [ ] **Full Flow Test** `internal/integration/vacation_flow_test.go`

```go
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVacationRequestFlow(t *testing.T) {
	// Setup full application
	app := setupTestApp(t)
	defer app.Cleanup()

	// 1. Login as employee
	loginResp := app.POST("/api/auth/login", map[string]string{
		"email":    "employee@test.com",
		"password": "password123",
	})
	require.Equal(t, http.StatusOK, loginResp.Code)

	var loginData map[string]interface{}
	json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["token"].(string)

	// 2. Create vacation request
	createResp := app.AuthPOST(token, "/api/vacation/request", map[string]string{
		"startDate": "15/01/2024",
		"endDate":   "19/01/2024",
		"reason":    "Vacation",
	})
	require.Equal(t, http.StatusCreated, createResp.Code)

	var request map[string]interface{}
	json.Unmarshal(createResp.Body.Bytes(), &request)
	requestID := request["id"].(string)

	// 3. Verify request appears in list
	listResp := app.AuthGET(token, "/api/vacation/requests")
	require.Equal(t, http.StatusOK, listResp.Code)

	var listData map[string]interface{}
	json.Unmarshal(listResp.Body.Bytes(), &listData)
	requests := listData["requests"].([]interface{})
	assert.Len(t, requests, 1)

	// 4. Login as admin
	adminLoginResp := app.POST("/api/auth/login", map[string]string{
		"email":    "admin@test.com",
		"password": "admin123",
	})
	var adminLoginData map[string]interface{}
	json.Unmarshal(adminLoginResp.Body.Bytes(), &adminLoginData)
	adminToken := adminLoginData["token"].(string)

	// 5. Approve request
	approveResp := app.AuthPUT(adminToken, "/api/admin/vacation/"+requestID+"/approve", nil)
	require.Equal(t, http.StatusOK, approveResp.Code)

	// 6. Verify balance updated
	meResp := app.AuthGET(token, "/api/auth/me")
	var meData map[string]interface{}
	json.Unmarshal(meResp.Body.Bytes(), &meData)
	// Original balance was 25, 5 days deducted
	assert.Equal(t, float64(20), meData["vacationBalance"])
}
```

---

## E2E Testing

### Playwright Setup

- [ ] **Playwright Config** `playwright.config.ts`

```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

### E2E Tests

- [ ] **Login Flow** `tests/e2e/auth.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
  test('successful login redirects to dashboard', async ({ page }) => {
    await page.goto('/');

    await page.fill('input[type="email"]', 'employee@test.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/employee');
    await expect(page.locator('h1')).toContainText('Welcome back');
  });

  test('invalid credentials shows error', async ({ page }) => {
    await page.goto('/');

    await page.fill('input[type="email"]', 'wrong@test.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    await page.click('button[type="submit"]');

    await expect(page.locator('[role="alert"]')).toBeVisible();
  });

  test('admin login redirects to admin dashboard', async ({ page }) => {
    await page.goto('/');

    await page.fill('input[type="email"]', 'admin@test.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');

    await expect(page).toHaveURL('/admin');
  });
});
```

- [ ] **Vacation Request Flow** `tests/e2e/vacation.spec.ts`

```typescript
import { test, expect } from '@playwright/test';

test.describe('Vacation Requests', () => {
  test.beforeEach(async ({ page }) => {
    // Login as employee
    await page.goto('/');
    await page.fill('input[type="email"]', 'employee@test.com');
    await page.fill('input[type="password"]', 'password123');
    await page.click('button[type="submit"]');
    await page.waitForURL('/employee');
  });

  test('can submit vacation request', async ({ page }) => {
    // Open request modal
    await page.click('text=Request Vacation');

    // Fill form
    await page.fill('input[type="date"]:first-of-type', '2024-01-15');
    await page.fill('input[type="date"]:last-of-type', '2024-01-19');
    await page.fill('textarea', 'Family vacation');

    // Submit
    await page.click('text=Submit Request');

    // Verify success
    await expect(page.locator('[role="alert"]')).toContainText('submitted');
    await expect(page.locator('text=15 Jan - 19 Jan')).toBeVisible();
  });

  test('can cancel pending request', async ({ page }) => {
    // Find and cancel request
    await page.click('[data-testid="cancel-request"]');
    await page.click('text=Confirm');

    await expect(page.locator('[role="alert"]')).toContainText('cancelled');
  });
});
```

---

## Test Data & Fixtures

### Backend Test Fixtures

```go
// internal/testutil/fixtures.go
package testutil

import (
	"github.com/yourorg/vacaytracker-api/internal/domain"
)

func TestUser(overrides ...func(*domain.User)) *domain.User {
	user := &domain.User{
		ID:               "test-user-id",
		Email:            "test@example.com",
		PasswordHash:     "$2a$10$...", // bcrypt hash of "password123"
		Name:             "Test User",
		Role:             domain.RoleEmployee,
		VacationBalance:  25,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}
	for _, override := range overrides {
		override(user)
	}
	return user
}

func TestAdmin(overrides ...func(*domain.User)) *domain.User {
	user := TestUser(func(u *domain.User) {
		u.ID = "test-admin-id"
		u.Email = "admin@example.com"
		u.Name = "Test Admin"
		u.Role = domain.RoleAdmin
	})
	for _, override := range overrides {
		override(user)
	}
	return user
}

func TestVacationRequest(overrides ...func(*domain.VacationRequest)) *domain.VacationRequest {
	request := &domain.VacationRequest{
		ID:        "test-request-id",
		UserID:    "test-user-id",
		StartDate: "2024-01-15",
		EndDate:   "2024-01-19",
		TotalDays: 5,
		Status:    domain.StatusPending,
	}
	for _, override := range overrides {
		override(request)
	}
	return request
}
```

### Frontend Test Fixtures

```typescript
// src/tests/fixtures.ts
import type { User, VacationRequest } from '$lib/types';

export function testUser(overrides: Partial<User> = {}): User {
  return {
    id: 'test-user-id',
    email: 'test@example.com',
    name: 'Test User',
    role: 'employee',
    vacationBalance: 25,
    emailPreferences: {
      vacationUpdates: true,
      weeklyDigest: false,
      teamNotifications: true,
    },
    ...overrides,
  };
}

export function testVacationRequest(overrides: Partial<VacationRequest> = {}): VacationRequest {
  return {
    id: 'test-request-id',
    userId: 'test-user-id',
    startDate: '2024-01-15',
    endDate: '2024-01-19',
    totalDays: 5,
    status: 'pending',
    createdAt: '2024-01-10T10:00:00Z',
    updatedAt: '2024-01-10T10:00:00Z',
    ...overrides,
  };
}
```

---

## Coverage Requirements

### Backend Coverage Targets

| Layer | Target | Critical Files |
|-------|--------|----------------|
| Services | 90% | auth.go, vacation.go, user.go |
| Handlers | 80% | auth.go, vacation.go, admin.go |
| Repositories | 75% | user.go, vacation.go |
| Middleware | 85% | auth.go |
| Domain | 70% | Models, helpers |

### Frontend Coverage Targets

| Category | Target | Critical Files |
|----------|--------|----------------|
| Stores | 90% | auth.svelte.ts, vacation.svelte.ts |
| API Modules | 85% | client.ts, auth.ts |
| UI Components | 70% | Button, Input, Card |
| Feature Components | 75% | RequestModal, PendingRequests |

### Running Coverage

**Backend:**
```bash
make test-coverage
# Output: coverage.html
```

**Frontend:**
```bash
npm run test:coverage
# Output: coverage/index.html
```

---

## Related Documents

- [04-backend-tasks.md](./04-backend-tasks.md) - Backend implementation
- [05-frontend-tasks.md](./05-frontend-tasks.md) - Frontend implementation
- [10-development-workflow.md](./10-development-workflow.md) - CI/CD integration
