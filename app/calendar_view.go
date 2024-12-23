package app

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

type CalendarView struct {
	startDate time.Time
	tasks     []Task
	width     int
	height    int
}

func NewCalendarView(tasks []Task, width, height int) CalendarView {
	// Start from 3 weeks ago + current week, beginning of week
	now := time.Now()
	startDate := now.AddDate(0, 0, -21) // Go back 3 weeks
	// Adjust to Monday of that week
	for startDate.Weekday() != time.Monday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	return CalendarView{
		startDate: startDate,
		tasks:     tasks,
		width:     width,
		height:    height,
	}
}

func (cv CalendarView) View() string {
	header := cv.renderHeader()
	grid := cv.renderGrid()
	return lipgloss.JoinVertical(lipgloss.Left, header, grid)
}

func (cv CalendarView) renderHeader() string {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri"}
	var headers []string
	cellWidth := (cv.width - 10) / 5 // Account for borders, 5 days only

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(summaryTitleColor).
		Width(cellWidth).
		Align(lipgloss.Center)

	for _, day := range days {
		headers = append(headers, headerStyle.Render(day))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, headers...)
}

func (cv CalendarView) renderGrid() string {
	var weeks []string
	currentDate := cv.startDate
	today := time.Now()

	cellWidth := (cv.width - 10) / 5
	cellStyle := lipgloss.NewStyle().
		Width(cellWidth).
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, true, false, true)

	for week := 0; week < 4; week++ {
		var days []string
		for day := 0; day < 5; day++ { // Only weekdays
			if !currentDate.After(today) { // Only show until today
				dayTasks := cv.getTasksForDate(currentDate)
				dayContent := cv.renderDay(currentDate, dayTasks)

				// Highlight current day
				if currentDate.Format("2006-01-02") == today.Format("2006-01-02") {
					cellStyle = cellStyle.BorderForeground(highlightedTextColor)
				} else {
					cellStyle = cellStyle.BorderForeground(lipgloss.Color("#9A9CCD"))
				}

				days = append(days, cellStyle.Render(dayContent))
			}
			currentDate = currentDate.AddDate(0, 0, 1)
		}
		// Skip weekend
		currentDate = currentDate.AddDate(0, 0, 2)
		weeks = append(weeks, lipgloss.JoinHorizontal(lipgloss.Top, days...))

		// Add spacing between weeks if not the last week
		if week < 3 {
			weeks = append(weeks, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, weeks...)
}

func (cv CalendarView) renderDay(date time.Time, tasks []Task) string {
	var content []string
	today := time.Now()

	// Style the date number with month
	dateStyle := lipgloss.NewStyle().
		Foreground(taskDateColor).
		Bold(true).
		MarginBottom(1)

	// If it's today, use green color and append (Today)
	if date.Format("2006-01-02") == today.Format("2006-01-02") {
		dateStyle = dateStyle.Foreground(completedColor)
		dateStr := date.Format("2 Jan") + " (Today)"
		content = append(content, dateStyle.Render(dateStr))
	} else {
		dateStr := date.Format("2 Jan")
		content = append(content, dateStyle.Render(dateStr))
	}

	for _, task := range tasks {
		status := task.StatusAtDate(date.Format("2006-01-02"))
		shouldShow := false
		var textStyle lipgloss.Style

		switch status {
		case completed:
			shouldShow = true
			textStyle = completedTextStyle
		case started:
			shouldShow = true
			textStyle = startedTextStyle
		case scheduled:
			shouldShow = true
			textStyle = scheduledTextStyle
		case overdue:
			shouldShow = true
			textStyle = overdueTextStyle
		default:
			textStyle = defaultTextStyle
		}

		if shouldShow {
			taskLine := truncateString(task.Summary(), 25)
			content = append(content, textStyle.Render(taskLine))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (cv CalendarView) getTasksForDate(date time.Time) []Task {
	dateStr := date.Format("2006-01-02")
	var dayTasks []Task

	for _, task := range cv.tasks {
		status := task.StatusAtDate(dateStr)
		if status != unscheduled {
			dayTasks = append(dayTasks, task)
		}
	}

	return dayTasks
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}
