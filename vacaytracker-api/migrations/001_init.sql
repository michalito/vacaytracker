-- ============================================
-- VacayTracker Initial Database Schema
-- Migration: 001_init
-- ============================================

-- Users table
-- Stores employee and admin accounts
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
-- Stores all vacation requests with their status
CREATE TABLE IF NOT EXISTS vacation_requests (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    total_days INTEGER NOT NULL,
    reason TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by TEXT,
    reviewed_at TEXT,
    rejection_reason TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Settings table (singleton)
-- Stores application-wide settings
CREATE TABLE IF NOT EXISTS settings (
    id TEXT PRIMARY KEY DEFAULT 'settings',
    weekend_policy TEXT NOT NULL DEFAULT '{"excludeWeekends":true,"excludedDays":[0,6]}',
    newsletter TEXT NOT NULL DEFAULT '{"enabled":false,"frequency":"monthly","dayOfMonth":1}',
    default_vacation_days INTEGER NOT NULL DEFAULT 25,
    vacation_reset_month INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- ============================================
-- Indexes for query performance
-- ============================================

-- Index for querying vacation requests by user
CREATE INDEX IF NOT EXISTS idx_vacation_requests_user_id ON vacation_requests(user_id);

-- Index for querying vacation requests by status (pending, approved, rejected)
CREATE INDEX IF NOT EXISTS idx_vacation_requests_status ON vacation_requests(status);

-- Index for querying users by role (admin, employee)
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Index for user email lookups (login)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ============================================
-- Triggers for automatic timestamp updates
-- ============================================

-- Update users.updated_at on modification
CREATE TRIGGER IF NOT EXISTS users_updated_at
    AFTER UPDATE ON users
    FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Update vacation_requests.updated_at on modification
CREATE TRIGGER IF NOT EXISTS vacation_requests_updated_at
    AFTER UPDATE ON vacation_requests
    FOR EACH ROW
BEGIN
    UPDATE vacation_requests SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Update settings.updated_at on modification
CREATE TRIGGER IF NOT EXISTS settings_updated_at
    AFTER UPDATE ON settings
    FOR EACH ROW
BEGIN
    UPDATE settings SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- ============================================
-- Default data
-- ============================================

-- Insert default settings row (singleton pattern)
INSERT OR IGNORE INTO settings (id) VALUES ('settings');
