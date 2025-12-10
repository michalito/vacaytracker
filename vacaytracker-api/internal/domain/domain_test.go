package domain

import (
	"testing"
)

// ============================================
// User Tests
// ============================================

func TestUserIsAdmin(t *testing.T) {
	admin := &User{Role: RoleAdmin}
	if !admin.IsAdmin() {
		t.Error("IsAdmin() should return true for admin role")
	}
	if admin.IsEmployee() {
		t.Error("IsEmployee() should return false for admin role")
	}

	employee := &User{Role: RoleEmployee}
	if employee.IsAdmin() {
		t.Error("IsAdmin() should return false for employee role")
	}
	if !employee.IsEmployee() {
		t.Error("IsEmployee() should return true for employee role")
	}
}

func TestDefaultEmailPreferences(t *testing.T) {
	prefs := DefaultEmailPreferences()

	if !prefs.VacationUpdates {
		t.Error("VacationUpdates should be true by default")
	}
	if prefs.WeeklyDigest {
		t.Error("WeeklyDigest should be false by default")
	}
	if !prefs.TeamNotifications {
		t.Error("TeamNotifications should be true by default")
	}
}

func TestParseEmailPreferences(t *testing.T) {
	// Test valid JSON
	json := `{"vacationUpdates":false,"weeklyDigest":true,"teamNotifications":false}`
	prefs, err := ParseEmailPreferences(json)
	if err != nil {
		t.Errorf("ParseEmailPreferences() error = %v", err)
	}
	if prefs.VacationUpdates {
		t.Error("VacationUpdates should be false")
	}
	if !prefs.WeeklyDigest {
		t.Error("WeeklyDigest should be true")
	}

	// Test empty string returns defaults
	prefs, err = ParseEmailPreferences("")
	if err != nil {
		t.Errorf("ParseEmailPreferences() error = %v", err)
	}
	if !prefs.VacationUpdates {
		t.Error("Should return defaults for empty string")
	}
}

func TestEmailPreferencesToJSONString(t *testing.T) {
	prefs := EmailPreferences{
		VacationUpdates:   true,
		WeeklyDigest:      false,
		TeamNotifications: true,
	}

	json, err := prefs.ToJSONString()
	if err != nil {
		t.Errorf("ToJSONString() error = %v", err)
	}
	if json == "" {
		t.Error("ToJSONString() should not return empty string")
	}
}

func TestIsValidRole(t *testing.T) {
	if !IsValidRole("admin") {
		t.Error("'admin' should be a valid role")
	}
	if !IsValidRole("employee") {
		t.Error("'employee' should be a valid role")
	}
	if IsValidRole("invalid") {
		t.Error("'invalid' should not be a valid role")
	}
}

// ============================================
// Vacation Tests
// ============================================

func TestVacationRequestStatus(t *testing.T) {
	pending := &VacationRequest{Status: StatusPending}
	if !pending.IsPending() {
		t.Error("IsPending() should return true for pending status")
	}
	if pending.IsApproved() {
		t.Error("IsApproved() should return false for pending status")
	}
	if pending.IsRejected() {
		t.Error("IsRejected() should return false for pending status")
	}

	approved := &VacationRequest{Status: StatusApproved}
	if approved.IsPending() {
		t.Error("IsPending() should return false for approved status")
	}
	if !approved.IsApproved() {
		t.Error("IsApproved() should return true for approved status")
	}

	rejected := &VacationRequest{Status: StatusRejected}
	if !rejected.IsRejected() {
		t.Error("IsRejected() should return true for rejected status")
	}
}

func TestVacationCanBeCancelled(t *testing.T) {
	pending := &VacationRequest{Status: StatusPending}
	if !pending.CanBeCancelled() {
		t.Error("Pending requests should be cancellable")
	}

	approved := &VacationRequest{Status: StatusApproved}
	if approved.CanBeCancelled() {
		t.Error("Approved requests should not be cancellable")
	}

	rejected := &VacationRequest{Status: StatusRejected}
	if rejected.CanBeCancelled() {
		t.Error("Rejected requests should not be cancellable")
	}
}

func TestIsValidStatus(t *testing.T) {
	if !IsValidStatus("pending") {
		t.Error("'pending' should be a valid status")
	}
	if !IsValidStatus("approved") {
		t.Error("'approved' should be a valid status")
	}
	if !IsValidStatus("rejected") {
		t.Error("'rejected' should be a valid status")
	}
	if IsValidStatus("invalid") {
		t.Error("'invalid' should not be a valid status")
	}
}

// ============================================
// Settings Tests
// ============================================

func TestDefaultWeekendPolicy(t *testing.T) {
	policy := DefaultWeekendPolicy()

	if !policy.ExcludeWeekends {
		t.Error("ExcludeWeekends should be true by default")
	}
	if len(policy.ExcludedDays) != 2 {
		t.Error("ExcludedDays should have 2 days by default")
	}
	// Check Sunday (0) and Saturday (6) are excluded
	hasZero := false
	hasSix := false
	for _, day := range policy.ExcludedDays {
		if day == 0 {
			hasZero = true
		}
		if day == 6 {
			hasSix = true
		}
	}
	if !hasZero || !hasSix {
		t.Error("ExcludedDays should include 0 (Sunday) and 6 (Saturday)")
	}
}

func TestDefaultNewsletterConfig(t *testing.T) {
	config := DefaultNewsletterConfig()

	if config.Enabled {
		t.Error("Newsletter should be disabled by default")
	}
	if config.Frequency != "monthly" {
		t.Errorf("Frequency should be 'monthly', got '%s'", config.Frequency)
	}
	if config.DayOfMonth != 1 {
		t.Errorf("DayOfMonth should be 1, got %d", config.DayOfMonth)
	}
}

func TestDefaultSettings(t *testing.T) {
	settings := DefaultSettings()

	if settings.ID != "settings" {
		t.Errorf("ID should be 'settings', got '%s'", settings.ID)
	}
	if settings.DefaultVacationDays != 25 {
		t.Errorf("DefaultVacationDays should be 25, got %d", settings.DefaultVacationDays)
	}
	if settings.VacationResetMonth != 1 {
		t.Errorf("VacationResetMonth should be 1, got %d", settings.VacationResetMonth)
	}
}

func TestWeekendPolicyIsDayExcluded(t *testing.T) {
	policy := DefaultWeekendPolicy()

	// Sunday (0) should be excluded
	if !policy.IsDayExcluded(0) {
		t.Error("Sunday (0) should be excluded")
	}

	// Saturday (6) should be excluded
	if !policy.IsDayExcluded(6) {
		t.Error("Saturday (6) should be excluded")
	}

	// Monday (1) should NOT be excluded
	if policy.IsDayExcluded(1) {
		t.Error("Monday (1) should not be excluded")
	}

	// Test with ExcludeWeekends = false
	policy.ExcludeWeekends = false
	if policy.IsDayExcluded(0) {
		t.Error("No day should be excluded when ExcludeWeekends is false")
	}
}

func TestParseWeekendPolicy(t *testing.T) {
	// Test valid JSON
	json := `{"excludeWeekends":false,"excludedDays":[1,2]}`
	policy, err := ParseWeekendPolicy(json)
	if err != nil {
		t.Errorf("ParseWeekendPolicy() error = %v", err)
	}
	if policy.ExcludeWeekends {
		t.Error("ExcludeWeekends should be false")
	}

	// Test empty string returns defaults
	policy, err = ParseWeekendPolicy("")
	if err != nil {
		t.Errorf("ParseWeekendPolicy() error = %v", err)
	}
	if !policy.ExcludeWeekends {
		t.Error("Should return defaults for empty string")
	}
}

func TestWeekendPolicyToJSONString(t *testing.T) {
	policy := WeekendPolicy{
		ExcludeWeekends: true,
		ExcludedDays:    []int{0, 6},
	}

	json, err := policy.ToJSONString()
	if err != nil {
		t.Errorf("ToJSONString() error = %v", err)
	}
	if json == "" {
		t.Error("ToJSONString() should not return empty string")
	}
}
