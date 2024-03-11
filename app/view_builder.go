package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func BuildSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int) string {
	titleStyle := summaryTitleStyle(width)
	progressTextStyle := startedTextStyle

	view := ""
	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]
		if isCategoryActive(m, category) {
			progressText := buildProgressText(m, category)
			taskTitle := category[0 : len(category)-len(".md")]
			tasksView := ""
			incompleteTaskCount := 0
			for _, task := range tasks {
				textStyle := startedTextStyle
				if task.Completed && !task.IsCompletedToday() {
					progressText = strings.Replace(progressText, " 🛫", "", -1)
					continue
				}
				incompleteTaskCount++
				tasks := ""
				text := ""
				text, textStyle, progressTextStyle, incompleteTaskCount = buildTaskView(task, progressText)

				tasks = textStyle.Render(text)
				tasksView = lipgloss.JoinVertical(lipgloss.Top, tasksView, tasks)
			}

			rightAlignedProgressText := progressTextStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
			taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
			taskTitleView := lipgloss.JoinHorizontal(lipgloss.Left, titleStyle.Render(taskTitle), rightAlignedProgressText)
			tasksView = lipgloss.JoinVertical(lipgloss.Top, taskTitleContainer(width).Render(taskTitleView), tasksView)
			view = lipgloss.JoinVertical(lipgloss.Top, view, tasksView)
		}
	}

	return view
}

func isCategoryActive(m *Model, category string) bool {
	if m.TaskManager.TaskCollection.IsInactive(category) {
		return false
	}

	completedTasks, totalTasks := m.TaskManager.TaskCollection.Progress(category)
	return float64(completedTasks)/float64(totalTasks) != 1
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

func buildTaskView(task Task, progressText string) (string, lipgloss.Style, lipgloss.Style, int) {
	var textStyle lipgloss.Style
	var progressTextStyle lipgloss.Style

	text := task.Summary()
	incompleteTaskCount := 0

	if task.Completed {
		incompleteTaskCount--
		text += " ✅ " + daysAgoFromString(task.CompletedDate)
		textStyle = completedTextStyle
		progressTextStyle = completedTextStyle
	} else if task.Started {
		text += " 🛫 " + daysAgoFromString(task.StartDate)
		if !strings.Contains(progressText, "🛫") {
			progressText = addIconToProgressText(progressText, "🛫")
			progressTextStyle = startedTextStyle
		}
	} else if task.Scheduled {
		text += " ⏳ " + daysAgoFromString(task.ScheduledDate)
		textStyle = scheduledTextStyle
		if !strings.Contains(progressText, "⏳") && !strings.Contains(progressText, "🚨") && !strings.Contains(progressText, "🛫") {
			progressText = addIconToProgressText(progressText, "⏳")
			progressTextStyle = scheduledTextStyle
		}
	}
	if task.IsOverdue() {
		text += " 🚨"
		textStyle = overdueTextStyle
		if !strings.Contains(progressText, "🚨") {
			progressText = addIconToProgressText(progressText, "🚨")
			progressTextStyle = overdueTextStyle
		}
	}

	return text, textStyle, progressTextStyle, incompleteTaskCount
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
