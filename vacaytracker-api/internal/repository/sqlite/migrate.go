package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all pending database migrations
func (db *DB) RunMigrations(migrationsDir string) error {
	// Create schema_migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of already applied migrations
	applied, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get list of migration files
	files, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Execute pending migrations
	for _, file := range files {
		version := extractVersion(file)
		if applied[version] {
			continue // Already applied
		}

		if err := db.executeMigration(migrationsDir, file, version); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		fmt.Printf("Applied migration: %s\n", file)
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`
	_, err := db.Exec(query)
	return err
}

// getAppliedMigrations returns a map of already applied migration versions
func (db *DB) getAppliedMigrations() (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// getMigrationFiles returns a sorted list of .sql files in the migrations directory
func getMigrationFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	// Sort by filename (numeric prefix like 001_, 002_)
	sort.Strings(files)

	return files, nil
}

// extractVersion extracts the version from a migration filename
// e.g., "001_init.sql" -> "001"
func extractVersion(filename string) string {
	// Remove .sql extension
	name := strings.TrimSuffix(filename, ".sql")
	// Split by underscore and take first part
	parts := strings.SplitN(name, "_", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return name
}

// executeMigration executes a single migration file within a transaction
func (db *DB) executeMigration(dir, filename, version string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration within a transaction
	return db.Transaction(func(tx *sql.Tx) error {
		// Execute the migration SQL
		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration SQL: %w", err)
		}

		// Record the migration as applied
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version) VALUES (?)",
			version,
		); err != nil {
			return fmt.Errorf("failed to record migration: %w", err)
		}

		return nil
	})
}
