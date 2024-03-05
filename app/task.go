package app

import (
	"regexp"
	"strings"
	"time"
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

func (t Task) Summary() string {
	return t.textWithoutDates()
}

func (t Task) IsOverdue() bool {
	today := time.Now()
	isOverdue := false

	if t.Scheduled {
		parsedScheduleDate, err := time.Parse("2006-01-02", t.ScheduledDate)
		if err != nil {
			isOverdue = false
		}
		scheduledDays := today.Sub(parsedScheduleDate).Hours() / 24
		if scheduledDays > 14 {
			isOverdue = true
		}
	}

	if t.Started {
		isOverdue = false
		parsedStartDate, err := time.Parse("2006-01-02", t.StartDate)
		if err != nil {
			return false
		}

		startedDays := today.Sub(parsedStartDate).Hours() / 24
		if startedDays > 14 {
			isOverdue = true
		}
	}

	return isOverdue
}

func (t Task) IsInactive() bool {
	return (!t.Started && !t.Scheduled) || t.Completed
}

func (t Task) textWithoutDates() string {
	return removeDatesFromText(t.Text)
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
	var date string
	for i := index + len(icon); i < len(text); i++ {
		if text[i] == ' ' {
			if len(date) > 0 {
				break
			}
			continue
		}
		date += string(text[i])
	}

	return date
}

func removeDatesFromText(text string) string {
	datesRegex := regexp.MustCompile(`[✅, ⏳, 🛫]\s+\d{4}-\d{2}-\d{2}`)

	text = datesRegex.ReplaceAllString(text, "")

	return strings.Trim(text, " ")
}
