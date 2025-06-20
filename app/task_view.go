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
	width  int
}

func (tv TaskView) RenderedText() string {
	status := tv.task.StatusAtDate(tv.date)
	priority := tv.task.Priority

	if tv.weekly {
		status = tv.task.WeeklyStatusAtDate(tv.date)
	}

	icon := tv.statusIcon(status)
	text := tv.task.Summary()
	textStyle := tv.textStyle(status)

	if priority != "" && status != completed {
		textStyle = priorityTextStyle
	}

	statusText := tv.statusText(status)

	renderedText := textStyle.Render(text)
	renderedStatusText := renderedStatusTextStyle.Render(statusText)

	return joinHorizontal(icon, renderedText, renderedStatusText)
}

func (tv TaskView) RenderedKanbanText() string {
	text := tv.RenderedText()

	return joinVertical(text)
}

func (tv TaskView) statusIcon(status status) string {
	icon := ""

	if status == completed {
		icon = "✅ "
	} else if status == started {
		icon = "🛫 "
	} else if status == scheduled {
		icon = "⏳ "
	} else if status == overdue {
		icon = "🚨 "
	}

	return iconStyle.Render(icon)
}

func (tv TaskView) statusText(status status) string {
	text := ""

	if status == completed {
		text += tv.daysAgoFromString(tv.task.CompletedDate)
	} else if status == started {
		text += tv.daysAgoFromString(tv.task.StartDate)
	} else if status == scheduled {
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

	if status == priority {
		textStyle = priorityTextStyle
	}

	return textStyle.Width(tv.width)
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
