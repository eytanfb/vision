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

func BuildKanbanSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int, date string) string {
	view := ""

	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]

		view = buildTaskFileView(m, category, width, date, view, tasks)
	}

	return view
}

func kanbanListView(tasks []Task, date string, boardNumber int) string {
	view := ""

	return view
}

func BuildTasksForFileView(m *Model, tasks []Task, date string, cursor int) string {
	view := ""

	for index, task := range tasks {
		view = buildTaskForFileView(m, task, date, view, cursor, index)
	}

	return view
}

func BuildFilesView(m *Model, hiddenSidebar bool) (string, string) {
	list := ""
	itemDetails := ""
	completedList := ""
	activeList := ""
	inactiveList := ""

	for index, file := range m.FileManager.Files {
		style := defaultTextStyle

		if index == m.FileManager.FilesCursor {
			style = highlightedTextStyle
			itemDetails = file.Content
			m.FileManager.SelectedFile = file
		}

		line := file.FileNameWithoutExtension()

		if m.DirectoryManager.SelectedCategory == "tasks" {
			if index > 9 {
				break
			}
			list, activeList, completedList, inactiveList = buildTaskFilesView(m, line, index, file, style, activeList, completedList, inactiveList)
		} else {
			list = joinVertical(list, style.Render(line))
		}
	}

	if hiddenSidebar {
		list = ""
	}

	if m.IsAddTaskView() {
		itemDetails = m.NewTaskInput.View()
	} else {
		markdown := renderMarkdown(itemDetails)
		m.Viewport.SetContent(markdown)
		itemDetails = m.Viewport.View()
	}

	return list, itemDetails
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
	progressText := buildProgressText(m, category)

	return progressTextStyle(titleStyle).Render(progressText)
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

func buildTaskForFileView(m *Model, task Task, date string, view string, cursor int, index int) string {
	width := m.ViewManager.DetailsViewWidth - 15

	if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
		return ""
	}

	tasksString := TaskView{
		task:   task,
		date:   date,
		weekly: true,
		width:  width,
	}.RenderedText()

	tasksString = taskStyle(width).Render(tasksString)

	if index == cursor {
		datesString := dateStyle().Render(task.HumanizedString())

		tasksString = joinVertical(tasksString, "\n", datesContainerStyle(width).Render(datesString))
		tasksString = tasksStringStyle(width).Render(tasksString)
	}

	view = joinVertical(view, tasksString)

	return view
}

func buildTaskFilesView(m *Model, line string, index int, file FileInfo, style lipgloss.Style, activeList string, completedList string, inactiveList string) (string, string, string, string) {
	isInactive := m.TaskManager.TaskCollection.IsInactive(file.Name)
	completed, total := m.TaskManager.TaskCollection.Progress(file.Name)

	if total > 0 && completed == total {
		if index != m.FileManager.FilesCursor {
			style = completedFileStyle
		}
		completedList = joinVertical(completedList, style.Render(line))
	} else if isInactive {
		if index != m.FileManager.FilesCursor {
			style = inactiveFileStyle
		}
		inactiveList = joinVertical(inactiveList, style.Render(line))
	} else {
		text := fmt.Sprintf("%d/%d", completed, total)
		line = "[" + text + "] " + line
		activeList = joinVertical(activeList, style.Render(line))
	}

	activeTitle := taskFileTitleStyle.Render("Active")
	inactiveTitle := inactiveTitleStyle().Render("Inactive")
	completeTitle := completedTitleStyle().Render("Complete")

	renderedActiveList := renderedActiveListStyle().Render(activeList)
	renderedInactiveList := renderedInactiveListStyle().Render(inactiveList)
	renderedCompletedList := renderedCompletedListStyle().Render(completedList)

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
