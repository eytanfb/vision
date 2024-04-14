package app

import (
	"sort"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	companyTextStyle     = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	selectedCompanyStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).Foreground(lipgloss.Color("#4CD137")).Bold(true)
	scheduledTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2D0A4"))
	startedTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0AAFC7"))
	completedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4CD137"))
	overdueTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EC4E20"))
)

func (m *Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m *Model) string {
	content := "Something is wrong"
	log.Info("ViewHandler is called")

	if m.IsCategoryView() {
		content = RenderList(m, m.CategoryNames(), "category")
	} else if m.IsDetailsView() {
		if m.IsTaskDetailsFocus() {
			content = RenderTasks(m)
		} else {
			content = RenderFiles(m)
		}
	}

	return joinVertical(RenderNavBar(m), content, RenderErrors(m))
}

func RenderErrors(m *Model) string {
	var errors strings.Builder
	for _, err := range m.Errors {
		errors.WriteString(err + "\n")
	}

	return errors.String()
}

func RenderAddTask(m *Model) string {
	return addTaskContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(m.NewTaskInput.View())
}

func RenderCompanies(m *Model, companies []string) string {
	result := ""

	for index, company := range companies {
		textStyle := companyTextStyle
		if index == m.DirectoryManager.CompaniesCursor {
			textStyle = selectedCompanyStyle
			result = joinHorizontal(result, textStyle.Render("["+company+"]"))
		}
	}

	return companiesContainerStyle(m.ViewManager.Width - 5).Render(result)
}

func RenderList(m *Model, items []string, title string) string {
	summaryView := ""
	list := ""

	cursor := m.GetCurrentCursor()

	for index, item := range items {
		list = joinVertical(list, createListItem(item, index, cursor))
	}

	view := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight).Render(joinVertical(list))

	if m.IsAddTaskView() {
		summaryView = m.NewTaskInput.View()
	} else if m.ViewManager.IsWeeklyView {
		summaryView = TaskSummaryToView(m, "weekly")
	} else {
		summaryView = TaskSummaryToView(m, "daily")
	}
	m.Viewport.SetContent(summaryView)
	summary := summaryContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.SummaryViewHeight).Render(m.Viewport.View())

	if m.ViewManager.HideSidebar {
		view = ""
	}

	return joinHorizontal(view, summary)
}

func TaskSummaryToView(m *Model, period string) string {
	tasksByFile, summaryDate := setDailySummaryValues(m)

	if period == "weekly" {
		tasksByFile, summaryDate = setWeeklySummaryValues(m)
	}

	keys := sortTaskKeys(tasksByFile)
	viewSort(keys, m)

	height := m.ViewManager.SummaryViewHeight
	containerTitleHeight := 2
	viewHeight := height - containerTitleHeight

	view := BuildSummaryView(m, keys, tasksByFile, m.ViewManager.DetailsViewWidth, summaryDate)

	containerTitle := taskSummaryContainerStyle(m.ViewManager.DetailsViewWidth, containerTitleHeight).Height(2).Render(summaryTitle(m, period))
	renderedView := taskSummaryContainerStyle(m.ViewManager.DetailsViewWidth, viewHeight).Render(view)

	return joinVertical(containerTitle, renderedView)
}

func setDailySummaryValues(m *Model) (map[string][]Task, string) {
	tasksByFile := m.TaskManager.Summary(m.GetCurrentCompanyName())
	summaryDate := m.TaskManager.DailySummaryDate

	return tasksByFile, summaryDate
}

func setWeeklySummaryValues(m *Model) (map[string][]Task, string) {
	startDate := m.TaskManager.WeeklySummaryStartDate
	endDate := m.TaskManager.WeeklySummaryEndDate
	tasksByFile := m.TaskManager.WeeklySummary(m.GetCurrentCompanyName(), startDate, endDate)
	summaryDate := m.TaskManager.WeeklySummaryEndDate

	return tasksByFile, summaryDate
}

func summaryTitle(m *Model, period string) string {
	title := "Daily Tasks for " + m.TaskManager.DailySummaryDate
	if period == "weekly" {
		title = "Weekly Tasks for " + m.TaskManager.WeeklySummaryStartDate + " - " + m.TaskManager.WeeklySummaryEndDate
	}

	return title
}

func sortTaskKeys(tasksByFile map[string][]Task) []string {
	keys := make([]string, 0, len(tasksByFile))
	for k := range tasksByFile {
		keys = append(keys, k)
	}

	return keys
}

func RenderFiles(m *Model) string {
	if !m.HasFiles() {
		return "No files found"
	}

	listContainerStyle := listContainerStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight, m.IsItemDetailsFocus())

	list, itemDetails := BuildFilesView(m)
	listContainer := listContainerStyle.Render(list)

	itemDetailsContainer := lipgloss.NewStyle().Width(m.ViewManager.DetailsViewWidth).MarginLeft(2).Border(lipgloss.NormalBorder()).Render(itemDetails)
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}

	container := lipgloss.NewStyle().Render(joinHorizontal(listContainer, itemDetailsContainer))

	return joinVertical(container)
}

func RenderNavBar(m *Model) string {
	companyColor := m.DirectoryManager.SelectedCompany.Color
	textStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(companyColor))
	container := lipgloss.NewStyle().Width(m.ViewManager.NavbarWidth).Padding(1).Border(lipgloss.NormalBorder())
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(m.ViewManager.NavbarWidth).Align(lipgloss.Center)
	navbar := textStyle.Render("Vision")

	if m.IsCategoryView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName())
	} else if m.IsDetailsView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > " + m.DirectoryManager.SelectedCategory + " > " + m.FileManager.SelectedFile.Name)
	}

	navbarView := lipgloss.JoinVertical(lipgloss.Top, style.Render(navbar))

	if m.ViewManager.ShowCompanies {
		navbarView = lipgloss.JoinHorizontal(lipgloss.Left, navbar, RenderCompanies(m, m.CategoryNames()))
	}

	view := container.Render(navbarView)

	return view
}

func RenderSidebar(m *Model) string {
	return ""
}

func RenderTasks(m *Model) string {
	itemDetailsContainerStyle := lipgloss.NewStyle().Width(m.ViewManager.DetailsViewWidth).Height(m.ViewManager.DetailsViewHeight).MarginLeft(2).Border(lipgloss.NormalBorder())

	if m.ViewManager.IsAddTaskView {
		addTaskView := m.NewTaskInput.View()

		m.Viewport.SetContent("Add new subtask\n" + addTaskView)

		return contentContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight).Render(m.Viewport.View())
	}

	date := m.TaskManager.DailySummaryDate

	tasks := m.TaskManager.TaskCollection.GetTasks(m.FileManager.currentFileName())

	tasksView := BuildTasksForFileView(m, tasks, date, m.TaskManager.TasksCursor)

	return joinVertical(itemDetailsContainerStyle.Render(tasksView))
}

func sidebarStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder())
}

func renderMarkdown(content string) string {
	out, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	markdown, err := out.Render(content)
	if err != nil {
		return ""
	}

	return markdown
}

func viewSort(filenames []string, m *Model) {
	sort.Slice(filenames, func(i, j int) bool {
		iInactive := m.TaskManager.TaskCollection.IsInactive(filenames[i])
		jInactive := m.TaskManager.TaskCollection.IsInactive(filenames[j])

		if iInactive {
			return false
		}

		if jInactive {
			return true
		}
		iCompletedTasks, iTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[i])
		jCompletedTasks, jTotalTasks := m.TaskManager.TaskCollection.Progress(filenames[j])

		iPercentage := float64(iCompletedTasks) / float64(iTotalTasks)
		jPercentage := float64(jCompletedTasks) / float64(jTotalTasks)

		iRoundedUp := int(iPercentage*10) * 10
		jRoundedUp := int(jPercentage*10) * 10

		if iRoundedUp == 100 {
			return false
		}

		if jRoundedUp == 100 {
			return true
		}

		if iRoundedUp == jRoundedUp {
			return filenames[i] < filenames[j]
		}

		return iRoundedUp > jRoundedUp
	})
}

func createListItem(item string, index int, cursor int) string {
	line := ""
	style := lipgloss.NewStyle()

	if index == cursor {
		style = style.Bold(true).Foreground(lipgloss.Color("#CB48B7"))
	}

	line += item + "\n"

	return style.Render(line)
}

func joinVertical(items ...string) string {
	return lipgloss.JoinVertical(lipgloss.Top, items...)
}

func joinHorizontal(items ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left, items...)
}

func itemDetailsContainerStyle(width, height int, isItemDetailsFocus bool) lipgloss.Style {
	style := contentContainerStyle(width, height)
	if isItemDetailsFocus {
		style = style.Copy().BorderForeground(lipgloss.Color("63"))
	}
	return style
}

func addTaskContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func summaryContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func contentContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Border(lipgloss.NormalBorder()).MarginLeft(2)
}

func taskSummaryContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1)
}

func taskTitleContainer(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(1).MarginTop(1)
}

func companiesContainerStyle(width int) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(width).Align(lipgloss.Center)
	return style
}

func listContainerStyle(width int, height int, isItemDetailsFocus bool) lipgloss.Style {
	style := sidebarStyle(width, height)
	if !isItemDetailsFocus {
		style = style.Copy().BorderForeground(lipgloss.Color("63"))
	}

	return style
}
