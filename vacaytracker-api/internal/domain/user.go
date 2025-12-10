package domain

import (
	"encoding/json"
	"time"
)

// Role represents a user's role in the system
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

// EmailPreferences holds user notification settings
type EmailPreferences struct {
	VacationUpdates   bool `json:"vacationUpdates"`
	WeeklyDigest      bool `json:"weeklyDigest"`
	TeamNotifications bool `json:"teamNotifications"`
}

// User represents an employee or admin in the system
type User struct {
	ID               string           `json:"id"`
	Email            string           `json:"email"`
	PasswordHash     string           `json:"-"` // Never expose password hash
	Name             string           `json:"name"`
	Role             Role             `json:"role"`
	VacationBalance  int              `json:"vacationBalance"`
	StartDate        *string          `json:"startDate,omitempty"`
	EmailPreferences EmailPreferences `json:"emailPreferences"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// IsAdmin returns true if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEmployee returns true if the user has employee role
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

// DefaultEmailPreferences returns default email notification settings
func DefaultEmailPreferences() EmailPreferences {
	return EmailPreferences{
		VacationUpdates:   true,
		WeeklyDigest:      false,
		TeamNotifications: true,
	}
}

// ParseEmailPreferences parses JSON string into EmailPreferences struct
func ParseEmailPreferences(data string) (EmailPreferences, error) {
	if data == "" {
		return DefaultEmailPreferences(), nil
	}

	var prefs EmailPreferences
	if err := json.Unmarshal([]byte(data), &prefs); err != nil {
		return DefaultEmailPreferences(), err
	}
	return prefs, nil
}

// ToJSONString converts EmailPreferences to JSON string for database storage
func (e EmailPreferences) ToJSONString() (string, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ValidRoles returns all valid role values
func ValidRoles() []Role {
	return []Role{RoleAdmin, RoleEmployee}
}

// IsValidRole checks if a role string is valid
func IsValidRole(role string) bool {
	for _, r := range ValidRoles() {
		if string(r) == role {
			return true
		}
	}
	return false
}
