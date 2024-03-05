package app

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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

func RenderCompanies(m *Model, companies []string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(m.ViewManager.Width - 5).Align(lipgloss.Center)
	result := ""

	for index, company := range companies {
		textStyle := lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
		if index == m.DirectoryManager.CompaniesCursor {
			textStyle = textStyle.Bold(true).Background(lipgloss.Color("63"))
		}
		result = lipgloss.JoinHorizontal(lipgloss.Left, result, textStyle.Render("["+company+"]"))
	}

	return style.Render(result)
}

func RenderList(m *Model, items []string, title string) string {
	list := ""

	cursor := m.GetCurrentCursor()

	for index, item := range items {
		list = lipgloss.JoinVertical(lipgloss.Top, list, createListItem(item, index, cursor))
	}

	view := sidebarStyle(m.ViewManager.SidebarWidth, m.ViewManager.SidebarHeight).Render(lipgloss.JoinVertical(lipgloss.Top, list))

	summaryStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.DetailsViewWidth).Height(m.ViewManager.DetailsViewHeight).Padding(1).Border(lipgloss.NormalBorder())
	summaryView := TaskSummaryToView(m)
	summary := summaryStyle.Render(summaryView)

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

	containerStyle := lipgloss.NewStyle().Width(width).Padding(1)
	titleStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(width - 40).Bold(true).Foreground(lipgloss.Color("63"))
	startedTextStyle := lipgloss.NewStyle()

	view := ""
	for _, key := range keys {
		category := key
		tasks := tasksByFile[key]
		isInactive := m.TaskManager.TaskCollection.IsInactive(category)
		if isInactive {
			continue
		}
		completedTasks, totalTasks := m.TaskManager.TaskCollection.Progress(category)
		percentage := float64(completedTasks) / float64(totalTasks)
		if percentage == 1 {
			continue
		}
		roundedUpPercentage := int(percentage*10) * 10

		titleContainer := lipgloss.NewStyle().Width(width).MarginTop(1)
		taskTitle := category
		progressText := progressBar(completedTasks, totalTasks) + " " + fmt.Sprintf("%d%%", roundedUpPercentage)
		tasksView := ""
		incompleteTaskCount := 0
		for _, task := range tasks {
			if task.Completed {
				continue
			}
			incompleteTaskCount++
			tasks := ""
			text := task.Summary()
			if task.Started {
				text += " 🛫 " + DaysAgoFromString(task.StartDate)
				if !strings.Contains(progressText, "🛫") {
					progressText = strings.Replace(progressText, " ⏳", "", -1)
					progressText += " 🛫"
				}
			} else if task.Scheduled {
				text += " ⏳ " + DaysAgoFromString(task.ScheduledDate)
				if !strings.Contains(progressText, "⏳") && !strings.Contains(progressText, "🚨") && !strings.Contains(progressText, "🛫") {
					progressText += " ⏳"
				}
			}
			if task.IsOverdue() {
				text += " 🚨"
				if !strings.Contains(progressText, "🚨") {
					progressText = strings.Replace(progressText, " ⏳", "", -1)
					progressText = strings.Replace(progressText, " 🛫", "", -1)
					progressText += " 🚨"
				}
			}

			tasks = startedTextStyle.Render(text)
			tasksView = lipgloss.JoinVertical(lipgloss.Top, tasksView, tasks)
		}

		rightAlignedProgressText := titleStyle.Copy().Width(30).Align(lipgloss.Right).Render(progressText)
		taskTitle += " (" + fmt.Sprintf("%d", incompleteTaskCount) + " tasks remaining)"
		taskTitleView := lipgloss.JoinHorizontal(lipgloss.Left, titleStyle.Render(taskTitle), rightAlignedProgressText)
		tasksView = lipgloss.JoinVertical(lipgloss.Top, titleContainer.Render(taskTitleView), tasksView)
		view = lipgloss.JoinVertical(lipgloss.Top, view, tasksView)
	}

	containerTitleStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(width).Padding(1)
	containerTitle := containerTitleStyle.Render("Daily Tasks for " + time.Now().Format("2006-01-02"))
	renderedView := containerStyle.Render(view)

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
	if days < 2 {
		daysString = "day"
	}
	return fmt.Sprintf("%.0f %s ago", days, daysString)
}

func createListItem(item string, index int, cursor int) string {
	line := "  "
	style := lipgloss.NewStyle()

	if index == cursor {
		line = "❯ "
		style = style.Bold(true)
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
			line += "❯ "
			style = style.Bold(true)
			itemDetails = file.Name + "\n" + file.Content
		} else {
			line += "  "
		}

		line += file.Name
		if m.DirectoryManager.SelectedCategory == "tasks" {
			isInactive := m.TaskManager.TaskCollection.IsInactive(file.Name)

			if isInactive {
				style = style.Copy().Foreground(lipgloss.Color("#A0A0A0"))
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

	markdown := renderMarkdown(m.ViewManager.Width, itemDetails)
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
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
		line := ""
		style = lipgloss.NewStyle()

		if index == m.FileManager.FilesCursor {
			line += "❯ "
			style = style.Bold(true)
			for index, task := range m.TaskManager.TaskCollection.GetTasks(file.Name) {
				if index == m.TaskManager.TasksCursor {
					tasks.WriteString("❯ ")
				}
				tasks.WriteString(task.String() + "\n\n")
			}
		} else {
			line += "  "
		}

		line += file.Name
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	listContainer := listContainerStyle.Render(list)

	markdown := renderMarkdown(m.ViewManager.Width, tasks.String())
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	if m.ViewManager.HideSidebar {
		listContainer = ""
	}
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

	return view
}

func sidebarStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder())
}

func renderMarkdown(width int, content string) string {
	background := "light"

	if lipgloss.HasDarkBackground() {
		background = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width-40),
		glamour.WithStandardStyle(background),
		glamour.WithStylesFromJSONFile("/Users/eytananjel/Code/vision/dark_modified.json"),
	)

	out, err := r.Render(content)
	if err != nil {
		return ""
	}

	return out
}
