package app

import (
	"fmt"

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
		incompleteTaskCount := len(m.TaskManager.TaskCollection.IncompleteTasks(category, date))

		for _, task := range tasks {
			if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
				continue
			}

			tasksString := TaskView{
				task:   task,
				date:   date,
				weekly: m.ViewManager.IsWeeklyView,
			}.RenderedText()

			tasksView = joinVertical(tasksView, tasksString)
		}

		rightAlignedProgressText := progressTextStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
		taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
		taskTitleView := joinHorizontal(titleStyle.Render(taskTitle), rightAlignedProgressText)
		tasksView = joinVertical(taskTitleContainer(width).Render(taskTitleView), tasksView)
		view = joinVertical(view, tasksView)
	}

	return view
}

func BuildTasksForFileView(m *Model, tasks []Task, date string, cursor int) string {
	view := ""
	background := "#474747"

	for index, task := range tasks {
		if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
			continue
		}

		taskStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth - 40)
		datesContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth - 40)
		dateStyle := lipgloss.NewStyle().Background(lipgloss.Color(background)).Foreground(lipgloss.Color("#9A9CCD")).PaddingRight(2)

		tasksString := TaskView{
			task:   task,
			date:   date,
			weekly: true,
		}.RenderedText()

		tasksString = taskStyle.Render(tasksString)

		if index == cursor {
			datesString := dateStyle.Render(task.HumanizedString())

			tasksString = joinVertical(tasksString, "\n", datesContainerStyle.Render(datesString))
			tasksString = lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(m.ViewManager.DetailsViewWidth - 40).PaddingTop(1).PaddingBottom(1).MarginTop(1).MarginBottom(1).Render(tasksString)
		}

		view = joinVertical(view, tasksString)
	}

	return view
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
