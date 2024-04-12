package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type TaskView struct {
	task   Task
	date   string
	weekly bool
}

func (tv TaskView) RenderedText() string {
	status := tv.task.StatusAtDate(tv.date)

	if tv.weekly {
		status = tv.task.WeeklyStatusAtDate(tv.date)
	}

	text := tv.text(status)
	textStyle := tv.textStyle(status)

	return textStyle.Render(text)
}

func (tv TaskView) text(status status) string {
	text := tv.task.Summary()

	if status == completed {
		text += " ✅ " + tv.daysAgoFromString(tv.task.CompletedDate)
	} else if status == started {
		text += " 🛫 " + tv.daysAgoFromString(tv.task.StartDate)
	} else if status == scheduled {
		text += " ⏳ " + tv.daysAgoFromString(tv.task.ScheduledDate)
	} else if status == overdue {
		text += " 🚨 "
		date := tv.task.ScheduledDate

		if tv.task.StartDate != "" {
			date = tv.task.StartDate
		}

		text += tv.daysAgoFromString(date)
	}

	return text
}

func (tv TaskView) textStyle(status status) lipgloss.Style {
	var textStyle lipgloss.Style

	if status == completed {
		textStyle = completedTextStyle
	} else if status == started {
		textStyle = startedTextStyle
	} else if status == scheduled {
		textStyle = scheduledTextStyle
	} else if status == overdue {
		textStyle = overdueTextStyle
	}

	return textStyle
}

func (tv TaskView) daysAgoFromString(date string) string {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}

	parsedTVDate, err := time.Parse("2006-01-02", tv.date)
	if err != nil {
		return ""
	}

	days := parsedTVDate.Sub(parsedDate).Hours() / 24
	daysString := "days"
	if days < 1 {
		return "today"
	} else if days < 2 {
		daysString = "day"
	}

	return fmt.Sprintf("%.0f %s ago", days, daysString)
}
