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
			} else {
				// Remaining tasks (including past-due scheduled but not started) go to inactive
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

func renderKanbanList(m *Model, kanbanList []KanbanItem, boardWidth int, selectedList bool) string {
	var renderedItems []string
	index := m.ViewManager.KanbanTaskCursor
	totalIndex := 0

	for _, kanbanItem := range kanbanList {
		tasks := kanbanItem.tasks
		filename := kanbanItem.filename

		// Add filename header
		renderedItems = append(renderedItems, renderFilename(filename, boardWidth))

		for _, task := range tasks {
			// Pure view logic - only check state, never mutate
			selected := false

			if selectedList && index == 0 {
				selected = true
			}

			// Collect rendered task
			renderedItems = append(renderedItems, renderKanbanTask(task, boardWidth, m.TaskManager.DailySummaryDate, selected, m.ViewManager.IsWeeklyView))

			index--
			totalIndex++
		}
	}

	// Join all items once
	renderedKanbanList := joinVertical(renderedItems...)

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
	boardWidth := (m.ViewManager.DetailsViewWidth - KanbanBoardSpacing) / KanbanColumnCount

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
	content, _, _ := buildTaskRowsForFileView(m, tasks, date, cursor)
	return content
}

type rowAccumulator struct {
	rows  []string
	line  int
	start int
	end   int
}

func newRowAccumulator(leading ...string) *rowAccumulator {
	acc := &rowAccumulator{start: -1, end: -1}
	for _, row := range leading {
		acc.add(row, false)
	}
	return acc
}

func (acc *rowAccumulator) add(row string, isCursor bool) {
	height := lipgloss.Height(row)
	if isCursor {
		acc.start, acc.end = acc.line, acc.line+height
	}
	acc.rows = append(acc.rows, row)
	acc.line += height
}

func (acc *rowAccumulator) content() string {
	return joinVertical(acc.rows...)
}

func buildTaskRowsForFileView(m *Model, tasks []Task, date string, cursor int) (string, int, int) {
	acc := newRowAccumulator("")
	for index, task := range tasks {
		row := renderTaskForFileView(m, task, date, cursor, index)
		if row == "" {
			continue
		}
		acc.add(row, index == cursor)
	}

	return acc.content(), acc.start, acc.end
}

func BuildFilesView(m *Model, hiddenSidebar bool) (string, string) {
	itemDetails := ""
	list := ""
	if !hiddenSidebar {
		list = m.FilesViewport.View()
	}

	if m.DirectoryManager.SelectedCategory == "tasks" {
		return list, m.Viewport.View()
	}

	if m.IsAddTaskView() || m.IsAddSubTaskView() {
		itemDetails = m.NewTaskInput.View()

		// Only render suggestions if they're active (state set in Update, not here)
		if hasUnclosedDoubleSquareBrackets(m.NewTaskInput.Value()) {
			filterValue := peopleFilterValue(m.NewTaskInput.Value())
			peopleOptions := m.FileManager.PeopleFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)
			taskOptions := m.FileManager.TaskFilenames(&m.DirectoryManager, &m.TaskManager, filterValue)

			// Use strings.Builder for efficient string concatenation
			if len(peopleOptions) > 0 {
				var peopleBuilder strings.Builder
				for i, option := range peopleOptions {
					person := strings.Split(option, m.FileManager.FileExtension)[0]
					peopleBuilder.WriteString(suggestionTextStyle.Render(person))
					if i < len(peopleOptions)-1 {
						peopleBuilder.WriteRune('\n')
					}
				}

				peopleOptionViewTitle := suggestionTitleStyle.Render("People")
				itemDetails += "\n" + peopleOptionViewTitle + "\n" + peopleBuilder.String() + "\n"
			}

			// Use strings.Builder for task options
			if len(taskOptions) > 0 {
				var taskBuilder strings.Builder
				for i, option := range taskOptions {
					task := strings.Split(option, m.FileManager.FileExtension)[0]
					taskBuilder.WriteString(suggestionTextStyle.Render(task))
					if i < len(taskOptions)-1 {
						taskBuilder.WriteRune('\n')
					}
				}

				taskOptionViewTitle := suggestionTitleStyle.Render("Tasks")
				itemDetails += "\n" + taskOptionViewTitle + "\n" + taskBuilder.String()
			}
		}
	} else {
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
		width:  m.ViewManager.DetailsViewWidth - TaskWidthPadding,
	}.RenderedText()

	return joinVertical(tasksView, tasksString)
}

func renderTaskForFileView(m *Model, task Task, date string, cursor int, index int) string {
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

	return tasksString
}

func buildSidebarFilesContent(m *Model, width int) (string, int, int) {
	m.FileManager.filesMu.RLock()
	files := m.FileManager.Files
	m.FileManager.filesMu.RUnlock()

	cursor := m.FileManager.FilesCursor
	wrap := lipgloss.NewStyle().Width(max(width, 1))

	if m.DirectoryManager.SelectedCategory == "tasks" {
		return buildTaskSidebarContent(m, files, cursor, wrap)
	}

	acc := newRowAccumulator()
	for index, file := range files {
		style := defaultTextStyle
		if index == cursor {
			style = highlightedTextStyle
		}

		acc.add(wrap.Render(style.Render(file.FileNameWithoutExtension(m.FileManager.FileExtension))), index == cursor)
	}

	return acc.content(), acc.start, acc.end
}

type taskSidebarSection struct {
	title string
	frame lipgloss.Style
	rows  *rowAccumulator
}

func newTaskSidebarSection(title string, frame lipgloss.Style) *taskSidebarSection {
	return &taskSidebarSection{title: title, frame: frame, rows: newRowAccumulator("")}
}

const (
	activeSection = iota
	inactiveSection
	completedSection
)

func buildTaskSidebarContent(m *Model, files []FileInfo, cursor int, wrap lipgloss.Style) (string, int, int) {
	sections := []*taskSidebarSection{
		activeSection:    newTaskSidebarSection(taskFileTitleStyle.Render("Active"), renderedActiveListStyle()),
		inactiveSection:  newTaskSidebarSection(inactiveTitleStyle().Render("Inactive"), renderedInactiveListStyle()),
		completedSection: newTaskSidebarSection(completedTitleStyle().Render("Complete"), renderedCompletedListStyle()),
	}

	for index, file := range files {
		style := defaultTextStyle
		if index == cursor {
			style = highlightedTextStyle
		}

		line := file.FileNameWithoutExtension(m.FileManager.FileExtension)
		completed, total := m.TaskManager.TaskCollection.Progress(file.Name)

		target := activeSection
		switch {
		case total > 0 && completed == total:
			target = completedSection
			if index != cursor {
				style = completedFileStyle
			}
		case m.TaskManager.TaskCollection.IsInactive(file.Name):
			target = inactiveSection
			if index != cursor {
				style = inactiveFileStyle
			}
		default:
			line = fmt.Sprintf("[%d/%d] ", completed, total) + line
		}

		sections[target].rows.add(wrap.Render(style.Render(line)), index == cursor)
	}

	blocks := make([]string, 0, len(sections)*2)
	line, start, end := 0, -1, -1
	for _, section := range sections {
		titleHeight := lipgloss.Height(section.title)
		rendered := section.frame.Render(section.rows.content())

		if section.rows.start >= 0 {
			offset := line + titleHeight + section.frame.GetMarginTop()
			start, end = offset+section.rows.start, offset+section.rows.end
		}

		blocks = append(blocks, section.title, rendered)
		line += titleHeight + lipgloss.Height(rendered)
	}

	return joinVertical(blocks...), start, end
}

func progressBar(completed, total int) string {
	if total == 0 {
		return "[          ]"
	}

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

	if totalTasksCount == 0 {
		return "[          ] 0%"
	}

	percentage := float64(completedTasksCount) / float64(totalTasksCount)
	roundedUpPercentage := int(percentage*10) * 10

	return progressBar(completedTasksCount, totalTasksCount) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
}
