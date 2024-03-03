package app

import (
	"fmt"
	"strings"

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

	view := sidebarStyle(m.ViewManager.Height).Render(lipgloss.JoinVertical(lipgloss.Top, list))

	return view
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
	listContainerStyle := sidebarStyle(m.ViewManager.Height)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.Width - 60).Height(m.ViewManager.Height - 20).Padding(1).Border(lipgloss.NormalBorder())
	list := ""
	itemDetails := ""

	if m.IsItemDetailsFocus() {
		itemDetailsContainerStyle = itemDetailsContainerStyle.BorderForeground(lipgloss.Color("63"))
	} else {
		listContainerStyle = listContainerStyle.BorderForeground(lipgloss.Color("63"))
	}

	for index, file := range m.FileManager.Files {
		line := ""
		style = lipgloss.NewStyle()

		if index == m.FileManager.FilesCursor {
			line += "❯ "
			style = style.Bold(true)
			itemDetails = file.Content
		} else {
			line += "  "
		}

		line += file.Name
		if m.DirectoryManager.SelectedCategory == "tasks" {
			completed, total := m.TaskManager.TaskCollection.Progress(file.Name)
			text := fmt.Sprint(completed) + "/" + fmt.Sprint(total)
			var completedText string
			if completed == total {
				completedText = lipgloss.NewStyle().Background(lipgloss.Color("#AAC0AF")).Render(text)
			} else {
				completedText = lipgloss.NewStyle().Render(text)
			}
			line += " " + completedText
		}
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	listContainer := listContainerStyle.Render(list)

	markdown := renderMarkdown(m.ViewManager.Width, itemDetails)
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

	return view
}

func RenderNavBar(m *Model) string {
	textStyle := lipgloss.NewStyle().Bold(true)
	container := lipgloss.NewStyle().Width(m.ViewManager.Width - 5).Padding(1).Border(lipgloss.NormalBorder())
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Width(m.ViewManager.Width - 5).Align(lipgloss.Center)
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
	listContainerStyle := sidebarStyle(m.ViewManager.Height)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.Width - 60).Height(m.ViewManager.Height - 5).Padding(1)
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
	fmt.Println(markdown)
	m.Viewport.SetContent(markdown)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

	return view
}

func sidebarStyle(windowHeight int) lipgloss.Style {
	return lipgloss.NewStyle().Width(40).Height(windowHeight - 10).Padding(1).Border(lipgloss.NormalBorder())
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
