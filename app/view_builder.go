package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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

type KanbanItem struct {
	filename string
	tasks    []Task
}

func BuildKanbanSummaryView(m *Model, keys []string, tasksByFile map[string][]Task, width int, date string) string {
	inactiveList, activeList, completedList := makeKanbanLists(m, keys, tasksByFile)
	setKanbanTasksCounts(inactiveList, activeList, completedList, m)

	activeBoard := renderBoard("Active", activeList, m, m.ViewManager.KanbanListCursor == 1)
	completedBoard := renderBoard("Completed", completedList, m, m.ViewManager.KanbanListCursor == 2)
	inactiveBoard := renderBoard("Inactive", inactiveList, m, m.ViewManager.KanbanListCursor == 0)

	return joinHorizontal(inactiveBoard, activeBoard, completedBoard)
}

func makeKanbanLists(m *Model, keys []string, tasksByFile map[string][]Task) ([]KanbanItem, []KanbanItem, []KanbanItem) {
	activeList := []KanbanItem{}
	completedList := []KanbanItem{}
	inactiveList := []KanbanItem{}

	for _, key := range keys {
		tasks := tasksByFile[key]

		for _, task := range tasks {
			if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) {
				inactiveList = addTaskOrCreateKanbanItem(inactiveList, key, task)
				continue
			}

			if task.Completed {
				completedList = addTaskOrCreateKanbanItem(completedList, key, task)
			} else if task.Started || task.Scheduled {
				activeList = addTaskOrCreateKanbanItem(activeList, key, task)
			} else if task.ScheduledDate == "" {
				inactiveList = addTaskOrCreateKanbanItem(inactiveList, key, task)
			}
		}
	}

	return inactiveList, activeList, completedList
}

func addTaskOrCreateKanbanItem(list []KanbanItem, filename string, task Task) []KanbanItem {
	for i, item := range list {
		if item.filename == filename {
			list[i].tasks = append(list[i].tasks, task)
			return list
		}
	}

	return append(list, KanbanItem{filename: filename, tasks: []Task{task}})
}

func setKanbanTasksCounts(inactiveList, activeList, completedList []KanbanItem, m *Model) {
	if m.ViewManager.KanbanListCursor == 0 {
		tasksCount := 0
		for _, item := range inactiveList {
			tasksCount += len(item.tasks)
		}
		m.ViewManager.KanbanTasksCount = tasksCount
	} else if m.ViewManager.KanbanListCursor == 1 {
		tasksCount := 0
		for _, item := range activeList {
			tasksCount += len(item.tasks)
		}
		m.ViewManager.KanbanTasksCount = tasksCount
	} else if m.ViewManager.KanbanListCursor == 2 {
		tasksCount := 0
		for _, item := range completedList {
			tasksCount += len(item.tasks)
		}
		m.ViewManager.KanbanTasksCount = tasksCount
	}
}

func renderKanbanList(m *Model, kanbanList []KanbanItem, boardWidth int, selectedList bool) string {
	renderedKanbanList := ""
	index := m.ViewManager.KanbanTaskCursor
	totalIndex := 0

	for _, kanbanItem := range kanbanList {
		tasks := kanbanItem.tasks
		filename := kanbanItem.filename

		renderedKanbanList = joinVertical(renderedKanbanList, renderFilename(filename, boardWidth))

		for _, task := range tasks {
			selected := false

			if m.ViewManager.IsKanbanTaskUpdated {
				if m.TaskManager.SelectedTask.textWithoutDates() == task.textWithoutDates() {
					selected = true
					m.SelectTask(task)
					m.FileManager.SelectFile(filename)
					m.ViewManager.IsKanbanTaskUpdated = false
					m.ViewManager.KanbanTaskCursor = totalIndex
				}
			} else if selectedList && index == 0 {
				selected = true
				m.FileManager.SelectFile(filename)
				m.SelectTask(task)
			}

			renderedKanbanList = joinVertical(renderedKanbanList, renderKanbanTask(task, boardWidth, m.TaskManager.DailySummaryDate, selected, m.ViewManager.IsWeeklyView))

			index--
			totalIndex++
		}
	}

	newViewport := viewport.Model{}
	newViewport.Width = boardWidth
	newViewport.Height = m.ViewManager.DetailsViewHeight
	newViewport.SetContent(renderedKanbanList)

	if selectedList {
		newViewport.LineDown(m.ViewManager.KanbanLineDownAmount())
	}

	return newViewport.View()
}

func renderBoard(title string, list []KanbanItem, m *Model, selectedBoard bool) string {
	boardWidth := (m.ViewManager.DetailsViewWidth - 6) / 3

	renderedTitle := kanbanBoardTitleStyle(colorForTitle(title)).Render(title)
	renderedList := renderKanbanList(m, list, boardWidth, selectedBoard)

	return boardContainerStyle(boardWidth, m.ViewManager.DetailsViewHeight, selectedBoard).Render(joinVertical(renderedTitle, renderedList))
}

func colorForTitle(title string) lipgloss.Color {
	if title == "Inactive" {
		return inactiveFileColor
	} else if title == "Active" {
		return activeFileColor
	} else if title == "Completed" {
		return completedFileColor
	}

	return white
}

func renderFilename(filename string, boardWidth int) string {
	return kanbanTaskTitleStyle.Render(filename)
}

func renderKanbanTask(task Task, boardWidth int, date string, selected bool, weekly bool) string {
	style := kanbanTaskStyle(boardWidth)

	if selected {
		style = highlightedKanbanTaskStyle(boardWidth)
	}

	return style.Render(TaskView{task: task, date: date, weekly: weekly, width: boardWidth}.RenderedKanbanText())
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

		line := file.FileNameWithoutExtension(m.FileManager.FileExtension)

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

	if m.IsAddTaskView() || m.IsAddSubTaskView() {
		itemDetails = m.NewTaskInput.View()

		if hasUnclosedDoubleSquareBrackets(m.NewTaskInput.Value()) {
			m.ViewManager.IsSuggestionsActive = true
			filterValue := peopleFilterValue(m.NewTaskInput.Value())
			peopleOptions := m.FileManager.PeopleFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)
			taskOptions := m.FileManager.TaskFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)

			peopleOptionsView := ""
			for _, option := range peopleOptions {
				person := strings.Split(option, m.FileManager.FileExtension)[0]
				peopleOptionsView = joinVertical(peopleOptionsView, suggestionTextStyle.Render(person))
			}

			peopleOptionViewTitle := suggestionTitleStyle.Render("People")

			if peopleOptionsView != "" {
				itemDetails = joinVertical(itemDetails, peopleOptionViewTitle, peopleOptionsView, "\n")
			}

			taskOptionsView := ""
			for _, option := range taskOptions {
				task := strings.Split(option, m.FileManager.FileExtension)[0]
				taskOptionsView = joinVertical(taskOptionsView, suggestionTextStyle.Render(task))
			}

			taskOptionViewTitle := suggestionTitleStyle.Render("Tasks")

			if taskOptionsView != "" {
				itemDetails = joinVertical(itemDetails, taskOptionViewTitle, taskOptionsView)
			}
		} else {
			m.ViewManager.IsSuggestionsActive = false
		}
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
	taskTitle := category[0 : len(category)-len(m.FileManager.FileExtension)]
	taskTitle += " (" + fmt.Sprintf("%d", activeTaskCount(m, category, date)) + " active, " + fmt.Sprintf("%d", incompleteTaskCount(m, category, date)) + " remaining)"

	return titleStyle.Render(taskTitle)
}

func incompleteTaskCount(m *Model, category string, date string) int {
	return len(m.TaskManager.TaskCollection.IncompleteTasks(category, date))
}

func activeTaskCount(m *Model, category string, date string) int {
	return len(m.TaskManager.TaskCollection.ActiveTasks(category, date))
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
		m.TaskManager.SelectedTask = task
		m.FileManager.SelectFile(task.FileName)
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
