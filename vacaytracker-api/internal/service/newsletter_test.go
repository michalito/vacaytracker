package service

import (
	"testing"
	"time"

	"vacaytracker-api/internal/domain"
)

func TestShouldSendNewsletter(t *testing.T) {
	tests := []struct {
		name     string
		settings domain.Settings
		now      time.Time
		expected bool
	}{
		{
			name: "monthly - correct day - not sent today",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    true,
					Frequency:  "monthly",
					DayOfMonth: 15,
					LastSentAt: nil,
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "monthly - wrong day",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    true,
					Frequency:  "monthly",
					DayOfMonth: 15,
					LastSentAt: nil,
				},
			},
			now:      time.Date(2025, 12, 14, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "monthly - already sent today",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    true,
					Frequency:  "monthly",
					DayOfMonth: 15,
					LastSentAt: timePtr(time.Date(2025, 12, 15, 8, 0, 0, 0, time.UTC)),
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "monthly - sent yesterday, today is the correct day",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    true,
					Frequency:  "monthly",
					DayOfMonth: 15,
					LastSentAt: timePtr(time.Date(2025, 11, 15, 8, 0, 0, 0, time.UTC)),
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "disabled",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    false,
					Frequency:  "monthly",
					DayOfMonth: 15,
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name: "weekly - Monday",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:   true,
					Frequency: "weekly",
				},
			},
			// December 15, 2025 is a Monday
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name: "weekly - Tuesday",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:   true,
					Frequency: "weekly",
				},
			},
			now:      time.Date(2025, 12, 16, 9, 0, 0, 0, time.UTC), // Tuesday
			expected: false,
		},
		{
			name: "weekly - Monday but already sent today",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:    true,
					Frequency:  "weekly",
					LastSentAt: timePtr(time.Date(2025, 12, 15, 6, 0, 0, 0, time.UTC)),
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC), // Monday
			expected: false,
		},
		{
			name: "invalid frequency",
			settings: domain.Settings{
				Newsletter: domain.NewsletterConfig{
					Enabled:   true,
					Frequency: "daily", // invalid
				},
			},
			now:      time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a scheduler with nil dependencies (we're only testing shouldSendNewsletter logic)
			s := &Scheduler{}
			result := s.shouldSendNewsletterAt(&tt.settings, tt.now)
			if result != tt.expected {
				t.Errorf("shouldSendNewsletterAt() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		{
			name:     "same day same time",
			t1:       time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			name:     "same day different time",
			t1:       time.Date(2025, 12, 15, 6, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 12, 15, 23, 59, 59, 0, time.UTC),
			expected: true,
		},
		{
			name:     "different days",
			t1:       time.Date(2025, 12, 15, 23, 59, 59, 0, time.UTC),
			t2:       time.Date(2025, 12, 16, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "different months",
			t1:       time.Date(2025, 11, 15, 9, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			name:     "different years",
			t1:       time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC),
			t2:       time.Date(2025, 12, 15, 9, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSameDay(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("isSameDay() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
