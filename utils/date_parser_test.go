package utils

import (
	"testing"
	"time"
)

func TestParseHashtagsToObsidianDates(t *testing.T) {
	// Mock current time for consistent testing
	now := time.Date(2025, 8, 11, 12, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name     string
		input    string
		expected string
		baseDate time.Time
	}{
		// Basic hashtag tests
		{
			name:     "today hashtag",
			input:    "Complete task #today",
			expected: "Complete task ⏳ 2025-08-11",
			baseDate: now,
		},
		{
			name:     "tomorrow hashtag",
			input:    "Review documents #tomorrow",
			expected: "Review documents ⏳ 2025-08-12",
			baseDate: now,
		},
		{
			name:     "nextweek hashtag",
			input:    "Weekly report #nextweek",
			expected: "Weekly report ⏳ 2025-08-18",
			baseDate: now,
		},

		// Next weekday tests (current day is Monday)
		{
			name:     "nextmonday (same as nextweek when today is Monday)",
			input:    "Meeting #nextmonday",
			expected: "Meeting ⏳ 2025-08-18",
			baseDate: now,
		},
		{
			name:     "nexttuesday (this week)",
			input:    "Call client #nexttuesday",
			expected: "Call client ⏳ 2025-08-12",
			baseDate: now,
		},
		{
			name:     "nextwednesday (this week)",
			input:    "Team standup #nextwednesday",
			expected: "Team standup ⏳ 2025-08-13",
			baseDate: now,
		},
		{
			name:     "nextthursday (this week)",
			input:    "Project review #nextthursday",
			expected: "Project review ⏳ 2025-08-14",
			baseDate: now,
		},
		{
			name:     "nextfriday (this week)",
			input:    "Deploy code #nextfriday",
			expected: "Deploy code ⏳ 2025-08-15",
			baseDate: now,
		},

		// Case insensitive tests
		{
			name:     "uppercase hashtag",
			input:    "Task #TODAY",
			expected: "Task ⏳ 2025-08-11",
			baseDate: now,
		},
		{
			name:     "mixed case hashtag",
			input:    "Task #Tomorrow",
			expected: "Task ⏳ 2025-08-12",
			baseDate: now,
		},

		// Multiple hashtags
		{
			name:     "multiple hashtags",
			input:    "Start #today, finish #tomorrow",
			expected: "Start ⏳ 2025-08-11, finish ⏳ 2025-08-12",
			baseDate: now,
		},

		// No hashtags
		{
			name:     "no hashtags",
			input:    "Regular task without dates",
			expected: "Regular task without dates",
			baseDate: now,
		},

		// Invalid hashtags (should remain unchanged)
		{
			name:     "invalid hashtag",
			input:    "Task #invalidday",
			expected: "Task #invalidday",
			baseDate: now,
		},

		// Edge case: hashtag at beginning
		{
			name:     "hashtag at beginning",
			input:    "#today complete the task",
			expected: "⏳ 2025-08-11 complete the task",
			baseDate: now,
		},

		// Edge case: hashtag at end
		{
			name:     "hashtag at end",
			input:    "Complete the task #today",
			expected: "Complete the task ⏳ 2025-08-11",
			baseDate: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseHashtagsToObsidianDatesWithTime(tt.input, tt.baseDate)
			if result != tt.expected {
				t.Errorf("ParseHashtagsToObsidianDates() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateNextWeekday(t *testing.T) {
	monday := time.Date(2025, 8, 11, 12, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name        string
		currentDay  time.Time
		targetDay   string
		expectedDay time.Time
	}{
		{
			name:        "Monday to Tuesday",
			currentDay:  monday,
			targetDay:   "tuesday",
			expectedDay: time.Date(2025, 8, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			name:        "Monday to Friday",
			currentDay:  monday,
			targetDay:   "friday",
			expectedDay: time.Date(2025, 8, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:        "Monday to next Monday",
			currentDay:  monday,
			targetDay:   "monday",
			expectedDay: time.Date(2025, 8, 18, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNextWeekday(tt.currentDay, tt.targetDay)
			if !result.Equal(tt.expectedDay) {
				t.Errorf("calculateNextWeekday() = %v, expected %v", result, tt.expectedDay)
			}
		})
	}
}

func TestParseHashtagsToObsidianDatesFromDifferentDays(t *testing.T) {
	// Test from Friday to ensure weekend handling
	friday := time.Date(2025, 8, 15, 12, 0, 0, 0, time.UTC) // Friday

	tests := []struct {
		name     string
		input    string
		expected string
		baseDate time.Time
	}{
		{
			name:     "Friday to next Monday",
			input:    "Meeting #nextmonday",
			expected: "Meeting ⏳ 2025-08-18",
			baseDate: friday,
		},
		{
			name:     "Friday to next Tuesday",
			input:    "Task #nexttuesday",
			expected: "Task ⏳ 2025-08-19",
			baseDate: friday,
		},
		{
			name:     "Friday to next Friday",
			input:    "Weekly #nextfriday",
			expected: "Weekly ⏳ 2025-08-22",
			baseDate: friday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseHashtagsToObsidianDatesWithTime(tt.input, tt.baseDate)
			if result != tt.expected {
				t.Errorf("ParseHashtagsToObsidianDates() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
