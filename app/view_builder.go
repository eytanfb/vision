package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func BuildSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int, date string) string {
	titleStyle := summaryTitleStyle(width)
	progressTextStyle := titleStyle

	view := ""
	for _, key := range keys {
		var progressText string
		category := key
		tasks := tasksByFile[key]
		taskTitle := category[0 : len(category)-len(".md")]
		tasksView := ""
		incompleteTaskCount := 0
		for _, task := range tasks {
			textStyle := startedTextStyle

			if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
				incompleteTaskCount++
				continue
			}

			incompleteTaskCount++
			tasks := ""
			text := ""
			status := task.StatusAtDate(date)
			if m.ViewManager.IsWeeklyView {
				status = task.WeeklyStatusAtDate(date)
			}
			if status == completed {
				incompleteTaskCount--
			}
			tasks = buildTaskView(task, date, status, progressText)
			progressText = buildProgressText(m, category, status)

			tasks = textStyle.Render(text)
			tasksView = joinVertical(tasksView, tasks)
		}

		rightAlignedProgressText := progressTextStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
		taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
		taskTitleView := lipgloss.JoinHorizontal(lipgloss.Left, titleStyle.Render(taskTitle), rightAlignedProgressText)
		tasksView = joinVertical(taskTitleContainer(width).Render(taskTitleView), tasksView)
		view = joinVertical(view, tasksView)
	}

	return view
}

func buildTaskView(task Task, date string, status status, progressText string) string {
	var textStyle lipgloss.Style

	text := task.Summary(date)

	if status == completed {
		textStyle = completedTextStyle
	} else if status == started {
		if !strings.Contains(progressText, "🛫") {
			progressText = addIconToProgressText(progressText, "🛫")
		}
	} else if status == scheduled {
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

	return textStyle.Render(text)
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

func buildProgressText(m *Model, category string, status status) string {
	completedTasksCount, totalTasksCount := m.TaskManager.TaskCollection.Progress(category)
	percentage := float64(completedTasksCount) / float64(totalTasksCount)
	roundedUpPercentage := int(percentage*10) * 10
	progressText := progressBar(completedTasksCount, totalTasksCount) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
	return progressText
}
