# 04 - Backend Tasks

> Complete Go implementation checklist with code examples

## Table of Contents

1. [Project Scaffolding](#project-scaffolding)
2. [Configuration](#configuration)
3. [Database Layer](#database-layer)
4. [Domain Models](#domain-models)
5. [Repository Layer](#repository-layer)
6. [Service Layer](#service-layer)
7. [Middleware](#middleware)
8. [Handlers](#handlers)
9. [DTOs](#dtos)
10. [Testing](#testing)

---

## Project Scaffolding

### Task 1.1: Initialize Go Module

- [ ] **Create go.mod** `go.mod`

```bash
go mod init github.com/yourorg/vacaytracker-api
```

**go.mod content:**
```go
module github.com/yourorg/vacaytracker-api

go 1.23

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/google/uuid v1.5.0
    github.com/joho/godotenv v1.5.1
    golang.org/x/crypto v0.18.0
    modernc.org/sqlite v1.28.0
)
```

**Verification:**
```bash
go mod tidy
go mod verify
```

---

### Task 1.2: Create Directory Structure

- [ ] **Create project directories**

```
vacaytracker-api/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   │   └── sqlite/
│   └── service/
├── migrations/
├── data/
├── .env
├── .env.example
├── .gitignore
├── Makefile
└── README.md
```

**Commands:**
```bash
mkdir -p cmd/server
mkdir -p internal/{config,domain,dto,handler,middleware,repository/sqlite,service}
mkdir -p migrations data
touch cmd/server/main.go
touch .env .env.example .gitignore Makefile
```

---

### Task 1.3: Create Makefile

- [ ] **Create Makefile** `Makefile`

```makefile
.PHONY: build run test test-coverage lint clean migrate

# Build settings
BINARY_NAME=vacaytracker-api
BUILD_DIR=./bin

# Go settings
GOFLAGS=-ldflags="-w -s"

build:
	CGO_ENABLED=0 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server/main.go

dev:
	air -c .air.toml

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

migrate:
	go run ./cmd/server/main.go -migrate

# Docker targets
docker-build:
	docker build -t vacaytracker-api .

docker-run:
	docker run -p 3000:3000 --env-file .env vacaytracker-api
```

---

### Task 1.4: Create .gitignore

- [ ] **Create .gitignore** `.gitignore`

```gitignore
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test artifacts
*.test
coverage.out
coverage.html

# IDE
.idea/
.vscode/
*.swp
*.swo

# Environment
.env
.env.local

# Database
data/*.db
data/*.db-shm
data/*.db-wal

# OS
.DS_Store
Thumbs.db

# Build
tmp/
```

---

### Task 1.5: Create .env.example

- [ ] **Create .env.example** `.env.example`

```env
# Server Configuration
PORT=3000
ENV=development
APP_URL=http://localhost:3000

# Database
DB_PATH=./data/vacaytracker.db

# Authentication
JWT_SECRET=your-secure-secret-key-minimum-32-characters
ADMIN_PASSWORD=admin123

# Email (Resend)
RESEND_API_KEY=re_xxxxxxxxxxxxx
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
EMAIL_FROM_NAME=VacayTracker

# Admin Setup
ADMIN_EMAIL=admin@company.com
ADMIN_NAME=Captain Admin
```

---

## Configuration

### Task 2.1: Implement Config Loader

- [ ] **Create config.go** `internal/config/config.go`

```go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string
	URL  string

	// Database
	DBPath string

	// Auth
	JWTSecret     string
	AdminPassword string
	AdminEmail    string
	AdminName     string

	// Email
	ResendAPIKey     string
	EmailFromAddress string
	EmailFromName    string
}

func Load() *Config {
	// Load .env file (ignore error if not exists in production)
	_ = godotenv.Load()

	cfg := &Config{
		// Server defaults
		Port: getEnv("PORT", "3000"),
		Env:  getEnv("ENV", "development"),
		URL:  getEnv("APP_URL", "http://localhost:3000"),

		// Database
		DBPath: getEnv("DB_PATH", "./data/vacaytracker.db"),

		// Auth (required)
		JWTSecret:     mustGetEnv("JWT_SECRET"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@company.com"),
		AdminName:     getEnv("ADMIN_NAME", "Admin"),

		// Email (optional)
		ResendAPIKey:     getEnv("RESEND_API_KEY", ""),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", ""),
		EmailFromName:    getEnv("EMAIL_FROM_NAME", "VacayTracker"),
	}

	// Validate JWT secret length
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	return cfg
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) EmailEnabled() bool {
	return c.ResendAPIKey != "" && c.EmailFromAddress != ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
```

---

## Database Layer

### Task 3.1: Create SQLite Connection

- [ ] **Create sqlite.go** `internal/repository/sqlite/sqlite.go`

```go
package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	// Ensure data directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database with WAL mode and foreign keys
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer
	db.SetMaxIdleConns(1)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to SQLite database: %s", dbPath)
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
```

---

### Task 3.2: Create Migration Runner

- [ ] **Create migrate.go** `internal/repository/sqlite/migrate.go`

```go
package sqlite

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (db *DB) RunMigrations(migrationsDir string) error {
	// Create migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}

	// Read migration files
	var migrations []string
	err = filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			migrations = append(migrations, d.Name())
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort migrations by name (assumes numeric prefix like 001_, 002_)
	sort.Strings(migrations)

	// Apply pending migrations
	for _, migration := range migrations {
		if applied[migration] {
			continue
		}

		log.Printf("Applying migration: %s", migration)

		content, err := os.ReadFile(filepath.Join(migrationsDir, migration))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", migration, err)
		}

		// Execute migration in transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration, err)
		}

		log.Printf("Applied migration: %s", migration)
	}

	return nil
}
```

---

### Task 3.3: Create Initial Migration

- [ ] **Create 001_init.sql** `migrations/001_init.sql`

```sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'employee' CHECK (role IN ('admin', 'employee')),
    vacation_balance INTEGER NOT NULL DEFAULT 25,
    start_date TEXT,
    email_preferences TEXT NOT NULL DEFAULT '{"vacationUpdates":true,"weeklyDigest":false,"teamNotifications":true}',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Vacation requests table
CREATE TABLE IF NOT EXISTS vacation_requests (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    total_days INTEGER NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TEXT,
    rejection_reason TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Settings table (singleton)
CREATE TABLE IF NOT EXISTS settings (
    id TEXT PRIMARY KEY DEFAULT 'settings',
    weekend_policy TEXT NOT NULL DEFAULT '{"excludeWeekends":true,"excludedDays":[0,6]}',
    newsletter TEXT NOT NULL DEFAULT '{"enabled":false,"frequency":"monthly","dayOfMonth":1}',
    default_vacation_days INTEGER NOT NULL DEFAULT 25,
    vacation_reset_month INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_vacation_requests_user_id ON vacation_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_vacation_requests_status ON vacation_requests(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Triggers for updated_at
CREATE TRIGGER IF NOT EXISTS users_updated_at
    AFTER UPDATE ON users
    FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS vacation_requests_updated_at
    AFTER UPDATE ON vacation_requests
    FOR EACH ROW
BEGIN
    UPDATE vacation_requests SET updated_at = datetime('now') WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS settings_updated_at
    AFTER UPDATE ON settings
    FOR EACH ROW
BEGIN
    UPDATE settings SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Insert default settings
INSERT OR IGNORE INTO settings (id) VALUES ('settings');
```

---

## Domain Models

### Task 4.1: Create User Domain

- [ ] **Create user.go** `internal/domain/user.go`

```go
package domain

import (
	"encoding/json"
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type EmailPreferences struct {
	VacationUpdates   bool `json:"vacationUpdates"`
	WeeklyDigest      bool `json:"weeklyDigest"`
	TeamNotifications bool `json:"teamNotifications"`
}

func DefaultEmailPreferences() EmailPreferences {
	return EmailPreferences{
		VacationUpdates:   true,
		WeeklyDigest:      false,
		TeamNotifications: true,
	}
}

type User struct {
	ID               string           `json:"id"`
	Email            string           `json:"email"`
	PasswordHash     string           `json:"-"` // Never expose in JSON
	Name             string           `json:"name"`
	Role             Role             `json:"role"`
	VacationBalance  int              `json:"vacationBalance"`
	StartDate        *string          `json:"startDate,omitempty"`
	EmailPreferences EmailPreferences `json:"emailPreferences"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

// MarshalEmailPreferences converts EmailPreferences to JSON string for storage
func (e EmailPreferences) MarshalJSON() ([]byte, error) {
	type Alias EmailPreferences
	return json.Marshal(Alias(e))
}

// ToJSONString returns the EmailPreferences as a JSON string
func (e EmailPreferences) ToJSONString() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// ParseEmailPreferences parses JSON string to EmailPreferences
func ParseEmailPreferences(data string) EmailPreferences {
	var prefs EmailPreferences
	if err := json.Unmarshal([]byte(data), &prefs); err != nil {
		return DefaultEmailPreferences()
	}
	return prefs
}
```

---

### Task 4.2: Create Vacation Domain

- [ ] **Create vacation.go** `internal/domain/vacation.go`

```go
package domain

import "time"

type VacationStatus string

const (
	StatusPending  VacationStatus = "pending"
	StatusApproved VacationStatus = "approved"
	StatusRejected VacationStatus = "rejected"
)

type VacationRequest struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	UserName        string         `json:"userName,omitempty"` // Populated by joins
	UserEmail       string         `json:"userEmail,omitempty"` // Populated by joins
	StartDate       string         `json:"startDate"` // YYYY-MM-DD
	EndDate         string         `json:"endDate"`   // YYYY-MM-DD
	TotalDays       int            `json:"totalDays"`
	Reason          *string        `json:"reason,omitempty"`
	Status          VacationStatus `json:"status"`
	ReviewedBy      *string        `json:"reviewedBy,omitempty"`
	ReviewedAt      *time.Time     `json:"reviewedAt,omitempty"`
	RejectionReason *string        `json:"rejectionReason,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

func (v *VacationRequest) IsPending() bool {
	return v.Status == StatusPending
}

func (v *VacationRequest) IsApproved() bool {
	return v.Status == StatusApproved
}

func (v *VacationRequest) IsRejected() bool {
	return v.Status == StatusRejected
}

func (v *VacationRequest) CanBeCancelled() bool {
	return v.Status == StatusPending
}

// TeamVacation is a simplified view for team calendar
type TeamVacation struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	TotalDays int    `json:"totalDays"`
}
```

---

### Task 4.3: Create Settings Domain

- [ ] **Create settings.go** `internal/domain/settings.go`

```go
package domain

import (
	"encoding/json"
	"time"
)

type WeekendPolicy struct {
	ExcludeWeekends bool  `json:"excludeWeekends"`
	ExcludedDays    []int `json:"excludedDays"` // 0 = Sunday, 6 = Saturday
}

func DefaultWeekendPolicy() WeekendPolicy {
	return WeekendPolicy{
		ExcludeWeekends: true,
		ExcludedDays:    []int{0, 6}, // Sunday and Saturday
	}
}

type NewsletterConfig struct {
	Enabled    bool   `json:"enabled"`
	Frequency  string `json:"frequency"` // "weekly" or "monthly"
	DayOfMonth int    `json:"dayOfMonth"` // 1-28 for monthly
}

func DefaultNewsletterConfig() NewsletterConfig {
	return NewsletterConfig{
		Enabled:    false,
		Frequency:  "monthly",
		DayOfMonth: 1,
	}
}

type Settings struct {
	ID                  string           `json:"id"`
	WeekendPolicy       WeekendPolicy    `json:"weekendPolicy"`
	Newsletter          NewsletterConfig `json:"newsletter"`
	DefaultVacationDays int              `json:"defaultVacationDays"`
	VacationResetMonth  int              `json:"vacationResetMonth"` // 1-12
	UpdatedAt           time.Time        `json:"updatedAt"`
}

func DefaultSettings() *Settings {
	return &Settings{
		ID:                  "settings",
		WeekendPolicy:       DefaultWeekendPolicy(),
		Newsletter:          DefaultNewsletterConfig(),
		DefaultVacationDays: 25,
		VacationResetMonth:  1, // January
	}
}

// ToJSONString helpers for database storage
func (w WeekendPolicy) ToJSONString() string {
	data, _ := json.Marshal(w)
	return string(data)
}

func (n NewsletterConfig) ToJSONString() string {
	data, _ := json.Marshal(n)
	return string(data)
}

// Parse helpers for database retrieval
func ParseWeekendPolicy(data string) WeekendPolicy {
	var policy WeekendPolicy
	if err := json.Unmarshal([]byte(data), &policy); err != nil {
		return DefaultWeekendPolicy()
	}
	return policy
}

func ParseNewsletterConfig(data string) NewsletterConfig {
	var config NewsletterConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return DefaultNewsletterConfig()
	}
	return config
}
```

---

## Repository Layer

### Task 5.1: Create User Repository

- [ ] **Create user.go** `internal/repository/sqlite/user.go`

```go
package sqlite

import (
	"context"
	"database/sql"
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
	query := `
		INSERT INTO users (id, email, password_hash, name, role, vacation_balance, start_date, email_preferences)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.VacationBalance,
		user.StartDate,
		user.EmailPreferences.ToJSONString(),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users WHERE id = ?
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
		FROM users WHERE email = ?
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *UserRepository) List(ctx context.Context, role *domain.Role, search string, limit, offset int) ([]*domain.User, int, error) {
	// Build query with filters
	baseQuery := "FROM users WHERE 1=1"
	args := []interface{}{}

	if role != nil {
		baseQuery += " AND role = ?"
		args = append(args, *role)
	}

	if search != "" {
		baseQuery += " AND (name LIKE ? OR email LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	selectQuery := `
		SELECT id, email, password_hash, name, role, vacation_balance, start_date, email_preferences, created_at, updated_at
	` + baseQuery + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user, err := r.scanUserRow(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = ?, name = ?, role = ?, vacation_balance = ?, start_date = ?, email_preferences = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Name,
		user.Role,
		user.VacationBalance,
		user.StartDate,
		user.EmailPreferences.ToJSONString(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	query := "UPDATE users SET password_hash = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, passwordHash, id)
	return err
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id string, balance int) error {
	query := "UPDATE users SET vacation_balance = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, balance, id)
	return err
}

func (r *UserRepository) UpdateEmailPreferences(ctx context.Context, id string, prefs domain.EmailPreferences) error {
	query := "UPDATE users SET email_preferences = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, prefs.ToJSONString(), id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) CountByRole(ctx context.Context, role domain.Role) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE role = ?", role).Scan(&count)
	return count, err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string, excludeID string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	args := []interface{}{email}

	if excludeID != "" {
		query += " AND id != ?"
		args = append(args, excludeID)
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count > 0, err
}

// Helper functions
func (r *UserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var startDate sql.NullString
	var emailPrefsJSON string
	var createdAt, updatedAt string

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.VacationBalance,
		&startDate,
		&emailPrefsJSON,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	if startDate.Valid {
		user.StartDate = &startDate.String
	}
	user.EmailPreferences = domain.ParseEmailPreferences(emailPrefsJSON)
	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

func (r *UserRepository) scanUserRow(rows *sql.Rows) (*domain.User, error) {
	var user domain.User
	var startDate sql.NullString
	var emailPrefsJSON string
	var createdAt, updatedAt string

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.VacationBalance,
		&startDate,
		&emailPrefsJSON,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	if startDate.Valid {
		user.StartDate = &startDate.String
	}
	user.EmailPreferences = domain.ParseEmailPreferences(emailPrefsJSON)
	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}
```

---

### Task 5.2: Create Vacation Repository

- [ ] **Create vacation.go** `internal/repository/sqlite/vacation.go`

```go
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
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
		INSERT INTO vacation_requests (id, user_id, start_date, end_date, total_days, reason, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		req.ID,
		req.UserID,
		req.StartDate,
		req.EndDate,
		req.TotalDays,
		req.Reason,
		req.Status,
	)
	return err
}

func (r *VacationRepository) GetByID(ctx context.Context, id string) (*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.id = ?
	`
	return r.scanRequest(r.db.QueryRowContext(ctx, query, id))
}

func (r *VacationRepository) ListByUser(ctx context.Context, userID string, status *domain.VacationStatus, year *int) ([]*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.user_id = ?
	`
	args := []interface{}{userID}

	if status != nil {
		query += " AND vr.status = ?"
		args = append(args, *status)
	}

	if year != nil {
		query += " AND strftime('%Y', vr.start_date) = ?"
		args = append(args, fmt.Sprintf("%d", *year))
	}

	query += " ORDER BY vr.created_at DESC"

	return r.queryRequests(ctx, query, args...)
}

func (r *VacationRepository) ListPending(ctx context.Context) ([]*domain.VacationRequest, error) {
	query := `
		SELECT vr.id, vr.user_id, u.name, u.email, vr.start_date, vr.end_date, vr.total_days,
		       vr.reason, vr.status, vr.reviewed_by, vr.reviewed_at, vr.rejection_reason,
		       vr.created_at, vr.updated_at
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.status = 'pending'
		ORDER BY vr.created_at ASC
	`
	return r.queryRequests(ctx, query)
}

func (r *VacationRepository) ListTeam(ctx context.Context, month, year int) ([]*domain.TeamVacation, error) {
	// Get start and end of month
	startOfMonth := fmt.Sprintf("%d-%02d-01", year, month)
	endOfMonth := fmt.Sprintf("%d-%02d-31", year, month)

	query := `
		SELECT vr.id, vr.user_id, u.name, vr.start_date, vr.end_date, vr.total_days
		FROM vacation_requests vr
		JOIN users u ON vr.user_id = u.id
		WHERE vr.status = 'approved'
		AND (
			(vr.start_date >= ? AND vr.start_date <= ?)
			OR (vr.end_date >= ? AND vr.end_date <= ?)
			OR (vr.start_date <= ? AND vr.end_date >= ?)
		)
		ORDER BY vr.start_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query,
		startOfMonth, endOfMonth,
		startOfMonth, endOfMonth,
		startOfMonth, endOfMonth,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacations []*domain.TeamVacation
	for rows.Next() {
		var v domain.TeamVacation
		if err := rows.Scan(&v.ID, &v.UserID, &v.UserName, &v.StartDate, &v.EndDate, &v.TotalDays); err != nil {
			return nil, err
		}
		vacations = append(vacations, &v)
	}

	return vacations, nil
}

func (r *VacationRepository) UpdateStatus(ctx context.Context, id string, status domain.VacationStatus, reviewedBy string, rejectionReason *string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE vacation_requests
		SET status = ?, reviewed_by = ?, reviewed_at = ?, rejection_reason = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, reviewedBy, now, rejectionReason, id)
	return err
}

func (r *VacationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM vacation_requests WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("request not found")
	}
	return nil
}

// Helper functions
func (r *VacationRepository) scanRequest(row *sql.Row) (*domain.VacationRequest, error) {
	var req domain.VacationRequest
	var reason, reviewedBy, rejectionReason sql.NullString
	var reviewedAt sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&req.ID,
		&req.UserID,
		&req.UserName,
		&req.UserEmail,
		&req.StartDate,
		&req.EndDate,
		&req.TotalDays,
		&reason,
		&req.Status,
		&reviewedBy,
		&reviewedAt,
		&rejectionReason,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if reason.Valid {
		req.Reason = &reason.String
	}
	if reviewedBy.Valid {
		req.ReviewedBy = &reviewedBy.String
	}
	if reviewedAt.Valid {
		t, _ := time.Parse(time.RFC3339, reviewedAt.String)
		req.ReviewedAt = &t
	}
	if rejectionReason.Valid {
		req.RejectionReason = &rejectionReason.String
	}
	req.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	req.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &req, nil
}

func (r *VacationRepository) queryRequests(ctx context.Context, query string, args ...interface{}) ([]*domain.VacationRequest, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.VacationRequest
	for rows.Next() {
		var req domain.VacationRequest
		var reason, reviewedBy, rejectionReason sql.NullString
		var reviewedAt sql.NullString
		var createdAt, updatedAt string

		err := rows.Scan(
			&req.ID,
			&req.UserID,
			&req.UserName,
			&req.UserEmail,
			&req.StartDate,
			&req.EndDate,
			&req.TotalDays,
			&reason,
			&req.Status,
			&reviewedBy,
			&reviewedAt,
			&rejectionReason,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reason.Valid {
			req.Reason = &reason.String
		}
		if reviewedBy.Valid {
			req.ReviewedBy = &reviewedBy.String
		}
		if reviewedAt.Valid {
			t, _ := time.Parse(time.RFC3339, reviewedAt.String)
			req.ReviewedAt = &t
		}
		if rejectionReason.Valid {
			req.RejectionReason = &rejectionReason.String
		}
		req.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		req.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		requests = append(requests, &req)
	}

	return requests, nil
}
```

---

### Task 5.3: Create Settings Repository

- [ ] **Create settings.go** `internal/repository/sqlite/settings.go`

```go
package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

type SettingsRepository struct {
	db *DB
}

func NewSettingsRepository(db *DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context) (*domain.Settings, error) {
	query := `
		SELECT id, weekend_policy, newsletter, default_vacation_days, vacation_reset_month, updated_at
		FROM settings
		WHERE id = 'settings'
	`

	var settings domain.Settings
	var weekendPolicyJSON, newsletterJSON string
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query).Scan(
		&settings.ID,
		&weekendPolicyJSON,
		&newsletterJSON,
		&settings.DefaultVacationDays,
		&settings.VacationResetMonth,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		// Return default settings if none exist
		return domain.DefaultSettings(), nil
	}
	if err != nil {
		return nil, err
	}

	settings.WeekendPolicy = domain.ParseWeekendPolicy(weekendPolicyJSON)
	settings.Newsletter = domain.ParseNewsletterConfig(newsletterJSON)
	settings.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &settings, nil
}

func (r *SettingsRepository) Update(ctx context.Context, settings *domain.Settings) error {
	query := `
		INSERT INTO settings (id, weekend_policy, newsletter, default_vacation_days, vacation_reset_month)
		VALUES ('settings', ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			weekend_policy = excluded.weekend_policy,
			newsletter = excluded.newsletter,
			default_vacation_days = excluded.default_vacation_days,
			vacation_reset_month = excluded.vacation_reset_month
	`

	_, err := r.db.ExecContext(ctx, query,
		settings.WeekendPolicy.ToJSONString(),
		settings.Newsletter.ToJSONString(),
		settings.DefaultVacationDays,
		settings.VacationResetMonth,
	)
	return err
}
```

---

## Service Layer

### Task 6.1: Create Auth Service

- [ ] **Create auth.go** `internal/service/auth.go`

```go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourorg/vacaytracker-api/internal/config"
	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo *sqlite.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *sqlite.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// VerifyPassword checks if password matches hash
func (s *AuthService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken creates a JWT token for a user
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// ValidateToken verifies and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*dto.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &dto.JWTClaims{
			UserID: claims["sub"].(string),
			Email:  claims["email"].(string),
			Name:   claims["name"].(string),
			Role:   domain.Role(claims["role"].(string)),
		}, nil
	}

	return nil, errors.New("invalid token")
}

// Login authenticates user and returns token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}

	if !s.VerifyPassword(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// ChangePassword updates user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if !s.VerifyPassword(currentPassword, user.PasswordHash) {
		return ErrInvalidCredentials
	}

	newHash, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, newHash)
}

// CreateInitialAdmin creates the admin user if none exists
func (s *AuthService) CreateInitialAdmin(ctx context.Context) error {
	count, err := s.userRepo.CountByRole(ctx, domain.RoleAdmin)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Admin already exists
	}

	hash, err := s.HashPassword(s.cfg.AdminPassword)
	if err != nil {
		return err
	}

	admin := &domain.User{
		ID:               uuid.New().String(),
		Email:            s.cfg.AdminEmail,
		PasswordHash:     hash,
		Name:             s.cfg.AdminName,
		Role:             domain.RoleAdmin,
		VacationBalance:  25,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	return s.userRepo.Create(ctx, admin)
}
```

---

### Task 6.2: Create Vacation Service

- [ ] **Create vacation.go** `internal/service/vacation.go`

```go
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
)

var (
	ErrInsufficientBalance    = errors.New("insufficient vacation balance")
	ErrInvalidDateRange       = errors.New("end date must be after start date")
	ErrDateInPast             = errors.New("start date cannot be in the past")
	ErrRequestNotFound        = errors.New("vacation request not found")
	ErrCannotCancelApproved   = errors.New("cannot cancel approved request")
	ErrCannotCancelRejected   = errors.New("cannot cancel rejected request")
	ErrAlreadyProcessed       = errors.New("request already processed")
)

type VacationService struct {
	vacationRepo *sqlite.VacationRepository
	userRepo     *sqlite.UserRepository
	settingsRepo *sqlite.SettingsRepository
}

func NewVacationService(
	vacationRepo *sqlite.VacationRepository,
	userRepo *sqlite.UserRepository,
	settingsRepo *sqlite.SettingsRepository,
) *VacationService {
	return &VacationService{
		vacationRepo: vacationRepo,
		userRepo:     userRepo,
		settingsRepo: settingsRepo,
	}
}

// Create creates a new vacation request
func (s *VacationService) Create(ctx context.Context, userID string, req dto.CreateVacationRequest) (*domain.VacationRequest, error) {
	// Parse dates (DD/MM/YYYY -> YYYY-MM-DD)
	startDate, err := parseDDMMYYYY(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := parseDDMMYYYY(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}

	// Check if start date is in the past
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if startDate.Before(today) {
		return nil, ErrDateInPast
	}

	// Get settings for business day calculation
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate business days
	totalDays := calculateBusinessDays(startDate, endDate, settings.WeekendPolicy)

	// Get user and check balance
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.VacationBalance < totalDays {
		return nil, ErrInsufficientBalance
	}

	// Create request
	vacation := &domain.VacationRequest{
		ID:        uuid.New().String(),
		UserID:    userID,
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		TotalDays: totalDays,
		Status:    domain.StatusPending,
	}

	if req.Reason != nil && *req.Reason != "" {
		vacation.Reason = req.Reason
	}

	if err := s.vacationRepo.Create(ctx, vacation); err != nil {
		return nil, err
	}

	return vacation, nil
}

// Cancel cancels a pending vacation request
func (s *VacationService) Cancel(ctx context.Context, requestID, userID string) error {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return ErrRequestNotFound
	}

	// Check ownership
	if request.UserID != userID {
		return errors.New("access denied")
	}

	// Check status
	if request.IsApproved() {
		return ErrCannotCancelApproved
	}
	if request.IsRejected() {
		return ErrCannotCancelRejected
	}

	return s.vacationRepo.Delete(ctx, requestID)
}

// Approve approves a pending request and deducts balance
func (s *VacationService) Approve(ctx context.Context, requestID, adminID string) (*domain.VacationRequest, error) {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrRequestNotFound
	}

	if !request.IsPending() {
		return nil, ErrAlreadyProcessed
	}

	// Update status
	if err := s.vacationRepo.UpdateStatus(ctx, requestID, domain.StatusApproved, adminID, nil); err != nil {
		return nil, err
	}

	// Deduct vacation balance
	user, err := s.userRepo.GetByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	newBalance := user.VacationBalance - request.TotalDays
	if newBalance < 0 {
		newBalance = 0
	}

	if err := s.userRepo.UpdateBalance(ctx, request.UserID, newBalance); err != nil {
		return nil, err
	}

	// Fetch updated request
	return s.vacationRepo.GetByID(ctx, requestID)
}

// Reject rejects a pending request
func (s *VacationService) Reject(ctx context.Context, requestID, adminID string, reason *string) (*domain.VacationRequest, error) {
	request, err := s.vacationRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrRequestNotFound
	}

	if !request.IsPending() {
		return nil, ErrAlreadyProcessed
	}

	if err := s.vacationRepo.UpdateStatus(ctx, requestID, domain.StatusRejected, adminID, reason); err != nil {
		return nil, err
	}

	return s.vacationRepo.GetByID(ctx, requestID)
}

// Helper functions

// parseDDMMYYYY parses DD/MM/YYYY format to time.Time
func parseDDMMYYYY(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "/")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid date format, expected DD/MM/YYYY")
	}

	// Rearrange to YYYY-MM-DD
	isoDate := fmt.Sprintf("%s-%s-%s", parts[2], parts[1], parts[0])
	return time.Parse("2006-01-02", isoDate)
}

// calculateBusinessDays counts business days between two dates
func calculateBusinessDays(start, end time.Time, policy domain.WeekendPolicy) int {
	if !policy.ExcludeWeekends {
		// Include all days
		return int(end.Sub(start).Hours()/24) + 1
	}

	count := 0
	current := start

	// Create a map of excluded weekdays for faster lookup
	excluded := make(map[time.Weekday]bool)
	for _, day := range policy.ExcludedDays {
		excluded[time.Weekday(day)] = true
	}

	for !current.After(end) {
		if !excluded[current.Weekday()] {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}

	return count
}
```

---

### Task 6.3: Create User Service

- [ ] **Create user.go** `internal/service/user.go`

```go
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
)

var (
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrCannotDeleteSelf       = errors.New("cannot delete own account")
	ErrCannotDeleteLastAdmin  = errors.New("cannot delete last admin")
	ErrCannotRemoveLastAdmin  = errors.New("cannot demote last admin")
	ErrCannotModifyOwnRole    = errors.New("cannot modify own role")
)

type UserService struct {
	userRepo    *sqlite.UserRepository
	authService *AuthService
}

func NewUserService(userRepo *sqlite.UserRepository, authService *AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*domain.User, error) {
	// Check if email exists
	exists, err := s.userRepo.EmailExists(ctx, req.Email, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hash, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set defaults
	balance := 25
	if req.VacationBalance != nil {
		balance = *req.VacationBalance
	}

	user := &domain.User{
		ID:               uuid.New().String(),
		Email:            req.Email,
		PasswordHash:     hash,
		Name:             req.Name,
		Role:             domain.Role(req.Role),
		VacationBalance:  balance,
		StartDate:        req.StartDate,
		EmailPreferences: domain.DefaultEmailPreferences(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user's information
func (s *UserService) Update(ctx context.Context, id string, req dto.UpdateUserRequest, currentUserID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check email uniqueness if changing
	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.userRepo.EmailExists(ctx, *req.Email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}

	// Check role change restrictions
	if req.Role != nil && domain.Role(*req.Role) != user.Role {
		// Cannot modify own role
		if id == currentUserID {
			return nil, ErrCannotModifyOwnRole
		}

		// Cannot demote if last admin
		if user.Role == domain.RoleAdmin && *req.Role == string(domain.RoleEmployee) {
			count, err := s.userRepo.CountByRole(ctx, domain.RoleAdmin)
			if err != nil {
				return nil, err
			}
			if count <= 1 {
				return nil, ErrCannotRemoveLastAdmin
			}
		}

		user.Role = domain.Role(*req.Role)
	}

	// Update other fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.VacationBalance != nil {
		user.VacationBalance = *req.VacationBalance
	}
	if req.StartDate != nil {
		user.StartDate = req.StartDate
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id, currentUserID string) error {
	// Cannot delete self
	if id == currentUserID {
		return ErrCannotDeleteSelf
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Cannot delete last admin
	if user.Role == domain.RoleAdmin {
		count, err := s.userRepo.CountByRole(ctx, domain.RoleAdmin)
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrCannotDeleteLastAdmin
		}
	}

	return s.userRepo.Delete(ctx, id)
}
```

---

## Middleware

### Task 7.1: Create Auth Middleware

- [ ] **Create auth.go** `internal/middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAuth validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTH_TOKEN_MISSING",
				Message: "Authorization header required",
			})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTH_TOKEN_INVALID",
				Message: "Invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTH_TOKEN_INVALID",
				Message: "Invalid or expired token",
			})
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userName", claims.Name)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// RequireAdmin checks that user has admin role
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "AUTH_TOKEN_MISSING",
				Message: "Authentication required",
			})
			return
		}

		if role.(domain.Role) != domain.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    "ADMIN_REQUIRED",
				Message: "Admin role required",
			})
			return
		}

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) string {
	if id, exists := c.Get("userID"); exists {
		return id.(string)
	}
	return ""
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) domain.Role {
	if role, exists := c.Get("userRole"); exists {
		return role.(domain.Role)
	}
	return domain.RoleEmployee
}
```

---

### Task 7.2: Create CORS Middleware

- [ ] **Create cors.go** `internal/middleware/cors.go`

```go
package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORS(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Allow specified origin or localhost in development
		if origin == allowOrigin ||
		   origin == "http://localhost:5173" ||
		   origin == "http://localhost:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
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

---

### Task 7.3: Create Error Handler Middleware

- [ ] **Create error.go** `internal/middleware/error.go`

```go
package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/dto"
)

// ErrorHandler recovers from panics and handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
				})
			}
		}()

		c.Next()

		// Check for errors set during request handling
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Request error: %v", err)

			if !c.Writer.Written() {
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
				})
			}
		}
	}
}
```

---

## Handlers

### Task 8.1: Create Health Handler

- [ ] **Create health.go** `internal/handler/health.go`

```go
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}
```

---

### Task 8.2: Create Auth Handler

- [ ] **Create auth.go** `internal/handler/auth.go`

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/middleware"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *sqlite.UserRepository
}

func NewAuthHandler(authService *service.AuthService, userRepo *sqlite.UserRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			})
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Login failed",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  dto.UserToResponse(user),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get user",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserToResponse(user))
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	// Validate password length
	if len(req.NewPassword) < 6 || len(req.NewPassword) > 72 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Password must be 6-72 characters",
		})
		return
	}

	userID := middleware.GetUserID(c)
	err := h.authService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "Current password is incorrect",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to change password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

func (h *AuthHandler) UpdateEmailPreferences(c *gin.Context) {
	var req dto.UpdateEmailPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	userID := middleware.GetUserID(c)
	err := h.userRepo.UpdateEmailPreferences(c.Request.Context(), userID, req.ToEmailPreferences())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"emailPreferences": req,
	})
}
```

---

### Task 8.3: Create Vacation Handler

- [ ] **Create vacation.go** `internal/handler/vacation.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/middleware"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

type VacationHandler struct {
	vacationService *service.VacationService
	vacationRepo    *sqlite.VacationRepository
}

func NewVacationHandler(vacationService *service.VacationService, vacationRepo *sqlite.VacationRepository) *VacationHandler {
	return &VacationHandler{
		vacationService: vacationService,
		vacationRepo:    vacationRepo,
	}
}

func (h *VacationHandler) Create(c *gin.Context) {
	var req dto.CreateVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	userID := middleware.GetUserID(c)
	vacation, err := h.vacationService.Create(c.Request.Context(), userID, req)
	if err != nil {
		switch err {
		case service.ErrInsufficientBalance:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "INSUFFICIENT_BALANCE",
				Message: "Insufficient vacation balance",
			})
		case service.ErrInvalidDateRange:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    "INVALID_DATE_RANGE",
				Message: "End date must be after start date",
			})
		case service.ErrDateInPast:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Code:    "DATE_IN_PAST",
				Message: "Start date cannot be in the past",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create request",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, vacation)
}

func (h *VacationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// Parse query parameters
	var status *domain.VacationStatus
	if s := c.Query("status"); s != "" {
		vs := domain.VacationStatus(s)
		status = &vs
	}

	var year *int
	if y := c.Query("year"); y != "" {
		if parsed, err := strconv.Atoi(y); err == nil {
			year = &parsed
		}
	}

	requests, err := h.vacationRepo.ListByUser(c.Request.Context(), userID, status, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"total":    len(requests),
	})
}

func (h *VacationHandler) Get(c *gin.Context) {
	requestID := c.Param("id")
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	request, err := h.vacationRepo.GetByID(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get request",
		})
		return
	}
	if request == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    "REQUEST_NOT_FOUND",
			Message: "Vacation request not found",
		})
		return
	}

	// Check access (employee can only see own requests)
	if userRole != domain.RoleAdmin && request.UserID != userID {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Code:    "ACCESS_DENIED",
			Message: "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, request)
}

func (h *VacationHandler) Cancel(c *gin.Context) {
	requestID := c.Param("id")
	userID := middleware.GetUserID(c)

	err := h.vacationService.Cancel(c.Request.Context(), requestID, userID)
	if err != nil {
		switch err {
		case service.ErrRequestNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "REQUEST_NOT_FOUND",
				Message: "Vacation request not found",
			})
		case service.ErrCannotCancelApproved:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_CANCEL_APPROVED",
				Message: "Cannot cancel approved request",
			})
		case service.ErrCannotCancelRejected:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_CANCEL_REJECTED",
				Message: "Cannot cancel rejected request",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to cancel request",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *VacationHandler) Team(c *gin.Context) {
	// Parse query parameters
	month, err := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid month parameter",
		})
		return
	}

	year, err := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid year parameter",
		})
		return
	}

	vacations, err := h.vacationRepo.ListTeam(c.Request.Context(), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get team calendar",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vacations": vacations,
		"month":     month,
		"year":      year,
	})
}

// Import time package at top
import "time"
```

---

### Task 8.4: Create Admin Handler

- [ ] **Create admin.go** `internal/handler/admin.go`

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yourorg/vacaytracker-api/internal/domain"
	"github.com/yourorg/vacaytracker-api/internal/dto"
	"github.com/yourorg/vacaytracker-api/internal/middleware"
	"github.com/yourorg/vacaytracker-api/internal/repository/sqlite"
	"github.com/yourorg/vacaytracker-api/internal/service"
)

type AdminHandler struct {
	userService     *service.UserService
	vacationService *service.VacationService
	userRepo        *sqlite.UserRepository
	vacationRepo    *sqlite.VacationRepository
	settingsRepo    *sqlite.SettingsRepository
}

func NewAdminHandler(
	userService *service.UserService,
	vacationService *service.VacationService,
	userRepo *sqlite.UserRepository,
	vacationRepo *sqlite.VacationRepository,
	settingsRepo *sqlite.SettingsRepository,
) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		vacationService: vacationService,
		userRepo:        userRepo,
		vacationRepo:    vacationRepo,
		settingsRepo:    settingsRepo,
	}
}

// User Management

func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Parse filters
	var role *domain.Role
	if r := c.Query("role"); r != "" {
		roleVal := domain.Role(r)
		role = &roleVal
	}
	search := c.Query("search")

	users, total, err := h.userRepo.List(c.Request.Context(), role, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": dto.UsersToResponse(users),
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req)
	if err != nil {
		switch err {
		case service.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    "EMAIL_ALREADY_EXISTS",
				Message: "Email already registered",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create user",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.UserToResponse(user))
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get user",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserToResponse(user))
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, req, currentUserID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			})
		case service.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Code:    "EMAIL_ALREADY_EXISTS",
				Message: "Email already taken",
			})
		case service.ErrCannotRemoveLastAdmin:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_REMOVE_LAST_ADMIN",
				Message: "Cannot demote last admin",
			})
		case service.ErrCannotModifyOwnRole:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_MODIFY_OWN_ROLE",
				Message: "Cannot modify own role",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to update user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.UserToResponse(user))
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)

	err := h.userService.Delete(c.Request.Context(), userID, currentUserID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			})
		case service.ErrCannotDeleteSelf:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_DELETE_SELF",
				Message: "Cannot delete own account",
			})
		case service.ErrCannotDeleteLastAdmin:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "CANNOT_DELETE_LAST_ADMIN",
				Message: "Cannot delete last admin",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to delete user",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// Vacation Management

func (h *AdminHandler) PendingRequests(c *gin.Context) {
	requests, err := h.vacationRepo.ListPending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list pending requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
		"total":    len(requests),
	})
}

func (h *AdminHandler) ApproveRequest(c *gin.Context) {
	requestID := c.Param("id")
	adminID := middleware.GetUserID(c)

	request, err := h.vacationService.Approve(c.Request.Context(), requestID, adminID)
	if err != nil {
		switch err {
		case service.ErrRequestNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "REQUEST_NOT_FOUND",
				Message: "Vacation request not found",
			})
		case service.ErrAlreadyProcessed:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "REQUEST_ALREADY_PROCESSED",
				Message: "Request already processed",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to approve request",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         request.ID,
		"status":     request.Status,
		"reviewedBy": request.ReviewedBy,
		"reviewedAt": request.ReviewedAt,
	})
}

func (h *AdminHandler) RejectRequest(c *gin.Context) {
	requestID := c.Param("id")
	adminID := middleware.GetUserID(c)

	var req dto.RejectVacationRequest
	c.ShouldBindJSON(&req) // Optional body

	request, err := h.vacationService.Reject(c.Request.Context(), requestID, adminID, req.Reason)
	if err != nil {
		switch err {
		case service.ErrRequestNotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    "REQUEST_NOT_FOUND",
				Message: "Vacation request not found",
			})
		case service.ErrAlreadyProcessed:
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Code:    "REQUEST_ALREADY_PROCESSED",
				Message: "Request already processed",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to reject request",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              request.ID,
		"status":          request.Status,
		"rejectionReason": request.RejectionReason,
		"reviewedBy":      request.ReviewedBy,
		"reviewedAt":      request.ReviewedAt,
	})
}

// Settings

func (h *AdminHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsRepo.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get settings",
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
		})
		return
	}

	settings := req.ToSettings()
	if err := h.settingsRepo.Update(c.Request.Context(), settings); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update settings",
		})
		return
	}

	// Fetch updated settings
	updated, _ := h.settingsRepo.Get(c.Request.Context())
	c.JSON(http.StatusOK, updated)
}
```

---

## DTOs

### Task 9.1: Create Request DTOs

- [ ] **Create request.go** `internal/dto/request.go`

```go
package dto

import "github.com/yourorg/vacaytracker-api/internal/domain"

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=72"`
}

type UpdateEmailPreferencesRequest struct {
	VacationUpdates   *bool `json:"vacationUpdates"`
	WeeklyDigest      *bool `json:"weeklyDigest"`
	TeamNotifications *bool `json:"teamNotifications"`
}

func (r *UpdateEmailPreferencesRequest) ToEmailPreferences() domain.EmailPreferences {
	prefs := domain.DefaultEmailPreferences()
	if r.VacationUpdates != nil {
		prefs.VacationUpdates = *r.VacationUpdates
	}
	if r.WeeklyDigest != nil {
		prefs.WeeklyDigest = *r.WeeklyDigest
	}
	if r.TeamNotifications != nil {
		prefs.TeamNotifications = *r.TeamNotifications
	}
	return prefs
}

// User DTOs
type CreateUserRequest struct {
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=6,max=72"`
	Name            string  `json:"name" binding:"required,min=1,max=100"`
	Role            string  `json:"role" binding:"required,oneof=admin employee"`
	VacationBalance *int    `json:"vacationBalance"`
	StartDate       *string `json:"startDate"` // DD/MM/YYYY
}

type UpdateUserRequest struct {
	Email           *string `json:"email" binding:"omitempty,email"`
	Name            *string `json:"name" binding:"omitempty,min=1,max=100"`
	Role            *string `json:"role" binding:"omitempty,oneof=admin employee"`
	VacationBalance *int    `json:"vacationBalance"`
	StartDate       *string `json:"startDate"`
}

// Vacation DTOs
type CreateVacationRequest struct {
	StartDate string  `json:"startDate" binding:"required"` // DD/MM/YYYY
	EndDate   string  `json:"endDate" binding:"required"`   // DD/MM/YYYY
	Reason    *string `json:"reason" binding:"omitempty,max=500"`
}

type RejectVacationRequest struct {
	Reason *string `json:"reason"`
}

// Settings DTOs
type UpdateSettingsRequest struct {
	WeekendPolicy       *domain.WeekendPolicy    `json:"weekendPolicy"`
	Newsletter          *domain.NewsletterConfig `json:"newsletter"`
	DefaultVacationDays *int                     `json:"defaultVacationDays"`
	VacationResetMonth  *int                     `json:"vacationResetMonth"`
}

func (r *UpdateSettingsRequest) ToSettings() *domain.Settings {
	settings := domain.DefaultSettings()
	if r.WeekendPolicy != nil {
		settings.WeekendPolicy = *r.WeekendPolicy
	}
	if r.Newsletter != nil {
		settings.Newsletter = *r.Newsletter
	}
	if r.DefaultVacationDays != nil {
		settings.DefaultVacationDays = *r.DefaultVacationDays
	}
	if r.VacationResetMonth != nil {
		settings.VacationResetMonth = *r.VacationResetMonth
	}
	return settings
}
```

---

### Task 9.2: Create Response DTOs

- [ ] **Create response.go** `internal/dto/response.go`

```go
package dto

import (
	"github.com/yourorg/vacaytracker-api/internal/domain"
)

// Auth Response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type JWTClaims struct {
	UserID string
	Email  string
	Name   string
	Role   domain.Role
}

// User Response
type UserResponse struct {
	ID               string                  `json:"id"`
	Email            string                  `json:"email"`
	Name             string                  `json:"name"`
	Role             domain.Role             `json:"role"`
	VacationBalance  int                     `json:"vacationBalance"`
	StartDate        *string                 `json:"startDate,omitempty"`
	EmailPreferences domain.EmailPreferences `json:"emailPreferences"`
	CreatedAt        string                  `json:"createdAt,omitempty"`
	UpdatedAt        string                  `json:"updatedAt,omitempty"`
}

func UserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:               user.ID,
		Email:            user.Email,
		Name:             user.Name,
		Role:             user.Role,
		VacationBalance:  user.VacationBalance,
		StartDate:        user.StartDate,
		EmailPreferences: user.EmailPreferences,
		CreatedAt:        user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func UsersToResponse(users []*domain.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = UserToResponse(user)
	}
	return responses
}
```

---

### Task 9.3: Create Error DTOs

- [ ] **Create errors.go** `internal/dto/errors.go`

```go
package dto

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Common error codes
const (
	// Auth errors
	ErrCodeAuthTokenMissing   = "AUTH_TOKEN_MISSING"
	ErrCodeAuthTokenInvalid   = "AUTH_TOKEN_INVALID"
	ErrCodeAuthTokenExpired   = "AUTH_TOKEN_EXPIRED"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeAccessDenied       = "ACCESS_DENIED"
	ErrCodeAdminRequired      = "ADMIN_REQUIRED"

	// Validation errors
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeInvalidDateFormat = "INVALID_DATE_FORMAT"
	ErrCodeInvalidDateRange  = "INVALID_DATE_RANGE"
	ErrCodeDateInPast        = "DATE_IN_PAST"

	// Resource errors
	ErrCodeUserNotFound    = "USER_NOT_FOUND"
	ErrCodeRequestNotFound = "REQUEST_NOT_FOUND"

	// Conflict errors
	ErrCodeEmailAlreadyExists = "EMAIL_ALREADY_EXISTS"

	// Business rule errors
	ErrCodeInsufficientBalance    = "INSUFFICIENT_BALANCE"
	ErrCodeCannotCancelApproved   = "CANNOT_CANCEL_APPROVED"
	ErrCodeCannotCancelRejected   = "CANNOT_CANCEL_REJECTED"
	ErrCodeRequestAlreadyProcessed = "REQUEST_ALREADY_PROCESSED"
	ErrCodeCannotDeleteSelf       = "CANNOT_DELETE_SELF"
	ErrCodeCannotDeleteLastAdmin  = "CANNOT_DELETE_LAST_ADMIN"
	ErrCodeCannotRemoveLastAdmin  = "CANNOT_REMOVE_LAST_ADMIN"
	ErrCodeCannotModifyOwnRole    = "CANNOT_MODIFY_OWN_ROLE"

	// Server errors
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeDatabaseError = "DATABASE_ERROR"
)
```

---

### Task 10: Create Main Entry Point

- [ ] **Create main.go** `cmd/server/main.go`

```go
package main

import (
	"context"
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
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := sqlite.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations("./migrations"); err != nil {
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

	// Create initial admin user
	if err := authService.CreateInitialAdmin(context.Background()); err != nil {
		log.Printf("Warning: Failed to create initial admin: %v", err)
	}

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService, userRepo)
	vacationHandler := handler.NewVacationHandler(vacationService, vacationRepo)
	adminHandler := handler.NewAdminHandler(userService, vacationService, userRepo, vacationRepo, settingsRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS(cfg.URL))

	// Public routes
	r.GET("/health", healthHandler.Check)

	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Authenticated routes
		authenticated := api.Group("/")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			authenticated.GET("/auth/me", authHandler.Me)
			authenticated.PUT("/auth/password", authHandler.ChangePassword)
			authenticated.PUT("/auth/email-preferences", authHandler.UpdateEmailPreferences)

			// Employee vacation routes
			vacation := authenticated.Group("/vacation")
			{
				vacation.POST("/request", vacationHandler.Create)
				vacation.GET("/requests", vacationHandler.List)
				vacation.GET("/requests/:id", vacationHandler.Get)
				vacation.DELETE("/requests/:id", vacationHandler.Cancel)
				vacation.GET("/team", vacationHandler.Team)
			}
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireAdmin())
		{
			// User management
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.GET("/users/:id", adminHandler.GetUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// Vacation management
			admin.GET("/vacation/pending", adminHandler.PendingRequests)
			admin.PUT("/vacation/:id/approve", adminHandler.ApproveRequest)
			admin.PUT("/vacation/:id/reject", adminHandler.RejectRequest)

			// Settings
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)
		}
	}

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
```

---

## Testing

### Task 10.1: Create Auth Service Tests

- [ ] **Create auth_test.go** `internal/service/auth_test.go`

```go
package service

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	// Test implementation
}

func TestVerifyPassword(t *testing.T) {
	// Test implementation
}

func TestGenerateToken(t *testing.T) {
	// Test implementation
}

func TestValidateToken(t *testing.T) {
	// Test implementation
}
```

### Task 10.2: Create Vacation Service Tests

- [ ] **Create vacation_test.go** `internal/service/vacation_test.go`

```go
package service

import (
	"testing"

	"github.com/yourorg/vacaytracker-api/internal/domain"
)

func TestCalculateBusinessDays(t *testing.T) {
	tests := []struct {
		name     string
		start    string
		end      string
		policy   domain.WeekendPolicy
		expected int
	}{
		{
			name:  "weekdays only, exclude weekends",
			start: "2024-01-15", // Monday
			end:   "2024-01-19", // Friday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			expected: 5,
		},
		{
			name:  "include weekend days",
			start: "2024-01-15", // Monday
			end:   "2024-01-21", // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: false,
			},
			expected: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}

func TestParseDDMMYYYY(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"15/01/2024", "2024-01-15", false},
		{"01/12/2024", "2024-12-01", false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

---

## Related Documents

- [01-database-schema.md](./01-database-schema.md) - Database design details
- [02-api-specification.md](./02-api-specification.md) - API endpoint details
- [03-implementation-roadmap.md](./03-implementation-roadmap.md) - Phase dependencies
- [07-testing-strategy.md](./07-testing-strategy.md) - Test coverage requirements
