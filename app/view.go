package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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

	return lipgloss.JoinVertical(lipgloss.Top, RenderNavBar(m), content, RenderErrors(m))
}

func RenderErrors(m *Model) string {
	var errors strings.Builder
	for _, err := range m.Errors {
		errors.WriteString(err + "\n")
	}

	return errors.String()
}

func RenderAddTask(m *Model) string {
	containerStyle := lipgloss.NewStyle().Width(m.ViewManager.DetailsViewWidth).Height(m.ViewManager.DetailsViewHeight).Padding(1).Border(lipgloss.NormalBorder())
	return containerStyle.Render(m.NewTaskInput.View())
}

func RenderCompanies(m *Model, companies []string) string {
	result := ""

	for index, company := range companies {
		textStyle := companyTextStyle
		if index == m.DirectoryManager.CompaniesCursor {
			textStyle = selectedCompanyStyle
		}
		result = lipgloss.JoinHorizontal(lipgloss.Left, result, textStyle.Render("["+company+"]"))
	}

	return companiesContainerStyle(m.ViewManager.Width - 5).Render(result)
}

func RenderList(m *Model, items []string, title string) string {
	summaryView := ""
	list := ""

	cursor := m.GetCurrentCursor()

	for index, item := range items {
		list = lipgloss.JoinVertical(lipgloss.Top, list, createListItem(item, index, cursor))
	}

	view := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight).Render(lipgloss.JoinVertical(lipgloss.Top, list))

	if m.IsAddTaskView() {
		summaryView = m.NewTaskInput.View()
	} else {
		summaryView = TaskSummaryToView(m)
	}
	summary := summaryContainerStyle(m.ViewManager.DetailsViewWidth, m.ViewManager.SummaryViewHeight).Render(summaryView)

	if m.ViewManager.HideSidebar {
		view = ""
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, view, summary)
}

func TaskSummaryToView(m *Model) string {
	tasksByFile := m.TaskManager.Summary(m.GetCurrentCompanyName())

	keys := make([]string, 0, len(tasksByFile))
	for k := range tasksByFile {
		keys = append(keys, k)
	}

	viewSort(keys, &tasksByFile, m)

	width := m.ViewManager.DetailsViewWidth

	view := BuildSummaryView(m, keys, tasksByFile, width)

	containerTitle := taskSummaryContainerStyle(width).Render("Daily Tasks for " + time.Now().Format("2006-01-02"))
	renderedView := taskSummaryContainerStyle(width).Render(view)

	return lipgloss.JoinVertical(lipgloss.Top, containerTitle, renderedView)
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

func DaysAgoFromString(date string) string {
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

func createListItem(item string, index int, cursor int) string {
	line := ""
	style := lipgloss.NewStyle()

	if index == cursor {
		style = style.Bold(true).Foreground(lipgloss.Color("#73A580"))
	}

	line += item + "\n"

	return style.Render(line)
}

func RenderFiles(m *Model) string {
	if !m.HasFiles() {
		return "No files found"
	}

	var style lipgloss.Style
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth).Height(m.ViewManager.DetailsViewHeight).Padding(1).Border(lipgloss.NormalBorder())
	list := ""
	itemDetails := ""

	if m.IsItemDetailsFocus() {
		itemDetailsContainerStyle = itemDetailsContainerStyle.BorderForeground(lipgloss.Color("63"))
	} else {
		listContainerStyle = listContainerStyle.BorderForeground(lipgloss.Color("63"))
	}

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

	completedList := ""
	incompleteList := ""
	for index, file := range sortedFiles {
		line := ""
		style = lipgloss.NewStyle()

		if index == m.FileManager.FilesCursor {
			style = style.Bold(true).Foreground(lipgloss.Color("#73A580"))
			itemDetails = file.FileNameWithoutExtension() + "\n" + file.Content
			m.FileManager.SelectedFile = file
		}

		line += file.FileNameWithoutExtension()
		if m.DirectoryManager.SelectedCategory == "tasks" {
			isInactive := m.TaskManager.TaskCollection.IsInactive(file.Name)

			if isInactive {
				if index != m.FileManager.FilesCursor {
					style = style.Copy().Foreground(lipgloss.Color("#A0A0A0"))
				}
				incompleteList = lipgloss.JoinVertical(lipgloss.Top, incompleteList, style.Render(line))
			} else {
				completed, total := m.TaskManager.TaskCollection.Progress(file.Name)
				text := fmt.Sprintf("%d/%d", completed, total)
				var completedText string
				if completed == total {
					style = style.Copy().Foreground(lipgloss.Color("#4DA165"))
					completedList = lipgloss.JoinVertical(lipgloss.Top, completedList, style.Render(line))
				} else {
					completedText = lipgloss.NewStyle().Render(text)
					line += " " + completedText
					incompleteList = lipgloss.JoinVertical(lipgloss.Top, incompleteList, style.Render(line))
				}
			}
			list = lipgloss.JoinVertical(lipgloss.Top, incompleteList, lipgloss.NewStyle().MarginTop(2).Render(completedList))
		} else {
			list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
		}
	}

	listContainer := listContainerStyle.Render(list)

	if m.IsAddTaskView() {
		itemDetails = m.NewTaskInput.View()
	} else {
		markdown := renderMarkdown(m.ViewManager.Width, itemDetails)
		m.Viewport.SetContent(markdown)
		itemDetails = m.Viewport.View()
	}

	itemDetailsContainer := itemDetailsContainerStyle.Render(itemDetails)
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

	return view
}

func RenderNavBar(m *Model) string {
	textStyle := lipgloss.NewStyle().Bold(true)
	container := lipgloss.NewStyle().Width(m.ViewManager.NavbarWidth).Padding(1).Border(lipgloss.NormalBorder())
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(m.ViewManager.NavbarWidth).Align(lipgloss.Center)
	navbar := textStyle.Render("Vision")

	if m.IsCategoryView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName())
	} else if m.IsDetailsView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > " + m.DirectoryManager.SelectedCategory)
	}

	view := container.Render(lipgloss.JoinVertical(lipgloss.Top, style.Render(navbar), RenderCompanies(m, m.CompanyNames())))

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
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	listContainer := listContainerStyle.Render(list)

	markdown := renderMarkdown(m.ViewManager.DetailsViewWidth, tasks.String())
	m.Errors = append(m.Errors, fmt.Sprintf("%d", m.ViewManager.DetailsViewHeight))
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

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

func renderMarkdown(width int, content string) string {
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

func companiesContainerStyle(width int) lipgloss.Style {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(width).Align(lipgloss.Center)
	return style
}

func summaryContainerStyle(width, height int) lipgloss.Style {
	summaryStyle := lipgloss.NewStyle().MarginLeft(2).Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder())

	return summaryStyle
}

func summaryTitleStyle(width int) lipgloss.Style {
	summaryStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(width - 40).Bold(true).Foreground(lipgloss.Color("63"))

	return summaryStyle
}

func taskSummaryContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Padding(1)
}

func taskSummaryContainerTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Padding(1)
}

func taskTitleContainer(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).MarginTop(1)
}

func addIconToProgressText(progressText, icon string) string {
	progressText = strings.Replace(progressText, " ⏳", "", -1)
	progressText = strings.Replace(progressText, " 🛫", "", -1)
	return progressText + " " + icon
}
