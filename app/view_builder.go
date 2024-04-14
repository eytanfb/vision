package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func BuildSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int, date string) string {
	view := ""
	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]
		view = buildTaskFileView(m, category, width, date, view, tasks)
	}

	return view
}

func buildTaskFileView(m *Model, category string, width int, date string, view string, tasks []Task) string {
	tasksView := ""

	for _, task := range tasks {
		tasksView = buildTaskView(m, task, date, tasksView)
	}

	tasksView = joinVertical(buildTaskTitleView(m, category, width, date), tasksView)
	view = joinVertical(view, tasksView)

	return view
}

func buildTaskTitleView(m *Model, category string, width int, date string) string {
	titleStyle := summaryTitleStyle(width)

	taskTitle := joinHorizontal(buildTaskTitle(m, category, date, titleStyle), buildRightAlignedProgressText(m, category, titleStyle))

	taskTitleView := taskTitleContainer(width).Render(taskTitle)

	return taskTitleView
}

func buildTaskTitle(m *Model, category string, date string, titleStyle lipgloss.Style) string {
	taskTitle := category[0 : len(category)-len(".md")]
	taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount(m, category, date)) + " tasks remaining)"
	return titleStyle.Render(taskTitle)
}

func incompleteTaskCount(m *Model, category string, date string) int {
	return len(m.TaskManager.TaskCollection.IncompleteTasks(category, date))
}

func buildRightAlignedProgressText(m *Model, category string, titleStyle lipgloss.Style) string {
	progressTextStyle := titleStyle

	progressText := buildProgressText(m, category)

	return progressTextStyle.Copy().Width(35).Align(lipgloss.Right).Render(progressText)
}

func buildTaskView(m *Model, task Task, date string, tasksView string) string {
	if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
		return ""
	}

	tasksString := TaskView{
		task:   task,
		date:   date,
		weekly: m.ViewManager.IsWeeklyView,
		width:  m.ViewManager.DetailsViewWidth - 25,
	}.RenderedText()

	return joinVertical(tasksView, tasksString)
}

func BuildTasksForFileView(m *Model, tasks []Task, date string, cursor int) string {
	view := ""

	for index, task := range tasks {
		view = buildTaskForFileView(m, task, date, view, cursor, index)
	}

	return view
}

func buildTaskForFileView(m *Model, task Task, date string, view string, cursor int, index int) string {
	background := "#474747"
	offset := 15

	if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
		return ""
	}

	taskStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth - offset)
	datesContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth - offset)
	dateStyle := lipgloss.NewStyle().Background(lipgloss.Color(background)).Foreground(lipgloss.Color("#9A9CCD")).PaddingRight(2)

	tasksString := TaskView{
		task:   task,
		date:   date,
		weekly: true,
		width:  m.ViewManager.DetailsViewWidth - offset,
	}.RenderedText()

	tasksString = taskStyle.Render(tasksString)

	if index == cursor {
		datesString := dateStyle.Render(task.HumanizedString())

		tasksString = joinVertical(tasksString, "\n", datesContainerStyle.Render(datesString))
		tasksString = lipgloss.NewStyle().Background(lipgloss.Color(background)).Width(m.ViewManager.DetailsViewWidth - offset).PaddingTop(1).PaddingBottom(1).MarginTop(1).MarginBottom(1).Render(tasksString)
	}

	view = joinVertical(view, tasksString)

	return view
}

func BuildTaskFilesView(m *Model, line string, index int, file FileInfo, style lipgloss.Style, activeList string, completedList string, inactiveList string) (string, string, string, string) {
	isInactive := m.TaskManager.TaskCollection.IsInactive(file.Name)

	completed, total := m.TaskManager.TaskCollection.Progress(file.Name)
	text := fmt.Sprintf("%d/%d", completed, total)
	var completedText string

	if total > 0 && completed == total {
		if index != m.FileManager.FilesCursor {
			style = style.Copy().Foreground(lipgloss.Color("#4DA165"))
		}
		completedList = joinVertical(completedList, style.Render(line))
	} else if isInactive {
		if index != m.FileManager.FilesCursor {
			style = style.Copy().Foreground(lipgloss.Color("#A0A0A0"))
		}
		inactiveList = joinVertical(inactiveList, style.Render(line))
	} else {
		completedText = lipgloss.NewStyle().Render(text)
		line = "[" + completedText + "] " + line
		activeList = joinVertical(activeList, style.Render(line))
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	activeTitle := titleStyle.Render("Active")
	inactiveTitle := titleStyle.Foreground(lipgloss.Color("#A0A0A0")).Render("Inactive")
	completeTitle := titleStyle.Foreground(lipgloss.Color("#4DA165")).Render("Complete")

	renderedActiveList := lipgloss.NewStyle().MarginTop(1).MarginBottom(2).Render(activeList)
	renderedInactiveList := lipgloss.NewStyle().MarginTop(1).MarginBottom(3).Render(inactiveList)
	renderedCompletedList := lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Render(completedList)

	return joinVertical(activeTitle, renderedActiveList, inactiveTitle, renderedInactiveList, completeTitle, renderedCompletedList), activeList, completedList, inactiveList
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
