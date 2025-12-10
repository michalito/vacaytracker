package domain

import (
	"encoding/json"
	"time"
)

// WeekendPolicy defines which days are excluded from business day calculations
type WeekendPolicy struct {
	ExcludeWeekends bool  `json:"excludeWeekends"`
	ExcludedDays    []int `json:"excludedDays"` // 0 = Sunday, 6 = Saturday
}

// NewsletterConfig holds newsletter scheduling settings
type NewsletterConfig struct {
	Enabled    bool       `json:"enabled"`
	Frequency  string     `json:"frequency"`  // "weekly" or "monthly"
	DayOfMonth int        `json:"dayOfMonth"` // 1-28 for monthly frequency
	LastSentAt *time.Time `json:"lastSentAt"` // Track last newsletter send time
}

// Settings holds application-wide configuration stored in the database
type Settings struct {
	ID                  string           `json:"id"` // Always "settings" (singleton)
	WeekendPolicy       WeekendPolicy    `json:"weekendPolicy"`
	Newsletter          NewsletterConfig `json:"newsletter"`
	DefaultVacationDays int              `json:"defaultVacationDays"`
	VacationResetMonth  int              `json:"vacationResetMonth"` // 1-12 (January = 1)
	UpdatedAt           time.Time        `json:"updatedAt"`
}

// DefaultWeekendPolicy returns the default weekend policy
// By default, weekends (Saturday and Sunday) are excluded
func DefaultWeekendPolicy() WeekendPolicy {
	return WeekendPolicy{
		ExcludeWeekends: true,
		ExcludedDays:    []int{0, 6}, // Sunday = 0, Saturday = 6
	}
}

// DefaultNewsletterConfig returns the default newsletter settings
// By default, newsletter is disabled
func DefaultNewsletterConfig() NewsletterConfig {
	return NewsletterConfig{
		Enabled:    false,
		Frequency:  "monthly",
		DayOfMonth: 1,
		LastSentAt: nil,
	}
}

// DefaultSettings returns a Settings struct with default values
func DefaultSettings() Settings {
	return Settings{
		ID:                  "settings",
		WeekendPolicy:       DefaultWeekendPolicy(),
		Newsletter:          DefaultNewsletterConfig(),
		DefaultVacationDays: 25,
		VacationResetMonth:  1, // January
		UpdatedAt:           time.Now(),
	}
}

// ParseWeekendPolicy parses JSON string into WeekendPolicy struct
func ParseWeekendPolicy(data string) (WeekendPolicy, error) {
	if data == "" {
		return DefaultWeekendPolicy(), nil
	}

	var policy WeekendPolicy
	if err := json.Unmarshal([]byte(data), &policy); err != nil {
		return DefaultWeekendPolicy(), err
	}
	return policy, nil
}

// ToJSONString converts WeekendPolicy to JSON string for database storage
func (w WeekendPolicy) ToJSONString() (string, error) {
	bytes, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseNewsletterConfig parses JSON string into NewsletterConfig struct
func ParseNewsletterConfig(data string) (NewsletterConfig, error) {
	if data == "" {
		return DefaultNewsletterConfig(), nil
	}

	var config NewsletterConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return DefaultNewsletterConfig(), err
	}
	return config, nil
}

// ToJSONString converts NewsletterConfig to JSON string for database storage
func (n NewsletterConfig) ToJSONString() (string, error) {
	bytes, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// IsDayExcluded checks if a given weekday is excluded from business day calculations
// weekday: 0 = Sunday, 1 = Monday, ..., 6 = Saturday
func (w WeekendPolicy) IsDayExcluded(weekday int) bool {
	if !w.ExcludeWeekends {
		return false
	}
	for _, day := range w.ExcludedDays {
		if day == weekday {
			return true
		}
	}
	return false
}
