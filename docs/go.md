# VacayTracker Backend Implementation Guide (Go)

> **Complete technical guide for building the VacayTracker API with Go**
>
> **Related Documentation**
> - [Features Documentation](./FEATURES.md) — Complete feature specifications
> - [Frontend Implementation Guide](./vacaytracker-frontend-guide.md) — Svelte 5 + Melt UI technical guide

---

## 1. Technology Stack

### 1.1 Core Technologies

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| Language | Go | 1.23+ | Core runtime |
| Web Framework | Gin | v1.10+ | HTTP routing, middleware |
| Database | SQLite | - | Embedded database |
| SQLite Driver | modernc.org/sqlite | v1.34+ | CGo-free SQLite driver |
| JWT | golang-jwt/jwt | v5 | Authentication tokens |
| Password Hashing | golang.org/x/crypto | - | bcrypt implementation |
| Email | resend-go | v2+ | Transactional emails |
| Configuration | godotenv | v1.5+ | Environment variables |
| Validation | go-playground/validator | v10 | Input validation |
| UUID | google/uuid | v1.6+ | ID generation |

### 1.2 Why These Choices?

**Gin** — Most popular Go web framework (75k+ GitHub stars), excellent performance, mature ecosystem, and gentle learning curve. Perfect for REST APIs.

**modernc.org/sqlite** — Pure Go SQLite implementation (no CGo). Enables cross-compilation, simpler deployment, and avoids CGo complexity. For a 5-person team app, SQLite provides simplicity without infrastructure overhead.

**golang-jwt/jwt v5** — The maintained successor to dgrijalva/jwt-go. Production-ready, supports all standard signing algorithms, and has excellent documentation.

**godotenv** — Simple and familiar pattern from other languages. Loads `.env` files for local development without the complexity of Viper (which would be overkill for this project).

---

## 2. Project Setup

### 2.1 Initialize Project

```bash
# Create project directory
mkdir vacaytracker-api
cd vacaytracker-api

# Initialize Go module
go mod init github.com/yourorg/vacaytracker-api

# Install dependencies
go get github.com/gin-gonic/gin
go get modernc.org/sqlite
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get github.com/google/uuid
go get github.com/resend/resend-go/v2
```

### 2.2 Project Structure

```
vacaytracker-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── domain/
│   │   ├── user.go              # User entity
│   │   ├── vacation.go          # Vacation request entity
│   │   └── settings.go          # Settings entity
│   ├── repository/
│   │   ├── repository.go        # Repository interfaces
│   │   ├── sqlite/
│   │   │   ├── sqlite.go        # SQLite connection
│   │   │   ├── user.go          # User repository
│   │   │   ├── vacation.go      # Vacation repository
│   │   │   └── settings.go      # Settings repository
│   ├── service/
│   │   ├── auth.go              # Authentication service
│   │   ├── user.go              # User business logic
│   │   ├── vacation.go          # Vacation business logic
│   │   └── email.go             # Email service
│   ├── handler/
│   │   ├── handler.go           # Handler dependencies
│   │   ├── auth.go              # Auth endpoints
│   │   ├── user.go              # User endpoints
│   │   ├── vacation.go          # Vacation endpoints
│   │   ├── admin.go             # Admin endpoints
│   │   └── health.go            # Health check
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── cors.go              # CORS configuration
│   │   └── error.go             # Error handling
│   └── dto/
│       ├── request.go           # Request DTOs
│       ├── response.go          # Response DTOs
│       └── errors.go            # Error types
├── pkg/
│   └── validator/
│       └── validator.go         # Custom validators
├── migrations/
│   └── 001_init.sql             # Database schema
├── data/
│   └── .gitkeep                 # Database directory
├── .env.example
├── .env
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 2.3 Makefile

```makefile
.PHONY: run build test clean migrate

# Development
run:
	go run cmd/server/main.go

# Build
build:
	CGO_ENABLED=0 go build -o bin/vacaytracker cmd/server/main.go

# Testing
test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Database
migrate:
	sqlite3 data/vacaytracker.db < migrations/001_init.sql

# Lint
lint:
	golangci-lint run

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out
```

---

## 3. Configuration

### 3.1 Environment Variables

Create `.env.example`:

```bash
# ===================
# Required
# ===================
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
ADMIN_PASSWORD=ChangeThisSecurePassword!

# ===================
# Email (Resend)
# ===================
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxx
EMAIL_FROM_ADDRESS=vacaytracker@yourcompany.com
EMAIL_FROM_NAME=VacayTracker

# ===================
# Server
# ===================
ENV=development
PORT=3000
APP_URL=http://localhost:3000

# ===================
# Database
# ===================
DB_PATH=./data/vacaytracker.db

# ===================
# Security
# ===================
TOKEN_EXPIRY=24h
BCRYPT_COST=10
CORS_ORIGINS=http://localhost:5173

# ===================
# Newsletter
# ===================
NEWSLETTER_DAY=1
NEWSLETTER_HOUR=9
```

### 3.2 Configuration Struct

`internal/config/config.go`:

```go
package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Env     string
	Port    string
	AppURL  string

	// Database
	DBPath string

	// Security
	JWTSecret     string
	TokenExpiry   time.Duration
	BcryptCost    int
	CORSOrigins   []string
	AdminPassword string

	// Email
	ResendAPIKey     string
	EmailFromAddress string
	EmailFromName    string

	// Newsletter
	NewsletterDay  int
	NewsletterHour int
}

func Load() *Config {
	// Load .env file (ignore error if not exists - production uses real env vars)
	_ = godotenv.Load()

	tokenExpiry, _ := time.ParseDuration(getEnv("TOKEN_EXPIRY", "24h"))
	bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "10"))
	newsletterDay, _ := strconv.Atoi(getEnv("NEWSLETTER_DAY", "1"))
	newsletterHour, _ := strconv.Atoi(getEnv("NEWSLETTER_HOUR", "9"))

	cfg := &Config{
		Env:     getEnv("ENV", "development"),
		Port:    getEnv("PORT", "3000"),
		AppURL:  getEnv("APP_URL", "http://localhost:3000"),

		DBPath: getEnv("DB_PATH", "./data/vacaytracker.db"),

		JWTSecret:     requireEnv("JWT_SECRET"),
		TokenExpiry:   tokenExpiry,
		BcryptCost:    bcryptCost,
		CORSOrigins:   strings.Split(getEnv("CORS_ORIGINS", "*"), ","),
		AdminPassword: requireEnv("ADMIN_PASSWORD"),

		ResendAPIKey:     getEnv("RESEND_API_KEY", ""),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", "noreply@example.com"),
		EmailFromName:    getEnv("EMAIL_FROM_NAME", "VacayTracker"),

		NewsletterDay:  newsletterDay,
		NewsletterHour: newsletterHour,
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsEmailEnabled returns true if email service is configured
func (c *Config) IsEmailEnabled() bool {
	return c.ResendAPIKey != ""
}
```

---

## 4. Domain Models

### 4.1 User Entity

`internal/domain/user.go`:

```go
package domain

import (
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type EmailPreferences struct {
	Enabled                      bool `json:"enabled"`
	VacationStatusUpdates        bool `json:"vacationStatusUpdates,omitempty"`        // Employee
	VacationRequestNotifications bool `json:"vacationRequestNotifications,omitempty"` // Admin
	UserCreatedNotifications     bool `json:"userCreatedNotifications,omitempty"`     // Admin
	VacationResetNotifications   bool `json:"vacationResetNotifications"`
	MonthlyVacationSummary       bool `json:"monthlyVacationSummary"`
}

type User struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Username         string           `json:"username"`
	Email            string           `json:"email,omitempty"`
	Password         string           `json:"-"` // Never serialize password
	Role             Role             `json:"role"`
	VacationDays     int              `json:"vacationDays,omitempty"`
	UsedVacationDays int              `json:"usedVacationDays,omitempty"`
	EmailPreferences EmailPreferences `json:"emailPreferences"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// RemainingDays calculates available vacation days
func (u *User) RemainingDays() int {
	return u.VacationDays - u.UsedVacationDays
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// DefaultEmployeeEmailPreferences returns default preferences for employees
func DefaultEmployeeEmailPreferences() EmailPreferences {
	return EmailPreferences{
		Enabled:                    true,
		VacationStatusUpdates:      true,
		VacationResetNotifications: true,
		MonthlyVacationSummary:     true,
	}
}

// DefaultAdminEmailPreferences returns default preferences for admins
func DefaultAdminEmailPreferences() EmailPreferences {
	return EmailPreferences{
		Enabled:                      true,
		VacationRequestNotifications: true,
		UserCreatedNotifications:     true,
		VacationResetNotifications:   true,
		MonthlyVacationSummary:       true,
	}
}
```

### 4.2 Vacation Request Entity

`internal/domain/vacation.go`:

```go
package domain

import (
	"time"
)

type VacationStatus string

const (
	StatusPending  VacationStatus = "pending"
	StatusApproved VacationStatus = "approved"
	StatusRejected VacationStatus = "rejected"
)

type VacationRequest struct {
	ID           string         `json:"id"`
	UserID       string         `json:"userId"`
	StartDate    string         `json:"startDate"`    // YYYY-MM-DD
	EndDate      string         `json:"endDate"`      // YYYY-MM-DD
	BusinessDays int            `json:"businessDays"`
	Status       VacationStatus `json:"status"`
	Reason       string         `json:"reason,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	ReviewedBy   *string        `json:"reviewedBy,omitempty"`
	ReviewedAt   *time.Time     `json:"reviewedAt,omitempty"`
}

// IsPending checks if request is still pending
func (v *VacationRequest) IsPending() bool {
	return v.Status == StatusPending
}

// CanBeReviewed checks if request can be approved/rejected
func (v *VacationRequest) CanBeReviewed() bool {
	return v.Status == StatusPending
}
```

### 4.3 Settings Entity

`internal/domain/settings.go`:

```go
package domain

import "time"

type WeekendPolicy struct {
	ExcludeWeekends bool `json:"excludeWeekends"`
}

type NewsletterSettings struct {
	Enabled    bool       `json:"enabled"`
	DayOfMonth int        `json:"dayOfMonth"`
	HourOfDay  int        `json:"hourOfDay"`
	LastSentAt *time.Time `json:"lastSentAt,omitempty"`
}

type Settings struct {
	WeekendPolicy WeekendPolicy      `json:"weekendPolicy"`
	Newsletter    NewsletterSettings `json:"newsletter"`
}

// DefaultSettings returns initial settings
func DefaultSettings() Settings {
	return Settings{
		WeekendPolicy: WeekendPolicy{
			ExcludeWeekends: true,
		},
		Newsletter: NewsletterSettings{
			Enabled:    false,
			DayOfMonth: 1,
			HourOfDay:  9,
		},
	}
}
```

---

## 5. Database Layer

### 5.1 Schema Migration

`migrations/001_init.sql`:

```sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'employee')),
    vacation_days INTEGER DEFAULT 0,
    used_vacation_days INTEGER DEFAULT 0,
    email_preferences TEXT NOT NULL DEFAULT '{}',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Vacation requests table
CREATE TABLE IF NOT EXISTS vacation_requests (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    business_days INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    reason TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_by TEXT,
    reviewed_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Settings table (single row)
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    weekend_policy TEXT NOT NULL DEFAULT '{"excludeWeekends": true}',
    newsletter TEXT NOT NULL DEFAULT '{"enabled": false, "dayOfMonth": 1, "hourOfDay": 9}'
);

-- Insert default settings
INSERT OR IGNORE INTO settings (id) VALUES (1);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_vacation_user ON vacation_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_vacation_status ON vacation_requests(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Trigger to update updated_at
CREATE TRIGGER IF NOT EXISTS update_user_timestamp 
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

### 5.2 SQLite Connection

`internal/repository/sqlite/sqlite.go`:

```go
package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database with WAL mode and other optimizations
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", dbPath)
	
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite
	db.SetMaxOpenConns(1)  // SQLite is single-writer
	db.SetMaxIdleConns(1)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Migrate runs database migrations
func (db *DB) Migrate(migrationPath string) error {
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
```

### 5.3 Repository Interfaces

`internal/repository/repository.go`:

```go
package repository

import (
	"context"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByRole(ctx context.Context, role domain.Role) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, id, hashedPassword string) error
	UpdateVacationDays(ctx context.Context, id string, used int) error
	Delete(ctx context.Context, id string) error
	UsernameExists(ctx context.Context, username string, excludeID string) (bool, error)
	ResetAllVacationDays(ctx context.Context, newDays int) error
}

type VacationRepository interface {
	Create(ctx context.Context, request *domain.VacationRequest) error
	GetByID(ctx context.Context, id string) (*domain.VacationRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.VacationRequest, error)
	GetAll(ctx context.Context) ([]domain.VacationRequest, error)
	GetPending(ctx context.Context) ([]domain.VacationRequest, error)
	Update(ctx context.Context, request *domain.VacationRequest) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	HasOverlap(ctx context.Context, userID, startDate, endDate string, excludeID string) (bool, error)
}

type SettingsRepository interface {
	Get(ctx context.Context) (*domain.Settings, error)
	Update(ctx context.Context, settings *domain.Settings) error
}
```

### 5.4 User Repository Implementation

`internal/repository/sqlite/user.go`:

```go
package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	prefs, err := json.Marshal(user.EmailPreferences)
	if err != nil {
		return fmt.Errorf("failed to marshal email preferences: %w", err)
	}

	query := `
		INSERT INTO users (id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err = r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.VacationDays,
		user.UsedVacationDays,
		string(prefs),
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at
		FROM users WHERE id = ?
	`
	
	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at
		FROM users WHERE username = ?
	`
	
	return r.scanUser(r.db.QueryRowContext(ctx, query, username))
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`
	
	return r.scanUsers(r.db.QueryContext(ctx, query))
}

func (r *UserRepository) GetByRole(ctx context.Context, role domain.Role) ([]domain.User, error) {
	query := `
		SELECT id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at
		FROM users WHERE role = ? ORDER BY created_at DESC
	`
	
	return r.scanUsers(r.db.QueryContext(ctx, query, role))
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	prefs, err := json.Marshal(user.EmailPreferences)
	if err != nil {
		return fmt.Errorf("failed to marshal email preferences: %w", err)
	}

	query := `
		UPDATE users 
		SET name = ?, username = ?, email = ?, role = ?, vacation_days = ?, used_vacation_days = ?, email_preferences = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Username,
		user.Email,
		user.Role,
		user.VacationDays,
		user.UsedVacationDays,
		string(prefs),
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, hashedPassword string) error {
	query := `UPDATE users SET password = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, id)
	return err
}

func (r *UserRepository) UpdateVacationDays(ctx context.Context, id string, used int) error {
	query := `UPDATE users SET used_vacation_days = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, used, id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *UserRepository) UsernameExists(ctx context.Context, username string, excludeID string) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != "" {
		query = `SELECT COUNT(*) FROM users WHERE username = ? AND id != ?`
		args = []interface{}{username, excludeID}
	} else {
		query = `SELECT COUNT(*) FROM users WHERE username = ?`
		args = []interface{}{username}
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) ResetAllVacationDays(ctx context.Context, newDays int) error {
	query := `UPDATE users SET vacation_days = ?, used_vacation_days = 0 WHERE role = 'employee'`
	_, err := r.db.ExecContext(ctx, query, newDays)
	return err
}

// Helper methods

func (r *UserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var email sql.NullString
	var prefsJSON string

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&email,
		&user.Password,
		&user.Role,
		&user.VacationDays,
		&user.UsedVacationDays,
		&prefsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user.Email = email.String
	_ = json.Unmarshal([]byte(prefsJSON), &user.EmailPreferences)
	
	return &user, nil
}

func (r *UserRepository) scanUsers(rows *sql.Rows, err error) ([]domain.User, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		var email sql.NullString
		var prefsJSON string

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&email,
			&user.Password,
			&user.Role,
			&user.VacationDays,
			&user.UsedVacationDays,
			&prefsJSON,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		user.Email = email.String
		_ = json.Unmarshal([]byte(prefsJSON), &user.EmailPreferences)
		users = append(users, user)
	}

	return users, rows.Err()
}
```

### 5.5 Vacation Repository Implementation

`internal/repository/sqlite/vacation.go`:

```go
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

type VacationRepository struct {
	db *DB
}

func NewVacationRepository(db *DB) *VacationRepository {
	return &VacationRepository{db: db}
}

func (r *VacationRepository) Create(ctx context.Context, req *domain.VacationRequest) error {
	query := `
		INSERT INTO vacation_requests (id, user_id, start_date, end_date, business_days, status, reason, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		req.ID,
		req.UserID,
		req.StartDate,
		req.EndDate,
		req.BusinessDays,
		req.Status,
		req.Reason,
		req.CreatedAt,
	)
	
	return err
}

func (r *VacationRepository) GetByID(ctx context.Context, id string) (*domain.VacationRequest, error) {
	query := `
		SELECT id, user_id, start_date, end_date, business_days, status, reason, created_at, reviewed_by, reviewed_at
		FROM vacation_requests WHERE id = ?
	`
	
	return r.scanRequest(r.db.QueryRowContext(ctx, query, id))
}

func (r *VacationRepository) GetByUserID(ctx context.Context, userID string) ([]domain.VacationRequest, error) {
	query := `
		SELECT id, user_id, start_date, end_date, business_days, status, reason, created_at, reviewed_by, reviewed_at
		FROM vacation_requests WHERE user_id = ? ORDER BY created_at DESC
	`
	
	return r.scanRequests(r.db.QueryContext(ctx, query, userID))
}

func (r *VacationRepository) GetAll(ctx context.Context) ([]domain.VacationRequest, error) {
	query := `
		SELECT id, user_id, start_date, end_date, business_days, status, reason, created_at, reviewed_by, reviewed_at
		FROM vacation_requests ORDER BY created_at DESC
	`
	
	return r.scanRequests(r.db.QueryContext(ctx, query))
}

func (r *VacationRepository) GetPending(ctx context.Context) ([]domain.VacationRequest, error) {
	query := `
		SELECT id, user_id, start_date, end_date, business_days, status, reason, created_at, reviewed_by, reviewed_at
		FROM vacation_requests WHERE status = 'pending' ORDER BY created_at ASC
	`
	
	return r.scanRequests(r.db.QueryContext(ctx, query))
}

func (r *VacationRepository) Update(ctx context.Context, req *domain.VacationRequest) error {
	query := `
		UPDATE vacation_requests 
		SET status = ?, reviewed_by = ?, reviewed_at = ?
		WHERE id = ?
	`
	
	result, err := r.db.ExecContext(ctx, query,
		req.Status,
		req.ReviewedBy,
		req.ReviewedAt,
		req.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *VacationRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM vacation_requests WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *VacationRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM vacation_requests WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *VacationRepository) HasOverlap(ctx context.Context, userID, startDate, endDate string, excludeID string) (bool, error) {
	var query string
	var args []interface{}

	if excludeID != "" {
		query = `
			SELECT COUNT(*) FROM vacation_requests 
			WHERE user_id = ? AND id != ? AND status != 'rejected'
			AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?))
		`
		args = []interface{}{userID, excludeID, endDate, startDate, startDate, startDate, startDate, endDate}
	} else {
		query = `
			SELECT COUNT(*) FROM vacation_requests 
			WHERE user_id = ? AND status != 'rejected'
			AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?))
		`
		args = []interface{}{userID, endDate, startDate, startDate, startDate, startDate, endDate}
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count > 0, err
}

// Helper methods

func (r *VacationRepository) scanRequest(row *sql.Row) (*domain.VacationRequest, error) {
	var req domain.VacationRequest
	var reason sql.NullString
	var reviewedBy sql.NullString
	var reviewedAt sql.NullTime

	err := row.Scan(
		&req.ID,
		&req.UserID,
		&req.StartDate,
		&req.EndDate,
		&req.BusinessDays,
		&req.Status,
		&reason,
		&req.CreatedAt,
		&reviewedBy,
		&reviewedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	req.Reason = reason.String
	if reviewedBy.Valid {
		req.ReviewedBy = &reviewedBy.String
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}
	
	return &req, nil
}

func (r *VacationRepository) scanRequests(rows *sql.Rows, err error) ([]domain.VacationRequest, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.VacationRequest
	for rows.Next() {
		var req domain.VacationRequest
		var reason sql.NullString
		var reviewedBy sql.NullString
		var reviewedAt sql.NullTime

		err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.StartDate,
			&req.EndDate,
			&req.BusinessDays,
			&req.Status,
			&reason,
			&req.CreatedAt,
			&reviewedBy,
			&reviewedAt,
		)
		if err != nil {
			return nil, err
		}

		req.Reason = reason.String
		if reviewedBy.Valid {
			req.ReviewedBy = &reviewedBy.String
		}
		if reviewedAt.Valid {
			req.ReviewedAt = &reviewedAt.Time
		}
		requests = append(requests, req)
	}

	return requests, rows.Err()
}
```

---

## 6. DTOs and Error Handling

### 6.1 Request DTOs

`internal/dto/request.go`:

```go
package dto

// Authentication
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=1"`
}

// User Management
type CreateUserRequest struct {
	Name         string `json:"name" binding:"required,min=1,max=100"`
	Username     string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email        string `json:"email" binding:"omitempty,email"`
	Password     string `json:"password" binding:"required,min=6,max=72"`
	Role         string `json:"role" binding:"required,oneof=admin employee"`
	VacationDays int    `json:"vacationDays" binding:"omitempty,min=0,max=365"`
}

type UpdateUserRequest struct {
	Name         string `json:"name" binding:"omitempty,min=1,max=100"`
	Username     string `json:"username" binding:"omitempty,min=3,max=50,alphanum"`
	Email        string `json:"email" binding:"omitempty,email"`
	Role         string `json:"role" binding:"omitempty,oneof=admin employee"`
	VacationDays *int   `json:"vacationDays" binding:"omitempty,min=0,max=365"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=72"`
}

type UpdateEmailPreferencesRequest struct {
	Enabled                      *bool `json:"enabled"`
	VacationStatusUpdates        *bool `json:"vacationStatusUpdates,omitempty"`
	VacationRequestNotifications *bool `json:"vacationRequestNotifications,omitempty"`
	UserCreatedNotifications     *bool `json:"userCreatedNotifications,omitempty"`
	VacationResetNotifications   *bool `json:"vacationResetNotifications,omitempty"`
	MonthlyVacationSummary       *bool `json:"monthlyVacationSummary,omitempty"`
}

// Vacation Requests
type CreateVacationRequest struct {
	StartDate string `json:"startDate" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"endDate" binding:"required,datetime=2006-01-02"`
	Reason    string `json:"reason" binding:"omitempty,max=200"`
}

type ReviewVacationRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

// Admin
type ResetVacationDaysRequest struct {
	Days int `json:"days" binding:"required,min=0,max=365"`
}

type UpdateSettingsRequest struct {
	WeekendPolicy *struct {
		ExcludeWeekends *bool `json:"excludeWeekends"`
	} `json:"weekendPolicy"`
	Newsletter *struct {
		Enabled    *bool `json:"enabled"`
		DayOfMonth *int  `json:"dayOfMonth" binding:"omitempty,min=1,max=31"`
		HourOfDay  *int  `json:"hourOfDay" binding:"omitempty,min=0,max=23"`
	} `json:"newsletter"`
}
```

### 6.2 Response DTOs

`internal/dto/response.go`:

```go
package dto

import "time"

// Standard response wrappers
type DataResponse struct {
	Data interface{} `json:"data"`
}

type DataWithMessageResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type ListResponse struct {
	Data interface{} `json:"data"`
	Meta MetaInfo    `json:"meta"`
}

type MetaInfo struct {
	Total int `json:"total"`
	Count int `json:"count"`
}

// Authentication
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// User
type UserResponse struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Username         string                   `json:"username"`
	Email            string                   `json:"email,omitempty"`
	Role             string                   `json:"role"`
	VacationDays     int                      `json:"vacationDays,omitempty"`
	UsedVacationDays int                      `json:"usedVacationDays,omitempty"`
	RemainingDays    int                      `json:"remainingDays,omitempty"`
	EmailPreferences EmailPreferencesResponse `json:"emailPreferences"`
	CreatedAt        time.Time                `json:"createdAt"`
}

type EmailPreferencesResponse struct {
	Enabled                      bool `json:"enabled"`
	VacationStatusUpdates        bool `json:"vacationStatusUpdates,omitempty"`
	VacationRequestNotifications bool `json:"vacationRequestNotifications,omitempty"`
	UserCreatedNotifications     bool `json:"userCreatedNotifications,omitempty"`
	VacationResetNotifications   bool `json:"vacationResetNotifications"`
	MonthlyVacationSummary       bool `json:"monthlyVacationSummary"`
}

// Vacation
type VacationRequestResponse struct {
	ID           string     `json:"id"`
	UserID       string     `json:"userId"`
	UserName     string     `json:"userName,omitempty"`
	StartDate    string     `json:"startDate"`
	EndDate      string     `json:"endDate"`
	BusinessDays int        `json:"businessDays"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	ReviewedBy   *string    `json:"reviewedBy,omitempty"`
	ReviewedAt   *time.Time `json:"reviewedAt,omitempty"`
}

// Settings
type SettingsResponse struct {
	WeekendPolicy struct {
		ExcludeWeekends bool `json:"excludeWeekends"`
	} `json:"weekendPolicy"`
	Newsletter struct {
		Enabled    bool       `json:"enabled"`
		DayOfMonth int        `json:"dayOfMonth"`
		HourOfDay  int        `json:"hourOfDay"`
		LastSentAt *time.Time `json:"lastSentAt,omitempty"`
	} `json:"newsletter"`
}

// Health
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version,omitempty"`
}
```

### 6.3 Error Types

`internal/dto/errors.go`:

```go
package dto

import (
	"fmt"
	"net/http"
)

// ErrorCode represents machine-readable error codes
type ErrorCode string

const (
	// Authentication Errors (401)
	ErrAuthTokenMissing    ErrorCode = "AUTH_TOKEN_MISSING"
	ErrAuthTokenInvalid    ErrorCode = "AUTH_TOKEN_INVALID"
	ErrAuthTokenExpired    ErrorCode = "AUTH_TOKEN_EXPIRED"
	ErrAuthCredentials     ErrorCode = "AUTH_CREDENTIALS_INVALID"
	ErrAuthRoleMismatch    ErrorCode = "AUTH_ROLE_MISMATCH"

	// Authorization Errors (403)
	ErrForbiddenAdminOnly    ErrorCode = "FORBIDDEN_ADMIN_ONLY"
	ErrForbiddenSelfDelete   ErrorCode = "FORBIDDEN_SELF_DELETE"
	ErrForbiddenSelfRole     ErrorCode = "FORBIDDEN_SELF_ROLE_CHANGE"

	// Validation Errors (400)
	ErrValidationRequired   ErrorCode = "VALIDATION_REQUIRED_FIELD"
	ErrValidationFormat     ErrorCode = "VALIDATION_INVALID_FORMAT"
	ErrValidationMinLength  ErrorCode = "VALIDATION_MIN_LENGTH"
	ErrValidationMaxLength  ErrorCode = "VALIDATION_MAX_LENGTH"
	ErrValidationDateFormat ErrorCode = "VALIDATION_INVALID_DATE"
	ErrValidationDateRange  ErrorCode = "VALIDATION_DATE_RANGE"
	ErrValidationPastDate   ErrorCode = "VALIDATION_PAST_DATE"

	// Resource Errors (404/409)
	ErrUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrVacationNotFound ErrorCode = "VACATION_REQUEST_NOT_FOUND"
	ErrUsernameTaken    ErrorCode = "USERNAME_TAKEN"

	// Business Logic Errors (422)
	ErrInsufficientDays    ErrorCode = "INSUFFICIENT_VACATION_DAYS"
	ErrAlreadyReviewed     ErrorCode = "VACATION_ALREADY_REVIEWED"
	ErrPasswordMismatch    ErrorCode = "PASSWORD_MISMATCH"
	ErrOverlappingVacation ErrorCode = "OVERLAPPING_VACATION"

	// Server Errors (500)
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrDatabase     ErrorCode = "DATABASE_ERROR"
	ErrEmailService ErrorCode = "EMAIL_SERVICE_ERROR"
)

// ErrorResponse represents the standard error format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// AppError is the application-level error type
type AppError struct {
	Code       ErrorCode
	Message    string
	Details    map[string]interface{}
	StatusCode int
	Err        error // underlying error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ToResponse converts AppError to ErrorResponse
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		},
	}
}

// Error constructors

func NewAuthError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewForbiddenError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewValidationError(code ErrorCode, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(code ErrorCode, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusConflict,
	}
}

func NewBusinessError(code ErrorCode, message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       ErrInternal,
		Message:    "An unexpected error occurred",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
```

---

## 7. Services (Business Logic)

### 7.1 Authentication Service

`internal/service/auth.go`:

```go
package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourorg/vacaytracker-api/internal/config"
	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type Claims struct {
	UserID   string      `json:"userId"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", dto.NewInternalError(err)
	}
	if user == nil {
		return nil, "", dto.NewAuthError(dto.ErrAuthCredentials, "Invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", dto.NewAuthError(dto.ErrAuthCredentials, "Invalid username or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", dto.NewInternalError(err)
	}

	return user, token, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, dto.NewAuthError(dto.ErrAuthTokenInvalid, "Invalid authentication token")
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, dto.NewAuthError(dto.ErrAuthTokenExpired, "Authentication token has expired")
		}
		return nil, dto.NewAuthError(dto.ErrAuthTokenInvalid, "Invalid authentication token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, dto.NewAuthError(dto.ErrAuthTokenInvalid, "Invalid authentication token")
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compares a password with a hash
func (s *AuthService) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) generateToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
```

### 7.2 User Service

`internal/service/user.go`:

```go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository"
)

type UserService struct {
	userRepo    repository.UserRepository
	authService *AuthService
}

func NewUserService(userRepo repository.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	// Check if username exists
	exists, err := s.userRepo.UsernameExists(ctx, req.Username, "")
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if exists {
		return nil, dto.NewConflictError(dto.ErrUsernameTaken, "Username already exists", map[string]interface{}{
			"username": req.Username,
		})
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}

	// Create user
	now := time.Now()
	role := domain.Role(req.Role)
	
	var emailPrefs domain.EmailPreferences
	if role == domain.RoleAdmin {
		emailPrefs = domain.DefaultAdminEmailPreferences()
	} else {
		emailPrefs = domain.DefaultEmployeeEmailPreferences()
	}

	user := &domain.User{
		ID:               "usr_" + uuid.New().String()[:8],
		Name:             req.Name,
		Username:         req.Username,
		Email:            req.Email,
		Password:         hashedPassword,
		Role:             role,
		VacationDays:     req.VacationDays,
		UsedVacationDays: 0,
		EmailPreferences: emailPrefs,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, dto.NewInternalError(err)
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if user == nil {
		return nil, dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}
	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	return users, nil
}

func (s *UserService) GetEmployees(ctx context.Context) ([]domain.User, error) {
	users, err := s.userRepo.GetByRole(ctx, domain.RoleEmployee)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	return users, nil
}

func (s *UserService) Update(ctx context.Context, id string, req dto.UpdateUserRequest, currentUserID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if user == nil {
		return nil, dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}

	// Check username uniqueness if changing
	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.UsernameExists(ctx, req.Username, id)
		if err != nil {
			return nil, dto.NewInternalError(err)
		}
		if exists {
			return nil, dto.NewConflictError(dto.ErrUsernameTaken, "Username already exists", map[string]interface{}{
				"username": req.Username,
			})
		}
		user.Username = req.Username
	}

	// Prevent self role change
	if req.Role != "" && req.Role != string(user.Role) && id == currentUserID {
		return nil, dto.NewForbiddenError(dto.ErrForbiddenSelfRole, "Cannot change your own role")
	}

	// Apply updates
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = domain.Role(req.Role)
	}
	if req.VacationDays != nil {
		user.VacationDays = *req.VacationDays
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, dto.NewInternalError(err)
	}

	return user, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, id string, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return dto.NewInternalError(err)
	}
	if user == nil {
		return dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}

	// Verify current password
	if !s.authService.ComparePassword(currentPassword, user.Password) {
		return dto.NewBusinessError(dto.ErrPasswordMismatch, "Current password is incorrect", nil)
	}

	// Hash new password
	hashedPassword, err := s.authService.HashPassword(newPassword)
	if err != nil {
		return dto.NewInternalError(err)
	}

	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}

func (s *UserService) UpdateEmailPreferences(ctx context.Context, id string, req dto.UpdateEmailPreferencesRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if user == nil {
		return nil, dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}

	// Apply preference updates
	if req.Enabled != nil {
		user.EmailPreferences.Enabled = *req.Enabled
	}
	if req.VacationStatusUpdates != nil {
		user.EmailPreferences.VacationStatusUpdates = *req.VacationStatusUpdates
	}
	if req.VacationRequestNotifications != nil {
		user.EmailPreferences.VacationRequestNotifications = *req.VacationRequestNotifications
	}
	if req.UserCreatedNotifications != nil {
		user.EmailPreferences.UserCreatedNotifications = *req.UserCreatedNotifications
	}
	if req.VacationResetNotifications != nil {
		user.EmailPreferences.VacationResetNotifications = *req.VacationResetNotifications
	}
	if req.MonthlyVacationSummary != nil {
		user.EmailPreferences.MonthlyVacationSummary = *req.MonthlyVacationSummary
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, dto.NewInternalError(err)
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id string, currentUserID string) error {
	if id == currentUserID {
		return dto.NewForbiddenError(dto.ErrForbiddenSelfDelete, "Cannot delete your own account")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return dto.NewInternalError(err)
	}
	if user == nil {
		return dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ResetVacationDays(ctx context.Context, days int) error {
	return s.userRepo.ResetAllVacationDays(ctx, days)
}
```

### 7.3 Vacation Service

`internal/service/vacation.go`:

```go
package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository"
)

type VacationService struct {
	vacationRepo repository.VacationRepository
	userRepo     repository.UserRepository
	settingsRepo repository.SettingsRepository
}

func NewVacationService(
	vacationRepo repository.VacationRepository,
	userRepo repository.UserRepository,
	settingsRepo repository.SettingsRepository,
) *VacationService {
	return &VacationService{
		vacationRepo: vacationRepo,
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
	}
}

func (s *VacationService) Create(ctx context.Context, userID string, req dto.CreateVacationRequest) (*domain.VacationRequest, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, dto.NewValidationError(dto.ErrValidationDateFormat, "Invalid date format", map[string]interface{}{
			"field":    "startDate",
			"expected": "YYYY-MM-DD",
		})
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, dto.NewValidationError(dto.ErrValidationDateFormat, "Invalid date format", map[string]interface{}{
			"field":    "endDate",
			"expected": "YYYY-MM-DD",
		})
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, dto.NewValidationError(dto.ErrValidationDateRange, "End date must be after start date", map[string]interface{}{
			"startDate": req.StartDate,
			"endDate":   req.EndDate,
		})
	}

	// Validate not in past
	today := time.Now().Truncate(24 * time.Hour)
	if startDate.Before(today) {
		return nil, dto.NewValidationError(dto.ErrValidationPastDate, "Start date cannot be in the past", map[string]interface{}{
			"field": "startDate",
		})
	}

	// Check for overlapping requests
	hasOverlap, err := s.vacationRepo.HasOverlap(ctx, userID, req.StartDate, req.EndDate, "")
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if hasOverlap {
		return nil, dto.NewBusinessError(dto.ErrOverlappingVacation, "Vacation dates overlap with existing request", nil)
	}

	// Calculate business days
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	businessDays := s.calculateBusinessDays(startDate, endDate, settings.WeekendPolicy.ExcludeWeekends)

	// Check sufficient days
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if user == nil {
		return nil, dto.NewNotFoundError(dto.ErrUserNotFound, "User not found")
	}

	available := user.RemainingDays()
	if businessDays > available {
		return nil, dto.NewBusinessError(dto.ErrInsufficientDays, "Not enough vacation days remaining", map[string]interface{}{
			"requested": businessDays,
			"available": available,
			"shortfall": businessDays - available,
		})
	}

	// Create request
	vacationReq := &domain.VacationRequest{
		ID:           "vac_" + uuid.New().String()[:8],
		UserID:       userID,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		BusinessDays: businessDays,
		Status:       domain.StatusPending,
		Reason:       req.Reason,
		CreatedAt:    time.Now(),
	}

	if err := s.vacationRepo.Create(ctx, vacationReq); err != nil {
		return nil, dto.NewInternalError(err)
	}

	return vacationReq, nil
}

func (s *VacationService) GetByID(ctx context.Context, id string) (*domain.VacationRequest, error) {
	req, err := s.vacationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if req == nil {
		return nil, dto.NewNotFoundError(dto.ErrVacationNotFound, "Vacation request not found")
	}
	return req, nil
}

func (s *VacationService) GetByUserID(ctx context.Context, userID string) ([]domain.VacationRequest, error) {
	requests, err := s.vacationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	return requests, nil
}

func (s *VacationService) GetAll(ctx context.Context) ([]domain.VacationRequest, error) {
	requests, err := s.vacationRepo.GetAll(ctx)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	return requests, nil
}

func (s *VacationService) GetPending(ctx context.Context) ([]domain.VacationRequest, error) {
	requests, err := s.vacationRepo.GetPending(ctx)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	return requests, nil
}

func (s *VacationService) Review(ctx context.Context, id string, status domain.VacationStatus, reviewerID string) (*domain.VacationRequest, error) {
	req, err := s.vacationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, dto.NewInternalError(err)
	}
	if req == nil {
		return nil, dto.NewNotFoundError(dto.ErrVacationNotFound, "Vacation request not found")
	}

	if !req.CanBeReviewed() {
		return nil, dto.NewBusinessError(dto.ErrAlreadyReviewed, "Request has already been processed", map[string]interface{}{
			"currentStatus": req.Status,
		})
	}

	// Update request
	now := time.Now()
	req.Status = status
	req.ReviewedBy = &reviewerID
	req.ReviewedAt = &now

	if err := s.vacationRepo.Update(ctx, req); err != nil {
		return nil, dto.NewInternalError(err)
	}

	// If approved, update user's used vacation days
	if status == domain.StatusApproved {
		user, err := s.userRepo.GetByID(ctx, req.UserID)
		if err != nil {
			return nil, dto.NewInternalError(err)
		}
		if user != nil {
			newUsed := user.UsedVacationDays + req.BusinessDays
			if err := s.userRepo.UpdateVacationDays(ctx, user.ID, newUsed); err != nil {
				return nil, dto.NewInternalError(err)
			}
		}
	}

	return req, nil
}

func (s *VacationService) Cancel(ctx context.Context, id string, userID string) error {
	req, err := s.vacationRepo.GetByID(ctx, id)
	if err != nil {
		return dto.NewInternalError(err)
	}
	if req == nil {
		return dto.NewNotFoundError(dto.ErrVacationNotFound, "Vacation request not found")
	}

	// Verify ownership
	if req.UserID != userID {
		return dto.NewForbiddenError(dto.ErrForbiddenAdminOnly, "Cannot cancel another user's request")
	}

	// Only pending requests can be cancelled
	if !req.IsPending() {
		return dto.NewBusinessError(dto.ErrAlreadyReviewed, "Only pending requests can be cancelled", map[string]interface{}{
			"currentStatus": req.Status,
		})
	}

	return s.vacationRepo.Delete(ctx, id)
}

// calculateBusinessDays calculates working days between two dates
func (s *VacationService) calculateBusinessDays(start, end time.Time, excludeWeekends bool) int {
	if !excludeWeekends {
		return int(end.Sub(start).Hours()/24) + 1
	}

	days := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			days++
		}
	}
	return days
}
```

### 7.4 Email Service

`internal/service/email.go`:

```go
package service

import (
	"fmt"

	"github.com/resend/resend-go/v2"

	"github.com/yourorg/vacaytracker-api/internal/config"
	"github.com/yourorg/vacaytracker-api/internal/domain"
)

type EmailService struct {
	client      *resend.Client
	cfg         *config.Config
	fromAddress string
	fromName    string
}

func NewEmailService(cfg *config.Config) *EmailService {
	var client *resend.Client
	if cfg.IsEmailEnabled() {
		client = resend.NewClient(cfg.ResendAPIKey)
	}

	return &EmailService{
		client:      client,
		cfg:         cfg,
		fromAddress: cfg.EmailFromAddress,
		fromName:    cfg.EmailFromName,
	}
}

func (s *EmailService) IsEnabled() bool {
	return s.client != nil
}

func (s *EmailService) from() string {
	return fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress)
}

// SendVacationApproved notifies employee their request was approved
func (s *EmailService) SendVacationApproved(user *domain.User, request *domain.VacationRequest) error {
	if !s.IsEnabled() || !user.EmailPreferences.Enabled || !user.EmailPreferences.VacationStatusUpdates || user.Email == "" {
		return nil
	}

	html := fmt.Sprintf(`
		<h2>🏖️ Vacation Request Approved!</h2>
		<p>Hi %s,</p>
		<p>Great news! Your vacation request has been approved.</p>
		<ul>
			<li><strong>Dates:</strong> %s to %s</li>
			<li><strong>Days:</strong> %d business days</li>
		</ul>
		<p>Enjoy your time off! 🌴</p>
		<p>— VacayTracker</p>
	`, user.Name, request.StartDate, request.EndDate, request.BusinessDays)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from(),
		To:      []string{user.Email},
		Subject: "✅ Vacation Request Approved",
		Html:    html,
	})

	return err
}

// SendVacationRejected notifies employee their request was rejected
func (s *EmailService) SendVacationRejected(user *domain.User, request *domain.VacationRequest) error {
	if !s.IsEnabled() || !user.EmailPreferences.Enabled || !user.EmailPreferences.VacationStatusUpdates || user.Email == "" {
		return nil
	}

	html := fmt.Sprintf(`
		<h2>Vacation Request Update</h2>
		<p>Hi %s,</p>
		<p>Unfortunately, your vacation request was not approved.</p>
		<ul>
			<li><strong>Dates:</strong> %s to %s</li>
			<li><strong>Days:</strong> %d business days</li>
		</ul>
		<p>Please contact your administrator for more information.</p>
		<p>— VacayTracker</p>
	`, user.Name, request.StartDate, request.EndDate, request.BusinessDays)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from(),
		To:      []string{user.Email},
		Subject: "Vacation Request Update",
		Html:    html,
	})

	return err
}

// SendNewVacationRequest notifies admins of a new pending request
func (s *EmailService) SendNewVacationRequest(admin *domain.User, employee *domain.User, request *domain.VacationRequest) error {
	if !s.IsEnabled() || !admin.EmailPreferences.Enabled || !admin.EmailPreferences.VacationRequestNotifications || admin.Email == "" {
		return nil
	}

	html := fmt.Sprintf(`
		<h2>📋 New Vacation Request</h2>
		<p>Hi %s,</p>
		<p><strong>%s</strong> has submitted a new vacation request.</p>
		<ul>
			<li><strong>Dates:</strong> %s to %s</li>
			<li><strong>Days:</strong> %d business days</li>
			<li><strong>Reason:</strong> %s</li>
		</ul>
		<p><a href="%s/admin">Review in VacayTracker</a></p>
		<p>— VacayTracker</p>
	`, admin.Name, employee.Name, request.StartDate, request.EndDate, request.BusinessDays, request.Reason, s.cfg.AppURL)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from(),
		To:      []string{admin.Email},
		Subject: fmt.Sprintf("New Vacation Request from %s", employee.Name),
		Html:    html,
	})

	return err
}

// SendUserCreated confirms new user creation to admin
func (s *EmailService) SendUserCreated(admin *domain.User, newUser *domain.User) error {
	if !s.IsEnabled() || !admin.EmailPreferences.Enabled || !admin.EmailPreferences.UserCreatedNotifications || admin.Email == "" {
		return nil
	}

	html := fmt.Sprintf(`
		<h2>👤 New User Created</h2>
		<p>Hi %s,</p>
		<p>A new user account has been created:</p>
		<ul>
			<li><strong>Name:</strong> %s</li>
			<li><strong>Username:</strong> %s</li>
			<li><strong>Role:</strong> %s</li>
		</ul>
		<p>— VacayTracker</p>
	`, admin.Name, newUser.Name, newUser.Username, newUser.Role)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from(),
		To:      []string{admin.Email},
		Subject: fmt.Sprintf("New User Created: %s", newUser.Name),
		Html:    html,
	})

	return err
}

// SendVacationReset notifies users about year-end reset
func (s *EmailService) SendVacationReset(user *domain.User, newDays int) error {
	if !s.IsEnabled() || !user.EmailPreferences.Enabled || !user.EmailPreferences.VacationResetNotifications || user.Email == "" {
		return nil
	}

	html := fmt.Sprintf(`
		<h2>🔄 Vacation Days Reset</h2>
		<p>Hi %s,</p>
		<p>Your vacation balance has been reset for the new year.</p>
		<ul>
			<li><strong>New Balance:</strong> %d days</li>
		</ul>
		<p>Happy new year! 🎉</p>
		<p>— VacayTracker</p>
	`, user.Name, newDays)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from(),
		To:      []string{user.Email},
		Subject: "Vacation Days Reset for New Year",
		Html:    html,
	})

	return err
}
```

---

## 8. HTTP Handlers

### 8.1 Handler Dependencies

`internal/handler/handler.go`:

```go
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/yourorg/vacaytracker-api/internal/config"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

type Handler struct {
	cfg             *config.Config
	authService     *service.AuthService
	userService     *service.UserService
	vacationService *service.VacationService
	emailService    *service.EmailService
	validate        *validator.Validate
}

func New(
	cfg *config.Config,
	authService *service.AuthService,
	userService *service.UserService,
	vacationService *service.VacationService,
	emailService *service.EmailService,
) *Handler {
	return &Handler{
		cfg:             cfg,
		authService:     authService,
		userService:     userService,
		vacationService: vacationService,
		emailService:    emailService,
		validate:        validator.New(),
	}
}

// Helper to send success responses
func (h *Handler) sendData(c *gin.Context, status int, data interface{}) {
	c.JSON(status, dto.DataResponse{Data: data})
}

func (h *Handler) sendDataWithMessage(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, dto.DataWithMessageResponse{Data: data, Message: message})
}

func (h *Handler) sendList(c *gin.Context, data interface{}, total, count int) {
	c.JSON(200, dto.ListResponse{
		Data: data,
		Meta: dto.MetaInfo{Total: total, Count: count},
	})
}

// Helper to send error responses
func (h *Handler) sendError(c *gin.Context, err error) {
	if appErr, ok := err.(*dto.AppError); ok {
		c.JSON(appErr.StatusCode, appErr.ToResponse())
		return
	}
	// Unknown error - wrap as internal
	internalErr := dto.NewInternalError(err)
	c.JSON(internalErr.StatusCode, internalErr.ToResponse())
}
```

### 8.2 Auth Handlers

`internal/handler/auth.go`:

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
)

// POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  h.toUserResponse(user),
	})
}

// GET /api/auth/me
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")
	
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	h.sendData(c, http.StatusOK, h.toUserResponse(user))
}

// Helper to convert domain user to response
func (h *Handler) toUserResponse(user *domain.User) dto.UserResponse {
	resp := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		EmailPreferences: dto.EmailPreferencesResponse{
			Enabled:                      user.EmailPreferences.Enabled,
			VacationStatusUpdates:        user.EmailPreferences.VacationStatusUpdates,
			VacationRequestNotifications: user.EmailPreferences.VacationRequestNotifications,
			UserCreatedNotifications:     user.EmailPreferences.UserCreatedNotifications,
			VacationResetNotifications:   user.EmailPreferences.VacationResetNotifications,
			MonthlyVacationSummary:       user.EmailPreferences.MonthlyVacationSummary,
		},
	}

	if user.Role == domain.RoleEmployee {
		resp.VacationDays = user.VacationDays
		resp.UsedVacationDays = user.UsedVacationDays
		resp.RemainingDays = user.RemainingDays()
	}

	return resp
}
```

### 8.3 User Handlers

`internal/handler/user.go`:

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/dto"
)

// PUT /api/users/:id/password
func (h *Handler) UpdatePassword(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	// Users can only change their own password
	if userID != currentUserID {
		h.sendError(c, dto.NewForbiddenError(dto.ErrForbiddenAdminOnly, "Cannot change another user's password"))
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	err := h.userService.UpdatePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// PUT /api/users/:id/email-preferences
func (h *Handler) UpdateEmailPreferences(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	// Users can only change their own preferences
	if userID != currentUserID {
		h.sendError(c, dto.NewForbiddenError(dto.ErrForbiddenAdminOnly, "Cannot change another user's preferences"))
		return
	}

	var req dto.UpdateEmailPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	user, err := h.userService.UpdateEmailPreferences(c.Request.Context(), userID, req)
	if err != nil {
		h.sendError(c, err)
		return
	}

	h.sendData(c, http.StatusOK, h.toUserResponse(user))
}
```

### 8.4 Vacation Handlers

`internal/handler/vacation.go`:

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
)

// POST /api/vacation
func (h *Handler) CreateVacationRequest(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.CreateVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	vacationReq, err := h.vacationService.Create(c.Request.Context(), userID, req)
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Send notification to admins (async, don't fail request if email fails)
	go h.notifyAdminsOfNewRequest(vacationReq)

	h.sendDataWithMessage(c, http.StatusCreated, h.toVacationResponse(vacationReq, nil), "Vacation request submitted successfully")
}

// GET /api/vacation
func (h *Handler) GetMyVacationRequests(c *gin.Context) {
	userID := c.GetString("userID")

	requests, err := h.vacationService.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	responses := make([]dto.VacationRequestResponse, len(requests))
	for i, req := range requests {
		responses[i] = h.toVacationResponse(&req, nil)
	}

	h.sendList(c, responses, len(responses), len(responses))
}

// DELETE /api/vacation/:id
func (h *Handler) CancelVacationRequest(c *gin.Context) {
	requestID := c.Param("id")
	userID := c.GetString("userID")

	err := h.vacationService.Cancel(c.Request.Context(), requestID, userID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper to convert domain vacation request to response
func (h *Handler) toVacationResponse(req *domain.VacationRequest, userName *string) dto.VacationRequestResponse {
	resp := dto.VacationRequestResponse{
		ID:           req.ID,
		UserID:       req.UserID,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		BusinessDays: req.BusinessDays,
		Status:       string(req.Status),
		Reason:       req.Reason,
		CreatedAt:    req.CreatedAt,
		ReviewedBy:   req.ReviewedBy,
		ReviewedAt:   req.ReviewedAt,
	}

	if userName != nil {
		resp.UserName = *userName
	}

	return resp
}

func (h *Handler) notifyAdminsOfNewRequest(vacationReq *domain.VacationRequest) {
	ctx := c.Request.Context()
	
	// Get employee
	employee, _ := h.userService.GetByID(ctx, vacationReq.UserID)
	if employee == nil {
		return
	}

	// Get all admins
	users, _ := h.userService.GetAll(ctx)
	for _, user := range users {
		if user.IsAdmin() {
			_ = h.emailService.SendNewVacationRequest(&user, employee, vacationReq)
		}
	}
}
```

### 8.5 Admin Handlers

`internal/handler/admin.go`:

```go
package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
)

// GET /api/admin/users
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAll(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = h.toUserResponse(&user)
	}

	h.sendList(c, responses, len(responses), len(responses))
}

// POST /api/admin/users
func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Notify admin of creation
	currentUserID := c.GetString("userID")
	go h.notifyAdminUserCreated(c.Request.Context(), currentUserID, user)

	h.sendDataWithMessage(c, http.StatusCreated, h.toUserResponse(user), "User created successfully")
}

// PUT /api/admin/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, req, currentUserID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	h.sendData(c, http.StatusOK, h.toUserResponse(user))
}

// DELETE /api/admin/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	err := h.userService.Delete(c.Request.Context(), userID, currentUserID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /api/admin/vacation
func (h *Handler) GetAllVacationRequests(c *gin.Context) {
	requests, err := h.vacationService.GetAll(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Build user map for names
	users, _ := h.userService.GetAll(c.Request.Context())
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	responses := make([]dto.VacationRequestResponse, len(requests))
	for i, req := range requests {
		name := userMap[req.UserID]
		responses[i] = h.toVacationResponse(&req, &name)
	}

	h.sendList(c, responses, len(responses), len(responses))
}

// GET /api/admin/vacation/pending
func (h *Handler) GetPendingVacationRequests(c *gin.Context) {
	requests, err := h.vacationService.GetPending(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Build user map for names
	users, _ := h.userService.GetAll(c.Request.Context())
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[user.ID] = user.Name
	}

	responses := make([]dto.VacationRequestResponse, len(requests))
	for i, req := range requests {
		name := userMap[req.UserID]
		responses[i] = h.toVacationResponse(&req, &name)
	}

	h.sendList(c, responses, len(responses), len(responses))
}

// PUT /api/admin/vacation/:id/review
func (h *Handler) ReviewVacationRequest(c *gin.Context) {
	requestID := c.Param("id")
	currentUserID := c.GetString("userID")

	var req dto.ReviewVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	status := domain.VacationStatus(req.Status)
	vacationReq, err := h.vacationService.Review(c.Request.Context(), requestID, status, currentUserID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Send notification to employee
	go h.notifyEmployeeOfReview(c.Request.Context(), vacationReq)

	h.sendData(c, http.StatusOK, h.toVacationResponse(vacationReq, nil))
}

// POST /api/admin/reset-vacation-days
func (h *Handler) ResetVacationDays(c *gin.Context) {
	var req dto.ResetVacationDaysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, dto.NewValidationError(dto.ErrValidationRequired, "Invalid request body", nil))
		return
	}

	err := h.userService.ResetVacationDays(c.Request.Context(), req.Days)
	if err != nil {
		h.sendError(c, err)
		return
	}

	// Notify all users
	go h.notifyUsersOfReset(c.Request.Context(), req.Days)

	c.JSON(http.StatusOK, gin.H{"message": "Vacation days reset successfully"})
}

// Helper functions
func (h *Handler) notifyAdminUserCreated(ctx context.Context, adminID string, newUser *domain.User) {
	admin, _ := h.userService.GetByID(ctx, adminID)
	if admin != nil {
		_ = h.emailService.SendUserCreated(admin, newUser)
	}
}

func (h *Handler) notifyEmployeeOfReview(ctx context.Context, vacationReq *domain.VacationRequest) {
	user, _ := h.userService.GetByID(ctx, vacationReq.UserID)
	if user == nil {
		return
	}

	if vacationReq.Status == domain.StatusApproved {
		_ = h.emailService.SendVacationApproved(user, vacationReq)
	} else if vacationReq.Status == domain.StatusRejected {
		_ = h.emailService.SendVacationRejected(user, vacationReq)
	}
}

func (h *Handler) notifyUsersOfReset(ctx context.Context, newDays int) {
	users, _ := h.userService.GetAll(ctx)
	for _, user := range users {
		_ = h.emailService.SendVacationReset(&user, newDays)
	}
}
```

### 8.6 Health Handler

`internal/handler/health.go`:

```go
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/dto"
)

const Version = "1.0.0"

// GET /health
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   Version,
	})
}
```

---

## 9. Middleware

### 9.1 JWT Authentication Middleware

`internal/middleware/auth.go`:

```go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendAuthError(c, dto.ErrAuthTokenMissing, "Authentication token required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendAuthError(c, dto.ErrAuthTokenInvalid, "Invalid authorization header format")
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			if appErr, ok := err.(*dto.AppError); ok {
				c.JSON(appErr.StatusCode, appErr.ToResponse())
			} else {
				sendAuthError(c, dto.ErrAuthTokenInvalid, "Invalid authentication token")
			}
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(domain.Role) != domain.RoleAdmin {
			c.JSON(403, dto.ErrorResponse{
				Error: dto.ErrorDetail{
					Code:    dto.ErrForbiddenAdminOnly,
					Message: "Admin access required",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func EmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(domain.Role) != domain.RoleEmployee {
			c.JSON(403, dto.ErrorResponse{
				Error: dto.ErrorDetail{
					Code:    dto.ErrAuthRoleMismatch,
					Message: "Employee access required",
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func sendAuthError(c *gin.Context, code dto.ErrorCode, message string) {
	c.JSON(401, dto.ErrorResponse{
		Error: dto.ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
	c.Abort()
}
```

### 9.2 CORS Middleware

`internal/middleware/cors.go`:

```go
package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/config"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, o := range cfg.CORSOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

### 9.3 Error Recovery Middleware

`internal/middleware/error.go`:

```go
package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/dto"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error: dto.ErrorDetail{
						Code:    dto.ErrInternal,
						Message: "An unexpected error occurred",
					},
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
```

---

## 10. Application Entry Point

### 10.1 Main Function

`cmd/server/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/config"
	"github.com/yourorg/vacaytracker-api/internal/handler"
	"github.com/yourorg/vacaytracker-api/internal/middleware"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate("migrations/001_init.sql"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := sqlite.NewUserRepository(db)
	vacationRepo := sqlite.NewVacationRepository(db)
	settingsRepo := sqlite.NewSettingsRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo, authService)
	vacationService := service.NewVacationService(vacationRepo, userRepo, settingsRepo)
	emailService := service.NewEmailService(cfg)

	// Create default admin if not exists
	if err := createDefaultAdmin(userService, cfg); err != nil {
		log.Printf("Warning: Could not create default admin: %v", err)
	}

	// Initialize handlers
	h := handler.New(cfg, authService, userService, vacationService, emailService)

	// Setup router
	router := setupRouter(cfg, h, authService)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(cfg *config.Config, h *handler.Handler, authService *service.AuthService) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check (no auth)
	router.GET("/health", h.Health)

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Auth
			protected.GET("/auth/me", h.GetCurrentUser)

			// User self-service
			protected.PUT("/users/:id/password", h.UpdatePassword)
			protected.PUT("/users/:id/email-preferences", h.UpdateEmailPreferences)

			// Employee vacation routes
			employee := protected.Group("/vacation")
			employee.Use(middleware.EmployeeMiddleware())
			{
				employee.GET("", h.GetMyVacationRequests)
				employee.POST("", h.CreateVacationRequest)
				employee.DELETE("/:id", h.CancelVacationRequest)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// User management
				admin.GET("/users", h.GetAllUsers)
				admin.POST("/users", h.CreateUser)
				admin.PUT("/users/:id", h.UpdateUser)
				admin.DELETE("/users/:id", h.DeleteUser)

				// Vacation management
				admin.GET("/vacation", h.GetAllVacationRequests)
				admin.GET("/vacation/pending", h.GetPendingVacationRequests)
				admin.PUT("/vacation/:id/review", h.ReviewVacationRequest)

				// Settings
				admin.POST("/reset-vacation-days", h.ResetVacationDays)
			}
		}
	}

	return router
}

func createDefaultAdmin(userService *service.UserService, cfg *config.Config) error {
	ctx := context.Background()

	// Check if admin exists
	_, err := userService.GetByID(ctx, "usr_admin001")
	if err == nil {
		return nil // Admin already exists
	}

	// Create default admin
	_, err = userService.Create(ctx, dto.CreateUserRequest{
		Name:     "Administrator",
		Username: "admin",
		Password: cfg.AdminPassword,
		Role:     "admin",
	})

	if err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	log.Println("Default admin user created")
	return nil
}
```

---

## 11. API Endpoints Summary

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| GET | `/health` | No | - | Health check |
| POST | `/api/auth/login` | No | - | User login |
| GET | `/api/auth/me` | Yes | Any | Get current user |
| PUT | `/api/users/:id/password` | Yes | Self | Change password |
| PUT | `/api/users/:id/email-preferences` | Yes | Self | Update email preferences |
| GET | `/api/vacation` | Yes | Employee | Get my vacation requests |
| POST | `/api/vacation` | Yes | Employee | Submit vacation request |
| DELETE | `/api/vacation/:id` | Yes | Employee | Cancel pending request |
| GET | `/api/admin/users` | Yes | Admin | List all users |
| POST | `/api/admin/users` | Yes | Admin | Create user |
| PUT | `/api/admin/users/:id` | Yes | Admin | Update user |
| DELETE | `/api/admin/users/:id` | Yes | Admin | Delete user |
| GET | `/api/admin/vacation` | Yes | Admin | List all vacation requests |
| GET | `/api/admin/vacation/pending` | Yes | Admin | List pending requests |
| PUT | `/api/admin/vacation/:id/review` | Yes | Admin | Approve/reject request |
| POST | `/api/admin/reset-vacation-days` | Yes | Admin | Reset all vacation days |

---

## 12. Testing

### 12.1 Testing Structure

```
vacaytracker-api/
├── internal/
│   ├── service/
│   │   ├── auth_test.go
│   │   ├── user_test.go
│   │   └── vacation_test.go
│   └── handler/
│       └── handler_test.go
└── test/
    └── integration/
        └── api_test.go
```

### 12.2 Service Unit Test Example

`internal/service/vacation_test.go`:

```go
package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

// Mock repositories would be defined here...

func TestVacationService_CalculateBusinessDays(t *testing.T) {
	tests := []struct {
		name            string
		start           string
		end             string
		excludeWeekends bool
		expected        int
	}{
		{
			name:            "full week with weekends excluded",
			start:           "2024-01-08", // Monday
			end:             "2024-01-12", // Friday
			excludeWeekends: true,
			expected:        5,
		},
		{
			name:            "week spanning weekend with weekends excluded",
			start:           "2024-01-08", // Monday
			end:             "2024-01-15", // Monday
			excludeWeekends: true,
			expected:        6,
		},
		{
			name:            "week with weekends included",
			start:           "2024-01-08",
			end:             "2024-01-14",
			excludeWeekends: false,
			expected:        7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, _ := time.Parse("2006-01-02", tt.start)
			end, _ := time.Parse("2006-01-02", tt.end)

			// Test the calculation logic
			days := calculateBusinessDays(start, end, tt.excludeWeekends)
			if days != tt.expected {
				t.Errorf("expected %d days, got %d", tt.expected, days)
			}
		})
	}
}

func calculateBusinessDays(start, end time.Time, excludeWeekends bool) int {
	if !excludeWeekends {
		return int(end.Sub(start).Hours()/24) + 1
	}

	days := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		weekday := d.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			days++
		}
	}
	return days
}
```

### 12.3 Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## 13. Deployment

### 13.1 Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /vacaytracker cmd/server/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /vacaytracker .
COPY --from=builder /app/migrations ./migrations

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 3000

# Run
CMD ["./vacaytracker"]
```

### 13.2 Docker Compose

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - ENV=production
      - PORT=3000
      - DB_PATH=/app/data/vacaytracker.db
      - JWT_SECRET=${JWT_SECRET}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - RESEND_API_KEY=${RESEND_API_KEY}
      - EMAIL_FROM_ADDRESS=${EMAIL_FROM_ADDRESS}
      - APP_URL=${APP_URL}
    volumes:
      - vacaytracker-data:/app/data
    restart: unless-stopped

volumes:
  vacaytracker-data:
```

### 13.3 Build Commands

```bash
# Build for current platform
go build -o bin/vacaytracker cmd/server/main.go

# Build for Linux (cross-compile)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/vacaytracker-linux cmd/server/main.go

# Build Docker image
docker build -t vacaytracker-api .

# Run with Docker Compose
docker-compose up -d
```

---

## 14. Development Workflow

### 14.1 Local Development

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env
# Edit .env with your values

# Run migrations
make migrate

# Start development server (with hot reload using air)
go install github.com/air-verse/air@latest
air

# Or run directly
go run cmd/server/main.go
```

### 14.2 VS Code Setup

`.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  }
}
```

### 14.3 Air Configuration (Hot Reload)

`.air.toml`:

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/server/main.go"
  bin = "./tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor", "data"]

[log]
  time = false

[misc]
  clean_on_exit = true
```

---

## 15. Quick Reference

### Go Commands
```bash
go mod init <module>     # Initialize module
go mod tidy              # Clean up dependencies
go get <package>         # Add dependency
go build                 # Compile
go run main.go           # Run without building
go test ./...            # Run all tests
go vet ./...             # Static analysis
```

### API Testing with curl
```bash
# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret"}'

# Get current user
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer <token>"

# Create vacation request
curl -X POST http://localhost:3000/api/vacation \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"startDate":"2024-08-01","endDate":"2024-08-05","reason":"Summer break"}'
```

---

*This guide provides complete implementation patterns for VacayTracker's Go backend. Adapt and extend as needed for your specific requirements.*