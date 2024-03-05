package app

import (
	"regexp"
	"strings"
)

type Task struct {
	IsDone        bool
	Text          string
	StartDate     string
	CompletedDate string
	ScheduledDate string
	LineNumber    int
	Completed     bool
	Started       bool
	Scheduled     bool
	FileName      string
}

func (t Task) String() string {
	var stringBuilder strings.Builder

	if t.IsDone {
		stringBuilder.WriteString("- [x] ")
	} else {
		stringBuilder.WriteString("- [ ] ")
	}

	stringBuilder.WriteString(t.Text)
	result := stringBuilder.String()
	resultWithoutDates := removeDatesFromText(result)
	stringBuilder.Reset()
	stringBuilder.WriteString(resultWithoutDates)

	if t.StartDate != "" || t.CompletedDate != "" || t.ScheduledDate != "" {
		stringBuilder.WriteString("\n")
		if t.ScheduledDate != "" {
			stringBuilder.WriteString("Scheduled: " + strings.Trim(t.ScheduledDate, " ") + "\n")
		}
		if t.StartDate != "" {
			stringBuilder.WriteString("Start: " + strings.Trim(t.StartDate, " ") + "\n")
		}
		if t.CompletedDate != "" {
			stringBuilder.WriteString("Completed: " + strings.Trim(t.CompletedDate, " ") + "\n")
		}
	}

	return stringBuilder.String()
}

func (t Task) textWithoutDates() string {
	return removeDatesFromText(t.Text)
}

func (t Task) Summary() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString(t.textWithoutDates() + " (" + t.FileName + ")")
	return stringBuilder.String()
}

func extractStartDateFromText(text string) string {
	startIcon := "🛫 "
	return extractDateFromText(text, startIcon)
}

func extractScheduledDateFromText(text string) string {
	scheduledIcon := "⏳"
	return extractDateFromText(text, scheduledIcon)
}

func extractCompletedDateFromText(text string) string {
	completedIcon := "✅ "
	return extractDateFromText(text, completedIcon)
}

func extractDateFromText(text string, icon string) string {
	index := strings.Index(text, icon)
	if index == -1 {
		return ""
	}
	// read date from the next 10 characters
	date := text[index : index+14]
	return date
}

func removeDatesFromText(text string) string {
	datesRegex := regexp.MustCompile(`[✅, ⏳, 🛫]\s+\d{4}-\d{2}-\d{2}`)

	text = datesRegex.ReplaceAllString(text, "")

	return strings.Trim(text, " ")
}
