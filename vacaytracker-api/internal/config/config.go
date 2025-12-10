package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port   string
	Env    string
	AppURL string

	// Database
	DBPath string

	// Authentication
	JWTSecret     string
	AdminPassword string
	AdminEmail    string
	AdminName     string

	// Email (Resend)
	ResendAPIKey     string
	EmailFromAddress string
	EmailFromName    string
}

// Load reads configuration from environment variables
// It loads .env file if present (for development)
func Load() *Config {
	// Load .env file if it exists (ignore error for production)
	_ = godotenv.Load()

	cfg := &Config{
		// Server defaults
		Port:   getEnv("PORT", "3000"),
		Env:    getEnv("ENV", "development"),
		AppURL: getEnv("APP_URL", "http://localhost:3000"),

		// Database defaults
		DBPath: getEnv("DB_PATH", "./data/vacaytracker.db"),

		// Authentication (required)
		JWTSecret:     mustGetEnv("JWT_SECRET"),
		AdminPassword: mustGetEnv("ADMIN_PASSWORD"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@company.com"),
		AdminName:     getEnv("ADMIN_NAME", "Admin"),

		// Email (optional)
		ResendAPIKey:     getEnv("RESEND_API_KEY", ""),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", ""),
		EmailFromName:    getEnv("EMAIL_FROM_NAME", "VacayTracker"),
	}

	// Validate JWT secret length
	if len(cfg.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	return cfg
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// EmailEnabled returns true if email configuration is complete
func (c *Config) EmailEnabled() bool {
	return c.ResendAPIKey != "" && c.EmailFromAddress != ""
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an environment variable as int with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvBool retrieves an environment variable as bool with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// mustGetEnv retrieves a required environment variable
// It logs a fatal error if the variable is not set
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}
