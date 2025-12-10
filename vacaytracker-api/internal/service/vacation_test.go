package service

import (
	"testing"
	"time"

	"vacaytracker-api/internal/domain"
)

func TestParseDDMMYYYY(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid date with leading zeros",
			input:   "25/12/2025",
			want:    time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date without leading zeros",
			input:   "1/1/2025",
			want:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date single digit day",
			input:   "5/12/2025",
			want:    time.Date(2025, 12, 5, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid date single digit month",
			input:   "15/3/2025",
			want:    time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format - wrong separator",
			input:   "25-12-2025",
			wantErr: true,
		},
		{
			name:    "invalid format - too few parts",
			input:   "25/12",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong year length",
			input:   "25/12/25",
			wantErr: true,
		},
		{
			name:    "invalid date - day too long",
			input:   "125/12/2025",
			wantErr: true,
		},
		{
			name:    "invalid date - month too long",
			input:   "25/123/2025",
			wantErr: true,
		},
		{
			name:    "invalid date - February 30",
			input:   "30/02/2025",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDDMMYYYY(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDDMMYYYY(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("parseDDMMYYYY(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCalculateBusinessDays(t *testing.T) {
	// Helper to create dates
	date := func(year, month, day int) time.Time {
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}

	tests := []struct {
		name   string
		start  time.Time
		end    time.Time
		policy domain.WeekendPolicy
		want   int
	}{
		{
			name:  "single weekday - Monday",
			start: date(2025, 12, 1), // Monday
			end:   date(2025, 12, 1), // Monday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6}, // Sunday, Saturday
			},
			want: 1,
		},
		{
			name:  "single weekend day - Saturday",
			start: date(2025, 12, 6), // Saturday
			end:   date(2025, 12, 6), // Saturday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			want: 0,
		},
		{
			name:  "full week Mon-Sun excluding weekends",
			start: date(2025, 12, 1), // Monday
			end:   date(2025, 12, 7), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			want: 5, // Mon-Fri
		},
		{
			name:  "full week Mon-Sun including weekends",
			start: date(2025, 12, 1), // Monday
			end:   date(2025, 12, 7), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: false,
				ExcludedDays:    []int{},
			},
			want: 7, // All days
		},
		{
			name:  "two weeks excluding weekends",
			start: date(2025, 12, 1),  // Monday
			end:   date(2025, 12, 14), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			want: 10, // 2 weeks x 5 weekdays
		},
		{
			name:  "Friday to Monday excluding weekends",
			start: date(2025, 12, 5),  // Friday
			end:   date(2025, 12, 8),  // Monday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			want: 2, // Friday + Monday
		},
		{
			name:  "only Saturday and Sunday",
			start: date(2025, 12, 6), // Saturday
			end:   date(2025, 12, 7), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{0, 6},
			},
			want: 0,
		},
		{
			name:  "custom excluded day - Friday only",
			start: date(2025, 12, 1), // Monday
			end:   date(2025, 12, 7), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{5}, // Friday only
			},
			want: 6, // Mon, Tue, Wed, Thu, Sat, Sun
		},
		{
			name:  "Middle East weekend - Friday and Saturday",
			start: date(2025, 12, 1), // Monday
			end:   date(2025, 12, 7), // Sunday
			policy: domain.WeekendPolicy{
				ExcludeWeekends: true,
				ExcludedDays:    []int{5, 6}, // Friday, Saturday
			},
			want: 5, // Sun, Mon, Tue, Wed, Thu
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateBusinessDays(tt.start, tt.end, tt.policy)
			if got != tt.want {
				t.Errorf("calculateBusinessDays() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalculateBusinessDays_EdgeCases(t *testing.T) {
	date := func(year, month, day int) time.Time {
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}

	standardPolicy := domain.WeekendPolicy{
		ExcludeWeekends: true,
		ExcludedDays:    []int{0, 6},
	}

	// Test a month with 31 days
	t.Run("full month January 2025", func(t *testing.T) {
		// January 2025: Wed Jan 1 to Fri Jan 31
		// Has 5 Saturdays (4, 11, 18, 25) and 5 Sundays (5, 12, 19, 26)
		// Wait - Jan 2025: 1st is Wednesday
		// Weekends: 4-5, 11-12, 18-19, 25-26 = 8 weekend days
		// 31 - 8 = 23 business days
		got := calculateBusinessDays(date(2025, 1, 1), date(2025, 1, 31), standardPolicy)
		if got != 23 {
			t.Errorf("January 2025 business days = %d, want 23", got)
		}
	})

	// Test leap year February
	t.Run("February 2024 leap year", func(t *testing.T) {
		// Feb 2024: Thu Feb 1 to Thu Feb 29 (leap year)
		// Weekends: 3-4, 10-11, 17-18, 24-25 = 8 weekend days
		// 29 - 8 = 21 business days
		got := calculateBusinessDays(date(2024, 2, 1), date(2024, 2, 29), standardPolicy)
		if got != 21 {
			t.Errorf("February 2024 business days = %d, want 21", got)
		}
	})

	// Test year crossing
	t.Run("year crossing Dec 30 to Jan 2", func(t *testing.T) {
		// Dec 30, 2025 is Tuesday, Dec 31 is Wednesday
		// Jan 1, 2026 is Thursday, Jan 2 is Friday
		// All weekdays = 4
		got := calculateBusinessDays(date(2025, 12, 30), date(2026, 1, 2), standardPolicy)
		if got != 4 {
			t.Errorf("Year crossing business days = %d, want 4", got)
		}
	})
}
