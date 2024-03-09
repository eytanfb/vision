package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func BuildSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int) string {
	titleStyle := summaryTitleStyle(width)
	progressTextStyle := startedTextStyle

	textStyle := startedTextStyle
	view := ""
	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]
		isInactive := m.TaskManager.TaskCollection.IsInactive(category)
		if !m.ViewManager.IsWeeklyView && isInactive {
			continue
		}
		completedTasks, totalTasks := m.TaskManager.TaskCollection.Progress(category)
		percentage := float64(completedTasks) / float64(totalTasks)
		if !m.ViewManager.IsWeeklyView && percentage == 1 {
			continue
		}
		roundedUpPercentage := int(percentage*10) * 10

		taskTitle := category[0 : len(category)-len(".md")]
		progressText := progressBar(completedTasks, totalTasks) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
		tasksView := ""
		incompleteTaskCount := 0
		for _, task := range tasks {
			if task.Completed && !task.IsCompletedToday() {
				progressText = strings.Replace(progressText, " 🛫", "", -1)
				continue
			}
			incompleteTaskCount++
			tasks := ""
			text := task.Summary()
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

			tasks = textStyle.Render(text)
			tasksView = lipgloss.JoinVertical(lipgloss.Top, tasksView, tasks)
		}

		rightAlignedProgressText := progressTextStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
		taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
		taskTitleView := lipgloss.JoinHorizontal(lipgloss.Left, titleStyle.Render(taskTitle), rightAlignedProgressText)
		tasksView = lipgloss.JoinVertical(lipgloss.Top, taskTitleContainer(width).Render(taskTitleView), tasksView)
		view = lipgloss.JoinVertical(lipgloss.Top, view, tasksView)
	}

	return view
}
