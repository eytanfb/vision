package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func BuildSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int, date string) string {
	titleStyle := summaryTitleStyle(width)
	progressTextStyle := titleStyle

	view := ""
	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]
		progressText := buildProgressText(m, category)
		taskTitle := category[0 : len(category)-len(".md")]
		tasksView := ""
		incompleteTaskCount := len(m.TaskManager.TaskCollection.IncompleteTasks(category))

		for _, task := range tasks {
			textStyle := startedTextStyle

			if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
				continue
			}

			tasks := ""
			text := ""
			status := task.StatusAtDate(date)

			if m.ViewManager.IsWeeklyView {
				status = task.WeeklyStatusAtDate(date)
			}

			text, textStyle = buildTaskView(task, progressText, status)

			tasks = textStyle.Render(text)
			tasksView = joinVertical(tasksView, tasks)
		}

		rightAlignedProgressText := progressTextStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
		taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
		taskTitleView := joinHorizontal(titleStyle.Render(taskTitle), rightAlignedProgressText)
		tasksView = joinVertical(taskTitleContainer(width).Render(taskTitleView), tasksView)
		view = joinVertical(view, tasksView)
	}

	return view
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

func buildTaskView(task Task, progressText string, status status) (string, lipgloss.Style) {
	var textStyle lipgloss.Style

	text := task.Summary()

	if status == completed {
		text += " ✅ " + daysAgoFromString(task.CompletedDate)
		textStyle = completedTextStyle
	} else if status == started {
		text += " 🛫 " + daysAgoFromString(task.StartDate)
		if !strings.Contains(progressText, "🛫") {
			progressText = addIconToProgressText(progressText, "🛫")
		}
	} else if status == scheduled {
		text += " ⏳ " + daysAgoFromString(task.ScheduledDate)
		textStyle = scheduledTextStyle
		if !strings.Contains(progressText, "⏳") && !strings.Contains(progressText, "🚨") && !strings.Contains(progressText, "🛫") {
			progressText = addIconToProgressText(progressText, "⏳")
		}
	}
	if task.IsOverdue() {
		text += " 🚨"
		textStyle = overdueTextStyle
		if !strings.Contains(progressText, "🚨") {
			progressText = addIconToProgressText(progressText, "🚨")
		}
	}

	return text, textStyle
}

func progressBar(completed, total int) string {
	percentage := float64(completed) / float64(total)
	numberOfBars := int(percentage * 10)
	progressBar := "["

	for i := 0; i < 10; i++ {
		if i < numberOfBars {
			progressBar += "#"
		} else {
			progressBar += " "
		}
	}
	progressBar += "]"

	return progressBar
}

func buildProgressText(m *Model, category string) string {
	completedTasksCount, totalTasksCount := m.TaskManager.TaskCollection.Progress(category)
	percentage := float64(completedTasksCount) / float64(totalTasksCount)
	roundedUpPercentage := int(percentage*10) * 10
	return progressBar(completedTasksCount, totalTasksCount) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
}

func summaryTitleStyle(width int) lipgloss.Style {
	summaryStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(width - 40).Bold(true).Foreground(lipgloss.Color("#9A9CCD"))

	return summaryStyle
}
