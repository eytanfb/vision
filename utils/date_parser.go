package utils

import (
	"regexp"
	"strings"
	"time"
)

// ParseHashtagsToObsidianDates converts hashtag time references to scheduled date format (⏳ YYYY-MM-DD)
// Supported hashtags: #today, #tomorrow, #nextmonday, #nexttuesday, #nextwednesday, #nextthursday, #nextfriday, #nextweek
func ParseHashtagsToObsidianDates(input string) string {
	return ParseHashtagsToObsidianDatesWithTime(input, time.Now())
}

// ParseHashtagsToObsidianDatesWithTime allows injecting a specific time for testing
func ParseHashtagsToObsidianDatesWithTime(input string, currentTime time.Time) string {
	// Define hashtag patterns (case insensitive)
	hashtagPatterns := map[string]func(time.Time) string{
		"#today":         func(t time.Time) string { return formatScheduledDate(t) },
		"#tomorrow":      func(t time.Time) string { return formatScheduledDate(t.AddDate(0, 0, 1)) },
		"#nextweek":      func(t time.Time) string { return formatScheduledDate(t.AddDate(0, 0, 7)) },
		"#nextmonday":    func(t time.Time) string { return formatScheduledDate(calculateNextWeekday(t, "monday")) },
		"#nexttuesday":   func(t time.Time) string { return formatScheduledDate(calculateNextWeekday(t, "tuesday")) },
		"#nextwednesday": func(t time.Time) string { return formatScheduledDate(calculateNextWeekday(t, "wednesday")) },
		"#nextthursday":  func(t time.Time) string { return formatScheduledDate(calculateNextWeekday(t, "thursday")) },
		"#nextfriday":    func(t time.Time) string { return formatScheduledDate(calculateNextWeekday(t, "friday")) },
	}

	result := input

	// Create case-insensitive regex for each hashtag
	for hashtag, dateFunc := range hashtagPatterns {
		// Create regex pattern that matches the hashtag case-insensitively
		pattern := "(?i)" + regexp.QuoteMeta(hashtag) + "\\b"
		regex := regexp.MustCompile(pattern)

		// Replace all matches with the calculated scheduled date
		result = regex.ReplaceAllStringFunc(result, func(match string) string {
			return dateFunc(currentTime)
		})
	}

	return result
}

// calculateNextWeekday calculates the next occurrence of the specified weekday
// If today is the target weekday, it returns the same weekday next week (+7 days)
func calculateNextWeekday(currentTime time.Time, targetWeekday string) time.Time {
	weekdayMap := map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}

	targetDay, exists := weekdayMap[strings.ToLower(targetWeekday)]
	if !exists {
		// Return current time if invalid weekday
		return currentTime
	}

	currentWeekday := currentTime.Weekday()

	// Calculate days until target weekday
	daysUntilTarget := int(targetDay - currentWeekday)

	// If target day is today or in the past this week, go to next week
	if daysUntilTarget <= 0 {
		daysUntilTarget += 7
	}

	return currentTime.AddDate(0, 0, daysUntilTarget)
}

// formatScheduledDate formats a time.Time to the app's scheduled date format ⏳ YYYY-MM-DD
func formatScheduledDate(t time.Time) string {
	return "⏳ " + t.Format("2006-01-02")
}
