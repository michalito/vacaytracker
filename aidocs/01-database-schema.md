# VacayTracker Database Schema

> **Purpose:** Complete database schema documentation for AI agent implementation
> **Database:** SQLite with modernc.org/sqlite driver (CGo-free)
> **Format:** Self-contained, explicit, copy-paste ready

---

## Table of Contents

1. [Entity Relationship Diagram](#1-entity-relationship-diagram)
2. [Table Definitions](#2-table-definitions)
3. [JSON Field Schemas](#3-json-field-schemas)
4. [Indexes](#4-indexes)
5. [Triggers](#5-triggers)
6. [Type Mappings](#6-type-mappings)
7. [Complete Migration SQL](#7-complete-migration-sql)
8. [Seed Data](#8-seed-data)

---

## 1. Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     VACAYTRACKER DATABASE SCHEMA                             │
│                     ────────────────────────────                             │
│                     SQLite with WAL Mode                                     │
└─────────────────────────────────────────────────────────────────────────────┘


    ┌──────────────────────────────────────────────────────────────────┐
    │                            users                                  │
    ├──────────────────────────────────────────────────────────────────┤
    │ PK   id                    TEXT         NOT NULL                 │
    │      name                  TEXT         NOT NULL                 │
    │ UQ   username              TEXT         NOT NULL  UNIQUE         │
    │      email                 TEXT         NULL                     │
    │      password              TEXT         NOT NULL  (bcrypt hash)  │
    │      role                  TEXT         NOT NULL  CHECK(admin|employee) │
    │      vacation_days         INTEGER      DEFAULT 0                │
    │      used_vacation_days    INTEGER      DEFAULT 0                │
    │      email_preferences     TEXT         NOT NULL  DEFAULT '{}'  (JSON) │
    │      created_at            DATETIME     NOT NULL  DEFAULT NOW    │
    │      updated_at            DATETIME     NOT NULL  DEFAULT NOW    │
    └────────────────────────────────┬─────────────────────────────────┘
                                     │
                                     │ 1:N
                                     │ ON DELETE CASCADE
                                     │
    ┌────────────────────────────────▼─────────────────────────────────┐
    │                      vacation_requests                            │
    ├──────────────────────────────────────────────────────────────────┤
    │ PK   id                    TEXT         NOT NULL                 │
    │ FK   user_id               TEXT         NOT NULL  → users.id     │
    │      start_date            TEXT         NOT NULL  (YYYY-MM-DD)   │
    │      end_date              TEXT         NOT NULL  (YYYY-MM-DD)   │
    │      business_days         INTEGER      NOT NULL                 │
    │      status                TEXT         NOT NULL  CHECK(pending|approved|rejected) │
    │      reason                TEXT         NULL      (max 200 chars)│
    │      created_at            DATETIME     NOT NULL  DEFAULT NOW    │
    │      reviewed_by           TEXT         NULL      → users.id     │
    │      reviewed_at           DATETIME     NULL                     │
    └──────────────────────────────────────────────────────────────────┘


    ┌──────────────────────────────────────────────────────────────────┐
    │                          settings                                 │
    ├──────────────────────────────────────────────────────────────────┤
    │ PK   id                    INTEGER      CHECK(id = 1)  (single row) │
    │      weekend_policy        TEXT         NOT NULL  (JSON)         │
    │      newsletter            TEXT         NOT NULL  (JSON)         │
    └──────────────────────────────────────────────────────────────────┘


    ═══════════════════════════════════════════════════════════════════
                              RELATIONSHIPS
    ═══════════════════════════════════════════════════════════════════

    users (1) ─────────────────── (N) vacation_requests
        │
        │   • One user can have zero or more vacation requests
        │   • Each vacation request belongs to exactly one user
        │   • CASCADE DELETE: Deleting user removes all their requests
        │
        └──► user_id (FK) references users.id

    users (1) ─────────────────── (0..1) vacation_requests.reviewed_by
        │
        │   • An admin may review multiple vacation requests
        │   • A vacation request may or may not have a reviewer
        │   • NULL until request is reviewed
        │
        └──► reviewed_by references users.id (nullable)

    ═══════════════════════════════════════════════════════════════════
                                INDEXES
    ═══════════════════════════════════════════════════════════════════

    idx_vacation_user    ON vacation_requests(user_id)
    idx_vacation_status  ON vacation_requests(status)
    idx_users_role       ON users(role)
```

---

## 2. Table Definitions

### 2.1 users Table

| Column | SQLite Type | Go Type | Constraints | Default | Description |
|--------|-------------|---------|-------------|---------|-------------|
| `id` | TEXT | `string` | PRIMARY KEY | - | Format: `usr_<8-char-uuid>` |
| `name` | TEXT | `string` | NOT NULL | - | Full display name (1-100 chars) |
| `username` | TEXT | `string` | NOT NULL, UNIQUE | - | Login username (3-50 chars, alphanumeric) |
| `email` | TEXT | `string` | - | NULL | Email for notifications (optional) |
| `password` | TEXT | `string` | NOT NULL | - | bcrypt hash (cost 10) |
| `role` | TEXT | `domain.Role` | NOT NULL, CHECK | - | `'admin'` or `'employee'` |
| `vacation_days` | INTEGER | `int` | - | 0 | Total annual days (employees only) |
| `used_vacation_days` | INTEGER | `int` | - | 0 | Days already used (employees only) |
| `email_preferences` | TEXT | `json.RawMessage` | NOT NULL | `'{}'` | JSON object (see section 3) |
| `created_at` | DATETIME | `time.Time` | NOT NULL | CURRENT_TIMESTAMP | Account creation |
| `updated_at` | DATETIME | `time.Time` | NOT NULL | CURRENT_TIMESTAMP | Last modification (via trigger) |

**ID Format:**
```
usr_a1b2c3d4
│   └─────── First 8 characters of UUID v4
└─────────── Prefix for entity type
```

**Role Values:**
- `admin` - Full system access (Captain)
- `employee` - Limited access (Crew Member)

### 2.2 vacation_requests Table

| Column | SQLite Type | Go Type | Constraints | Default | Description |
|--------|-------------|---------|-------------|---------|-------------|
| `id` | TEXT | `string` | PRIMARY KEY | - | Format: `vac_<8-char-uuid>` |
| `user_id` | TEXT | `string` | NOT NULL, FK | - | References users.id |
| `start_date` | TEXT | `string` | NOT NULL | - | Format: `YYYY-MM-DD` |
| `end_date` | TEXT | `string` | NOT NULL | - | Format: `YYYY-MM-DD` |
| `business_days` | INTEGER | `int` | NOT NULL | - | Calculated working days |
| `status` | TEXT | `domain.VacationStatus` | NOT NULL, CHECK | - | `'pending'`, `'approved'`, `'rejected'` |
| `reason` | TEXT | `string` | - | NULL | Optional notes (max 200 chars) |
| `created_at` | DATETIME | `time.Time` | NOT NULL | CURRENT_TIMESTAMP | Request submission time |
| `reviewed_by` | TEXT | `*string` | - | NULL | Admin user ID who reviewed |
| `reviewed_at` | DATETIME | `*time.Time` | - | NULL | Review timestamp |

**ID Format:**
```
vac_x9y8z7w6
│   └─────── First 8 characters of UUID v4
└─────────── Prefix for entity type
```

**Status Values:**
- `pending` - Awaiting admin review
- `approved` - Approved by admin
- `rejected` - Rejected by admin

**Date Format:**
```
2024-07-15
│    │  └── Day (01-31)
│    └───── Month (01-12)
└────────── Year (4 digits)
```

### 2.3 settings Table

| Column | SQLite Type | Go Type | Constraints | Default | Description |
|--------|-------------|---------|-------------|---------|-------------|
| `id` | INTEGER | - | PRIMARY KEY, CHECK(id=1) | - | Enforces single row |
| `weekend_policy` | TEXT | `json.RawMessage` | NOT NULL | (see default) | Weekend handling config |
| `newsletter` | TEXT | `json.RawMessage` | NOT NULL | (see default) | Newsletter scheduling config |

**Singleton Pattern:** The `CHECK(id = 1)` constraint ensures only one row can exist.

---

## 3. JSON Field Schemas

### 3.1 email_preferences (users table)

**For Employee Users:**
```json
{
  "enabled": true,
  "vacationStatusUpdates": true,
  "vacationResetNotifications": true,
  "monthlyVacationSummary": true
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `true` | Master toggle for all notifications |
| `vacationStatusUpdates` | boolean | `true` | Notify when request approved/rejected |
| `vacationResetNotifications` | boolean | `true` | Notify on year-end reset |
| `monthlyVacationSummary` | boolean | `true` | Receive monthly newsletter |

**For Admin Users:**
```json
{
  "enabled": true,
  "vacationRequestNotifications": true,
  "userCreatedNotifications": true,
  "vacationResetNotifications": true,
  "monthlyVacationSummary": true
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `true` | Master toggle for all notifications |
| `vacationRequestNotifications` | boolean | `true` | Notify when employees submit requests |
| `userCreatedNotifications` | boolean | `true` | Confirm new user creation |
| `vacationResetNotifications` | boolean | `true` | Confirm year-end reset |
| `monthlyVacationSummary` | boolean | `true` | Receive monthly newsletter |

### 3.2 weekend_policy (settings table)

```json
{
  "excludeWeekends": true
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `excludeWeekends` | boolean | `true` | If true, weekends don't count as vacation days |

**Impact on Business Day Calculation:**
- `true`: Mon-Fri only (5 days/week)
- `false`: All 7 days count

### 3.3 newsletter (settings table)

```json
{
  "enabled": false,
  "dayOfMonth": 1,
  "hourOfDay": 9,
  "lastSentAt": null
}
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | boolean | `false` | Enable automatic monthly newsletters |
| `dayOfMonth` | integer | `1` | Day to send (1-31) |
| `hourOfDay` | integer | `9` | Hour to send in 24h format (0-23) |
| `lastSentAt` | string/null | `null` | ISO 8601 timestamp of last send |

---

## 4. Indexes

### 4.1 Index Definitions

| Index Name | Table | Column(s) | Purpose |
|------------|-------|-----------|---------|
| `idx_vacation_user` | vacation_requests | `user_id` | Fast lookup of user's vacation requests |
| `idx_vacation_status` | vacation_requests | `status` | Fast filtering by status (pending/approved/rejected) |
| `idx_users_role` | users | `role` | Fast filtering by role (admin/employee) |

### 4.2 Query Performance

**Optimized Queries:**
```sql
-- Fast: Uses idx_vacation_user
SELECT * FROM vacation_requests WHERE user_id = ?;

-- Fast: Uses idx_vacation_status
SELECT * FROM vacation_requests WHERE status = 'pending';

-- Fast: Uses idx_users_role
SELECT * FROM users WHERE role = 'admin';
```

---

## 5. Triggers

### 5.1 update_user_timestamp

**Purpose:** Automatically update `updated_at` when a user record is modified.

```sql
CREATE TRIGGER IF NOT EXISTS update_user_timestamp
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
```

**Behavior:**
- Fires after any UPDATE on users table
- Sets `updated_at` to current timestamp
- Uses SQLite's CURRENT_TIMESTAMP

---

## 6. Type Mappings

### 6.1 SQLite to Go Type Mapping

| SQLite Type | Go Type | Nullable Go Type | Notes |
|-------------|---------|------------------|-------|
| `TEXT` (id) | `string` | - | UUID-based: `usr_xxxxxxxx` or `vac_xxxxxxxx` |
| `TEXT` (date) | `string` | - | ISO format: `YYYY-MM-DD` |
| `TEXT` (enum) | Custom type | - | `domain.Role`, `domain.VacationStatus` |
| `TEXT` (json) | `struct` | - | Marshal/unmarshal via `encoding/json` |
| `INTEGER` | `int` | - | Standard Go integer |
| `DATETIME` | `time.Time` | `*time.Time` | Parsed by SQLite driver |
| `TEXT` (nullable) | - | `sql.NullString` | Use in scanning, convert to `string` |
| `DATETIME` (nullable) | - | `sql.NullTime` | Use in scanning, convert to `*time.Time` |

### 6.2 Go Domain Types

```go
// internal/domain/user.go
type Role string
const (
    RoleAdmin    Role = "admin"
    RoleEmployee Role = "employee"
)

// internal/domain/vacation.go
type VacationStatus string
const (
    StatusPending  VacationStatus = "pending"
    StatusApproved VacationStatus = "approved"
    StatusRejected VacationStatus = "rejected"
)
```

### 6.3 JSON Marshaling

**EmailPreferences Struct:**
```go
type EmailPreferences struct {
    Enabled                      bool `json:"enabled"`
    VacationStatusUpdates        bool `json:"vacationStatusUpdates,omitempty"`        // Employee
    VacationRequestNotifications bool `json:"vacationRequestNotifications,omitempty"` // Admin
    UserCreatedNotifications     bool `json:"userCreatedNotifications,omitempty"`     // Admin
    VacationResetNotifications   bool `json:"vacationResetNotifications"`
    MonthlyVacationSummary       bool `json:"monthlyVacationSummary"`
}
```

**Repository Pattern for JSON:**
```go
// Storing
prefs, _ := json.Marshal(user.EmailPreferences)
// Use string(prefs) in INSERT/UPDATE

// Reading
var prefsJSON string
row.Scan(&prefsJSON)
json.Unmarshal([]byte(prefsJSON), &user.EmailPreferences)
```

---

## 7. Complete Migration SQL

### 7.1 migrations/001_init.sql

```sql
-- ============================================================================
-- VacayTracker Database Migration: 001_init
-- ============================================================================
-- Description: Initial database schema
-- Created: 2024-12-09
-- ============================================================================

-- Enable WAL mode for better concurrency and performance
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA foreign_keys=ON;

-- ============================================================================
-- TABLE: users
-- ============================================================================
-- Stores user accounts with role-based access control and notification preferences
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id                    TEXT PRIMARY KEY,
    name                  TEXT NOT NULL,
    username              TEXT UNIQUE NOT NULL,
    email                 TEXT,
    password              TEXT NOT NULL,
    role                  TEXT NOT NULL CHECK (role IN ('admin', 'employee')),
    vacation_days         INTEGER DEFAULT 0,
    used_vacation_days    INTEGER DEFAULT 0,
    email_preferences     TEXT NOT NULL DEFAULT '{}',
    created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- TABLE: vacation_requests
-- ============================================================================
-- Stores vacation requests with status tracking and review information
-- ============================================================================

CREATE TABLE IF NOT EXISTS vacation_requests (
    id                    TEXT PRIMARY KEY,
    user_id               TEXT NOT NULL,
    start_date            TEXT NOT NULL,
    end_date              TEXT NOT NULL,
    business_days         INTEGER NOT NULL,
    status                TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    reason                TEXT,
    created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_by           TEXT,
    reviewed_at           DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- TABLE: settings
-- ============================================================================
-- Single-row table for system-wide configuration
-- ============================================================================

CREATE TABLE IF NOT EXISTS settings (
    id                    INTEGER PRIMARY KEY CHECK (id = 1),
    weekend_policy        TEXT NOT NULL DEFAULT '{"excludeWeekends": true}',
    newsletter            TEXT NOT NULL DEFAULT '{"enabled": false, "dayOfMonth": 1, "hourOfDay": 9}'
);

-- Insert default settings row
INSERT OR IGNORE INTO settings (id) VALUES (1);

-- ============================================================================
-- INDEXES
-- ============================================================================

CREATE INDEX IF NOT EXISTS idx_vacation_user ON vacation_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_vacation_status ON vacation_requests(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- ============================================================================
-- TRIGGERS
-- ============================================================================

CREATE TRIGGER IF NOT EXISTS update_user_timestamp
AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- ============================================================================
-- END OF MIGRATION 001_init
-- ============================================================================
```

### 7.2 Running Migrations (Go Code)

```go
// internal/repository/sqlite/sqlite.go

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
```

---

## 8. Seed Data

### 8.1 Default Admin User

```sql
-- Create default admin user (password from ADMIN_PASSWORD env var, bcrypt hashed)
INSERT INTO users (id, name, username, password, role, email_preferences, created_at, updated_at)
VALUES (
    'usr_admin001',
    'Admin',
    'admin',
    '$2a$10$...', -- bcrypt hash of ADMIN_PASSWORD
    'admin',
    '{"enabled":true,"vacationRequestNotifications":true,"userCreatedNotifications":true,"vacationResetNotifications":true,"monthlyVacationSummary":true}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

### 8.2 Test Users (Development Only)

```sql
-- Test employee: john / john123
INSERT INTO users (id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at)
VALUES (
    'usr_john0001',
    'John Doe',
    'john',
    'john@example.com',
    '$2a$10$...', -- bcrypt hash of 'john123'
    'employee',
    25,
    0,
    '{"enabled":true,"vacationStatusUpdates":true,"vacationResetNotifications":true,"monthlyVacationSummary":true}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Test employee: jane / jane123
INSERT INTO users (id, name, username, email, password, role, vacation_days, used_vacation_days, email_preferences, created_at, updated_at)
VALUES (
    'usr_jane0001',
    'Jane Smith',
    'jane',
    'jane@example.com',
    '$2a$10$...', -- bcrypt hash of 'jane123'
    'employee',
    25,
    0,
    '{"enabled":true,"vacationStatusUpdates":true,"vacationResetNotifications":true,"monthlyVacationSummary":true}',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);
```

### 8.3 Seeding in Go (main.go)

```go
func seedDefaultAdmin(ctx context.Context, userRepo repository.UserRepository, authService *service.AuthService, adminPassword string) error {
    // Check if admin exists
    existing, err := userRepo.GetByUsername(ctx, "admin")
    if err != nil {
        return err
    }
    if existing != nil {
        return nil // Admin already exists, skip
    }

    // Hash password
    hashedPassword, err := authService.HashPassword(adminPassword)
    if err != nil {
        return err
    }

    // Create admin user
    admin := &domain.User{
        ID:               "usr_admin001",
        Name:             "Admin",
        Username:         "admin",
        Password:         hashedPassword,
        Role:             domain.RoleAdmin,
        EmailPreferences: domain.DefaultAdminEmailPreferences(),
        CreatedAt:        time.Now(),
        UpdatedAt:        time.Now(),
    }

    return userRepo.Create(ctx, admin)
}
```

---

## Implementation Checklist

- [ ] Create `migrations/` directory
- [ ] Create `migrations/001_init.sql` with complete schema
- [ ] Implement `sqlite.New()` with WAL mode configuration
- [ ] Implement `sqlite.Migrate()` to run SQL files
- [ ] Create domain types (`Role`, `VacationStatus`, structs)
- [ ] Implement repository scanning with NULL handling
- [ ] Add seed logic for default admin in `main.go`

---

*Document Version: 1.0*
*Generated for AI Agent Implementation*
