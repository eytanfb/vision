package app

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type status int

const (
	unscheduled status = iota
	scheduled
	started
	completed
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

const (
	StartedIcon   = "🛫 "
	CompletedIcon = "✅ "
	ScheduledIcon = "⏳"
)

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

func (t Task) Summary(date string) string {
	text := t.textWithoutDates()
	status := t.StatusAtDate(date)
	if status == completed {
		text += t.daysAgoWithIcon(t.CompletedDate, CompletedIcon)
	} else if status == started {
		text += t.daysAgoWithIcon(t.StartDate, StartedIcon)
	} else if status == scheduled {
		text += t.daysAgoWithIcon(t.ScheduledDate, ScheduledIcon)
	}
	return text
}

func (t Task) IsOverdue() bool {
	if t.Completed {
		return false
	}

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
	if t.ScheduledDate != "" {
		scheduledDate, err := time.Parse("2006-01-02", t.ScheduledDate)
		if err != nil {
			fmt.Println("Error parsing scheduled date", t.ScheduledDate)
			return false
		}
		if scheduledDate.After(time.Now()) {
			return true
		}
	}

	return (!t.Started && !t.Scheduled) || t.Completed
}

func (t Task) IsCompletedToday() bool {
	if t.CompletedDate == "" {
		return false
	}

	parsedCompletedDate, err := time.Parse("2006-01-02", t.CompletedDate)
	if err != nil {
		return false
	}

	return parsedCompletedDate.Day() == time.Now().Day()
}

func (t Task) IsScheduledForFuture(date string) bool {
	if t.ScheduledDate == "" {
		return false
	}

	return t.ScheduledDate > date
}

func (t Task) IsStarted() bool {
	return t.Started && !t.Completed
}

func (t Task) IsScheduled() bool {
	return t.Scheduled && !t.Completed && !t.Started
}

func (t Task) StatusAtDate(date string) status {
	if t.CompletedDate != "" && date == t.CompletedDate {
		return completed
	}

	if t.StartDate != "" && date >= t.StartDate {
		if t.CompletedDate == "" || date < t.CompletedDate {
			return started
		}
	}

	if t.ScheduledDate != "" && date >= t.ScheduledDate {
		if (t.StartDate == "" || date < t.StartDate) && (t.CompletedDate == "" || date < t.CompletedDate) {
			return scheduled
		}
	}

	return unscheduled
}

func (t Task) WeeklyStatusAtDate(date string) status {
	if t.CompletedDate != "" && date >= t.CompletedDate {
		return completed
	}

	if t.StartDate != "" && date >= t.StartDate {
		if t.CompletedDate == "" || date < t.CompletedDate {
			return started
		}
	}

	if t.ScheduledDate != "" && date >= t.ScheduledDate {
		if (t.StartDate == "" || date < t.StartDate) && (t.CompletedDate == "" || date < t.CompletedDate) {
			return scheduled
		}
	}

	return unscheduled
}

func (t Task) textWithoutDates() string {
	return removeDatesFromText(t.Text)
}

func extractStartDateFromText(text string) string {
	return extractDateFromText(text, StartedIcon)
}

func extractScheduledDateFromText(text string) string {
	return extractDateFromText(text, ScheduledIcon)
}

func extractCompletedDateFromText(text string) string {
	return extractDateFromText(text, CompletedIcon)
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

func daysAgoFromString(date string) string {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}

	today := time.Now()
	days := today.Sub(parsedDate).Hours() / 24
	daysString := "days"
	if days < 1 {
		return "today"
	} else if days < 2 {
		daysString = "day"
	}

	return fmt.Sprintf("%.0f %s ago", days, daysString)
}

func (t Task) daysAgoWithIcon(date, icon string) string {
	return " " + icon + " " + daysAgoFromString(date)
}
