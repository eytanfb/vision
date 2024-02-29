package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m *Model) string {
	content := "Something is wrong"

	if m.IsCompanyView() {
		content = RenderList(m, m.CompanyNames(), "company")
	} else if m.IsCategoryView() {
		content = RenderList(m, m.CategoryNames(), "category")
	} else if m.IsDetailsView() {
		if m.IsTaskDetailsFocus() {
			content = RenderTasks(m)
		} else {
			content = RenderFiles(m)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, RenderNavBar(m), content, RenderErrors(m.FileManager))
}

func RenderErrors(fm FileManager) string {
	var errors strings.Builder
	for _, err := range fm.Errors {
		errors.WriteString(err + "\n")
	}

	return errors.String()
}

func RenderList(m *Model, items []string, title string) string {
	list := ""
	containerStyle := lipgloss.NewStyle().Width(40).Height(m.ViewManager.Height - 20).Padding(1)

	cursor := m.GetCurrentCursor()

	for index, item := range items {
		line := "  "
		style := lipgloss.NewStyle()

		if index == cursor {
			line = "❯ "
			style = style.Bold(true)
		}

		line += item + "\n"
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	view := containerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, "Select a "+title+":\n", list))

	return view
}

func RenderFiles(m *Model) string {
	if !m.HasFiles() {
		return "No files found"
	}

	var style lipgloss.Style
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := lipgloss.NewStyle().Width(40).Height(m.ViewManager.Height - 20).Padding(1).Border(lipgloss.RoundedBorder())
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.Width - 60).Height(m.ViewManager.Height - 20).Padding(1).Border(lipgloss.RoundedBorder())
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
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	listContainer := listContainerStyle.Render(list)

	m.Viewport.SetContent(itemDetails)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	titleStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	title := titleStyle.Render("Select a file (" + fmt.Sprint(len(m.FileManager.Files)) + "):")
	view := lipgloss.JoinVertical(lipgloss.Top, title, container)

	return view
}

func RenderNavBar(m *Model) string {
	navbar := ""
	textStyle := lipgloss.NewStyle().Bold(true)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#000")).MarginBottom(1).Width(m.ViewManager.Width).Padding(1)
	if m.IsCompanyView() {
		navbar = textStyle.Render("Company selection")
	} else if m.IsCategoryView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > Category selection")
	} else if m.IsDetailsView() {
		navbar = textStyle.Render(m.GetCurrentCompanyName() + " > " + m.DirectoryManager.SelectedCategory)
	}

	return style.Render(navbar)
}

func RenderTasks(m *Model) string {
	if len(m.FileManager.Files) == 0 {
		return "No files found"
	}

	var style lipgloss.Style
	var tasks strings.Builder
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := lipgloss.NewStyle().Width(40).Height(m.ViewManager.Height - 20).Padding(1)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.ViewManager.Width - 60).Height(m.ViewManager.Height - 20).Padding(1)
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
			for index, task := range m.TaskManager.TaskCollection.Tasks {
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

	m.Viewport.SetContent(tasks.String())

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	view := lipgloss.JoinVertical(lipgloss.Top, container)

	return view
}
