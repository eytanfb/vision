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
	overdue
	completed
	completed_past
)

type Task struct {
	Company       string
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
	return t.Text
}

func (t Task) HumanizedString() string {
	var humanizedString string

	if t.ScheduledDate != "" {
		scheduledDate, _ := time.Parse("2006-01-02", t.ScheduledDate)

		today := time.Now()

		if scheduledDate.Day() == today.Day() && scheduledDate.Month() == today.Month() && scheduledDate.Year() == today.Year() {
			humanizedString += "Scheduled for today"
		} else {
			formattedScheduledDate := scheduledDate.Format("Jan 2")

			humanizedString += "Scheduled for " + formattedScheduledDate
		}
	}

	if t.StartDate != "" {
		startedDate, _ := time.Parse("2006-01-02", t.StartDate)
		formattedStartedDate := startedDate.Format("Jan 2")

		if humanizedString != "" {
			humanizedString += ", "
		}

		if t.ScheduledDate != "" {
			humanizedString += "and started on " + formattedStartedDate
		} else {
			humanizedString += "Started on " + formattedStartedDate
		}
	}

	if t.CompletedDate != "" {
		parsedCompletedDate, _ := time.Parse("2006-01-02", t.CompletedDate)
		formattedCompletedDate := parsedCompletedDate.Format("Jan 2")

		if t.StartDate != "" {
			parsedStartDate, _ := time.Parse("2006-01-02", t.StartDate)
			formattedStartDate := parsedStartDate.Format("Jan 2")

			days := parsedCompletedDate.Sub(parsedStartDate).Hours() / 24
			if days == 0 {
				humanizedString = fmt.Sprintf("Started on %s, completed the same day", formattedStartDate)
			} else if days == 1 {
				humanizedString = fmt.Sprintf("Started on %s, took 1 day, completed on %s", formattedStartDate, formattedCompletedDate)
			} else {
				humanizedString = fmt.Sprintf("Started on %s, took %.0f days, completed on %s", formattedStartDate, days, formattedCompletedDate)
			}
		} else {
			humanizedString = "Completed on " + formattedCompletedDate
		}
	}

	return humanizedString
}

func (t Task) Summary() string {
	return t.textWithoutDates()
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

	return parsedCompletedDate.Day() == time.Now().Day() && parsedCompletedDate.Month() == time.Now().Month() && parsedCompletedDate.Year() == time.Now().Year()
}

func (t Task) IsScheduledForDay(date string) bool {
	if t.ScheduledDate == "" {
		return false
	}

	parsedScheduledDate, err := time.Parse("2006-01-02", t.ScheduledDate)
	if err != nil {
		return false
	}

	dateDay, _ := time.Parse("2006-01-02", date)

	return parsedScheduledDate.Day() == dateDay.Day()
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
	status := unscheduled

	if t.Completed {
		if date == t.CompletedDate {
			return completed
		} else if date > t.CompletedDate {
			return completed_past
		}
	}

	if status != completed && t.Started && date >= t.StartDate {
		if t.CompletedDate == "" || date < t.CompletedDate {
			status = started
		}
	}

	if status != completed && status != started && t.Scheduled && date >= t.ScheduledDate {
		if (t.StartDate == "" || date < t.StartDate) && (t.CompletedDate == "" || date < t.CompletedDate) {
			status = scheduled
		}
	}

	if status == scheduled {
		parsedDate, _ := time.Parse("2006-01-02", date)
		parsedScheduleDate, _ := time.Parse("2006-01-02", t.ScheduledDate)

		scheduledDays := parsedDate.Sub(parsedScheduleDate).Hours() / 24
		if scheduledDays > 14 {
			status = overdue
		}
	} else if status == started {
		parsedDate, _ := time.Parse("2006-01-02", date)
		parsedStartDate, _ := time.Parse("2006-01-02", t.StartDate)

		startedDays := parsedDate.Sub(parsedStartDate).Hours() / 24
		if startedDays > 14 {
			status = overdue
		}
	}

	return status
}

func (t Task) WeeklyStatusAtDate(date string) status {
	status := unscheduled

	if t.CompletedDate != "" && date >= t.CompletedDate {
		status = completed
	}

	if t.StartDate != "" && date >= t.StartDate {
		if t.CompletedDate == "" || date < t.CompletedDate {
			status = started
		}
	}

	if t.ScheduledDate != "" && date >= t.ScheduledDate {
		if (t.StartDate == "" || date < t.StartDate) && (t.CompletedDate == "" || date < t.CompletedDate) {
			status = scheduled
		}
	}

	if status == scheduled {
		parsedDate, _ := time.Parse("2006-01-02", date)
		parsedScheduleDate, _ := time.Parse("2006-01-02", t.ScheduledDate)

		scheduledDays := parsedDate.Sub(parsedScheduleDate).Hours() / 24
		if scheduledDays > 14 {
			status = overdue
		}
	} else if status == started {
		parsedDate, _ := time.Parse("2006-01-02", date)
		parsedStartDate, _ := time.Parse("2006-01-02", t.StartDate)

		startedDays := parsedDate.Sub(parsedStartDate).Hours() / 24
		if startedDays > 14 {
			status = overdue
		}
	}

	return status
}

func (t Task) LastUpdatedAt() string {
	if t.CompletedDate != "" {
		return t.CompletedDate
	}

	if t.StartDate != "" {
		return t.StartDate
	}

	if t.ScheduledDate != "" {
		return t.ScheduledDate
	}

	return ""
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
