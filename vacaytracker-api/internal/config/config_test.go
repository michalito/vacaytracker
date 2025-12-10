package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	if got := getEnv("TEST_VAR", "default"); got != "test_value" {
		t.Errorf("getEnv() = %v, want %v", got, "test_value")
	}

	// Test with non-existing env var
	if got := getEnv("NON_EXISTING_VAR", "default"); got != "default" {
		t.Errorf("getEnv() = %v, want %v", got, "default")
	}
}

func TestGetEnvInt(t *testing.T) {
	// Test with valid int
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	if got := getEnvInt("TEST_INT", 0); got != 42 {
		t.Errorf("getEnvInt() = %v, want %v", got, 42)
	}

	// Test with invalid int
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	if got := getEnvInt("TEST_INVALID_INT", 100); got != 100 {
		t.Errorf("getEnvInt() = %v, want %v", got, 100)
	}

	// Test with non-existing var
	if got := getEnvInt("NON_EXISTING_INT", 50); got != 50 {
		t.Errorf("getEnvInt() = %v, want %v", got, 50)
	}
}

func TestGetEnvBool(t *testing.T) {
	// Test with true value
	os.Setenv("TEST_BOOL_TRUE", "true")
	defer os.Unsetenv("TEST_BOOL_TRUE")

	if got := getEnvBool("TEST_BOOL_TRUE", false); got != true {
		t.Errorf("getEnvBool() = %v, want %v", got, true)
	}

	// Test with false value
	os.Setenv("TEST_BOOL_FALSE", "false")
	defer os.Unsetenv("TEST_BOOL_FALSE")

	if got := getEnvBool("TEST_BOOL_FALSE", true); got != false {
		t.Errorf("getEnvBool() = %v, want %v", got, false)
	}

	// Test with non-existing var
	if got := getEnvBool("NON_EXISTING_BOOL", true); got != true {
		t.Errorf("getEnvBool() = %v, want %v", got, true)
	}
}

func TestConfigHelperMethods(t *testing.T) {
	// Test IsDevelopment
	cfg := &Config{Env: "development"}
	if !cfg.IsDevelopment() {
		t.Error("IsDevelopment() should return true for development env")
	}
	if cfg.IsProduction() {
		t.Error("IsProduction() should return false for development env")
	}

	// Test IsProduction
	cfg = &Config{Env: "production"}
	if !cfg.IsProduction() {
		t.Error("IsProduction() should return true for production env")
	}
	if cfg.IsDevelopment() {
		t.Error("IsDevelopment() should return false for production env")
	}
}

func TestEmailEnabled(t *testing.T) {
	// Both set - should be enabled
	cfg := &Config{
		ResendAPIKey:     "re_test_key",
		EmailFromAddress: "test@example.com",
	}
	if !cfg.EmailEnabled() {
		t.Error("EmailEnabled() should return true when both API key and from address are set")
	}

	// Only API key set - should be disabled
	cfg = &Config{
		ResendAPIKey:     "re_test_key",
		EmailFromAddress: "",
	}
	if cfg.EmailEnabled() {
		t.Error("EmailEnabled() should return false when from address is empty")
	}

	// Only from address set - should be disabled
	cfg = &Config{
		ResendAPIKey:     "",
		EmailFromAddress: "test@example.com",
	}
	if cfg.EmailEnabled() {
		t.Error("EmailEnabled() should return false when API key is empty")
	}

	// Neither set - should be disabled
	cfg = &Config{
		ResendAPIKey:     "",
		EmailFromAddress: "",
	}
	if cfg.EmailEnabled() {
		t.Error("EmailEnabled() should return false when both are empty")
	}
}
