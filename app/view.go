package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	companyTextStyle     = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	selectedCompanyStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).Foreground(lipgloss.Color("#C0DFA1")).Bold(true)
	startedTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#035E7B"))
	completedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#1C3A13"))
	scheduledTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2D0A4"))
	overdueTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EC4E20"))
)

func (m *Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m *Model) string {
	content := "Something is wrong"

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
	tasksByFile := m.TaskManager.Summary(m.GetCurrentCompanyName())

	if period == "weekly" {
		startDate := m.TaskManager.WeeklySummaryStartDate
		endDate := m.TaskManager.WeeklySummaryEndDate
		tasksByFile = m.TaskManager.WeeklySummary(m.GetCurrentCompanyName(), startDate, endDate)
	}

	keys := sortTaskKeys(tasksByFile)
	viewSort(keys, &tasksByFile, m)

	width := m.ViewManager.DetailsViewWidth

	view := BuildSummaryView(m, keys, tasksByFile, width)

	containerTitle := taskSummaryContainerStyle(width).Render("Daily Tasks for " + time.Now().Format("2006-01-02"))
	renderedView := taskSummaryContainerStyle(width).Render(view)

	return joinVertical(containerTitle, renderedView)
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

	containerStyle := lipgloss.NewStyle()
	listContainerStyle := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight)
	itemDetailsContainerStyle := itemDetailsContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight)
	list := ""
	itemDetails := ""

	if m.IsItemDetailsFocus() {
		itemDetailsContainerStyle = itemDetailsContainerStyle.BorderForeground(lipgloss.Color("63"))
	} else {
		listContainerStyle = listContainerStyle.Copy().BorderForeground(lipgloss.Color("63"))
	}

	completedList := ""
	incompleteList := ""
	for index, file := range sortedFiles(m) {
		line := ""
		style := lipgloss.NewStyle()

		if index == m.FileManager.FilesCursor {
			style = style.Bold(true).Foreground(lipgloss.Color("#CB48B7"))
			itemDetails = file.FileNameWithoutExtension() + "\n" + file.Content
			m.FileManager.SelectedFile = file
		}

		line += file.FileNameWithoutExtension()
		if m.DirectoryManager.SelectedCategory == "tasks" {
			list, incompleteList, completedList = buildTasksView(m, line, index, list, file, style, incompleteList, completedList)
		} else {
			list = joinVertical(list, style.Render(line))
		}
	}

	listContainer := listContainerStyle.Render(list)

	if m.IsAddTaskView() {
		itemDetails = m.NewTaskInput.View()
	} else {
		markdown := renderMarkdown(wordwrap.String(itemDetails, m.ViewManager.DetailsViewWidth))
		m.Viewport.SetContent(markdown)
		itemDetails = m.Viewport.View()
	}

	itemDetailsContainer := itemDetailsContainerStyle.Render(itemDetails)
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}

	container := containerStyle.Render(joinHorizontal(listContainer, itemDetailsContainer))

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
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > " + m.DirectoryManager.SelectedCategory)
	}

	view := container.Render(joinVertical(style.Render(navbar), RenderCompanies(m, m.CompanyNames())))

	return view
}

func RenderSidebar(m *Model) string {
	return ""
}

func RenderTasks(m *Model) string {
	if len(m.FileManager.Files) == 0 {
		return "No files found"
	}

	var style lipgloss.Style
	var tasks strings.Builder
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth).Height(m.ViewManager.DetailsViewHeight).Padding(1)
	list := ""

	if m.IsItemDetailsFocus() {
		itemDetailsContainerStyle = itemDetailsContainerStyle.Border(lipgloss.RoundedBorder()).Padding(1)
	} else {
		listContainerStyle = listContainerStyle.Border(lipgloss.RoundedBorder()).Padding(1)
	}

	for index, file := range m.FileManager.Files {
		line := "  "
		style = lipgloss.NewStyle()

		if index == m.FileManager.FilesCursor {
			line += "❯ "
			style = style.Bold(true)
			for index, task := range m.TaskManager.TaskCollection.GetTasks(file.Name) {
				tasks.WriteString(writeTaskString(task, index, m.TaskManager.TasksCursor))
			}
		}

		line += file.Name
		list = joinVertical(list, style.Render(line))
	}

	listContainer := listContainerStyle.Render(list)

	markdown := renderMarkdown(tasks.String())
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}
	container := containerStyle.Render(joinHorizontal(listContainer, itemDetailsContainer))

	view := joinVertical(container)

	return view
}

func writeTaskString(task Task, index int, cursor int) string {
	viewCursor := "  "
	if index == cursor {
		viewCursor = "❯ "
	}

	return viewCursor + task.String() + "\n\n"
}

func sidebarStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder())
}

func renderMarkdown(content string) string {
	background := "light"

	if lipgloss.HasDarkBackground() {
		background = "dark"
	}

	out, err := glamour.Render(content, background)
	if err != nil {
		return ""
	}

	return out
}

func addIconToProgressText(progressText, icon string) string {
	progressText = strings.Replace(progressText, " ⏳", "", -1)
	progressText = strings.Replace(progressText, " 🛫", "", -1)
	progressText = strings.Replace(progressText, " ✅", "", -1)
	progressText = strings.Replace(progressText, " 🚨", "", -1)
	return progressText + " " + icon
}

func viewSort(filenames []string, tasksByFile *map[string][]Task, m *Model) {
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

func itemDetailsContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func addTaskContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func summaryContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func contentContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder()).MarginLeft(2)
}

func taskSummaryContainerTitleStyle(width int) lipgloss.Style {
	return taskSummaryContainerStyle(width).Copy().Align(lipgloss.Center)
}

func taskSummaryContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Padding(1)
}

func taskTitleContainer(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).MarginTop(1)
}

func companiesContainerStyle(width int) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(width).Align(lipgloss.Center)
	return style
}

func summaryTitleStyle(width int) lipgloss.Style {
	summaryStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(width - 40).Bold(true).Foreground(lipgloss.Color("63"))

	return summaryStyle
}

func sortedFiles(m *Model) []FileInfo {
	filenames := []string{}
	for _, file := range m.FileManager.Files {
		filenames = append(filenames, file.Name)
	}

	viewSort(filenames, &m.TaskManager.TaskCollection.TasksByFile, m)

	sortedFiles := []FileInfo{}
	for _, filename := range filenames {
		for _, file := range m.FileManager.Files {
			if file.Name == filename {
				sortedFiles = append(sortedFiles, file)
			}
		}
	}

	return sortedFiles
}

func buildTasksView(m *Model, line string, index int, list string, file FileInfo, style lipgloss.Style, incompleteList string, completedList string) (string, string, string) {
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
		incompleteList = joinVertical(incompleteList, style.Render(line))
	} else {
		completedText = lipgloss.NewStyle().Render(text)
		line += " " + completedText
		incompleteList = joinVertical(incompleteList, style.Render(line))
	}

	return joinVertical(incompleteList, lipgloss.NewStyle().MarginTop(2).Render(completedList)), incompleteList, completedList
}
