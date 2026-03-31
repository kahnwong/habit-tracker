package calendar

import (
	"testing"
	"time"
)

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{2024, true},  // divisible by 4
		{2023, false}, // not divisible by 4
		{2000, true},  // divisible by 400
		{1900, false}, // divisible by 100 but not 400
		{2100, false}, // divisible by 100 but not 400
		{2004, true},  // divisible by 4
		{2001, false}, // not divisible by 4
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.year)), func(t *testing.T) {
			result := isLeapYear(tt.year)
			if result != tt.expected {
				t.Errorf("isLeapYear(%d) = %v, expected %v", tt.year, result, tt.expected)
			}
		})
	}
}

func TestGetDaysInMonth(t *testing.T) {
	tests := []struct {
		year     int
		month    time.Month
		expected int
	}{
		{2024, time.January, 31},
		{2024, time.February, 29}, // leap year
		{2023, time.February, 28}, // non-leap year
		{2024, time.March, 31},
		{2024, time.April, 30},
		{2024, time.May, 31},
		{2024, time.June, 30},
		{2024, time.July, 31},
		{2024, time.August, 31},
		{2024, time.September, 30},
		{2024, time.October, 31},
		{2024, time.November, 30},
		{2024, time.December, 31},
		{2000, time.February, 29}, // leap year (divisible by 400)
		{1900, time.February, 28}, // non-leap year (divisible by 100 but not 400)
	}

	for _, tt := range tests {
		t.Run(tt.month.String(), func(t *testing.T) {
			result := getDaysInMonth(tt.year, tt.month)
			if result != tt.expected {
				t.Errorf("getDaysInMonth(%d, %s) = %d, expected %d", tt.year, tt.month, result, tt.expected)
			}
		})
	}
}

func TestIsHighlighted(t *testing.T) {
	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	highlightDates := []time.Time{date1, date2}

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"exact match", date1, true},
		{"another match", date2, true},
		{"no match - different month", date3, false},
		{"no match - different day", time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC), false},
		{"match with different time", time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHighlighted(tt.date, highlightDates)
			if result != tt.expected {
				t.Errorf("isHighlighted(%v) = %v, expected %v", tt.date, result, tt.expected)
			}
		})
	}
}

func TestGetMaxWeeks(t *testing.T) {
	tests := []struct {
		name     string
		months   []time.Time
		expected int
	}{
		{
			name:     "single month starting on Sunday",
			months:   []time.Time{time.Date(2024, 9, 1, 0, 0, 0, 0, time.Local)}, // Sept 2024 starts on Sunday
			expected: 5,
		},
		{
			name:     "single month starting on Friday",
			months:   []time.Time{time.Date(2024, 3, 1, 0, 0, 0, 0, time.Local)}, // March 2024 starts on Friday
			expected: 6,
		},
		{
			name: "multiple months",
			months: []time.Time{
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMaxWeeks(tt.months)
			if result != tt.expected {
				t.Errorf("getMaxWeeks() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

func TestGenerateMonths(t *testing.T) {
	// Save original now and restore after test
	originalNow := now
	defer func() { now = originalNow }()

	tests := []struct {
		name           string
		currentDate    time.Time
		lookbackMonths int
		expectedMonths []time.Month
	}{
		{
			name:           "3 months from March 31 (edge case)",
			currentDate:    time.Date(2026, 3, 31, 0, 0, 0, 0, time.Local),
			lookbackMonths: 3,
			expectedMonths: []time.Month{time.January, time.February, time.March},
		},
		{
			name:           "3 months from January 31",
			currentDate:    time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
			lookbackMonths: 3,
			expectedMonths: []time.Month{time.November, time.December, time.January},
		},
		{
			name:           "2 months from May 15",
			currentDate:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.Local),
			lookbackMonths: 2,
			expectedMonths: []time.Month{time.April, time.May},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now = tt.currentDate
			result := generateMonths(tt.lookbackMonths)

			if len(result) != tt.lookbackMonths {
				t.Errorf("generateMonths() returned %d months, expected %d", len(result), tt.lookbackMonths)
			}

			for i, month := range result {
				if month.Month() != tt.expectedMonths[i] {
					t.Errorf("generateMonths()[%d] = %s, expected %s", i, month.Month(), tt.expectedMonths[i])
				}
			}
		})
	}
}
