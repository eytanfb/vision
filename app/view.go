package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	return ViewHandler(m)
}

func ViewHandler(m Model) string {
	content := "Something is wrong"

	if m.CurrentView == "companies" {
		companyNames := make([]string, len(m.Companies))
		for i, company := range m.Companies {
			companyNames[i] = company.DisplayName
		}

		content = RenderList(m, companyNames, "company")
	} else if m.CurrentView == "categories" {
		content = RenderList(m, m.Categories, "category")
	} else if m.CurrentView == "details" {
		content = RenderFiles(m)
	}

	return lipgloss.JoinVertical(lipgloss.Top, RenderNavBar(m), content)
}

func RenderList(m Model, items []string, title string) string {
	var style lipgloss.Style
	list := ""
	containerStyle := lipgloss.NewStyle().Width(40)
	cursor := m.CompaniesCursor
	if title == "category" {
		cursor = m.CategoriesCursor
	}

	for index, item := range items {
		line := ""
		style = lipgloss.NewStyle()

		if index == cursor {
			line += "❯ "
			style = style.Bold(true)
		} else {
			line += "  "
		}

		line += item + "\n"
		list = lipgloss.JoinVertical(lipgloss.Top, list, style.Render(line))
	}

	view := containerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, "Select a "+title+":\n", list))

	return view
}

func RenderFiles(m Model) string {
	if len(m.Files) == 0 {
		return "No files found"
	}

	var style lipgloss.Style
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := lipgloss.NewStyle().Width(40).Height(m.Height - 20)
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2).Width(m.Width - 60).Height(m.Height - 20)
	list := ""
	itemDetails := ""

	if m.ItemDetailsFocus {
		itemDetailsContainerStyle = itemDetailsContainerStyle.Border(lipgloss.RoundedBorder())
	} else {
		listContainerStyle = listContainerStyle.Border(lipgloss.RoundedBorder())
	}

	for index, file := range m.Files {
		line := ""
		style = lipgloss.NewStyle()

		if index == m.FilesCursor {
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

	m.Viewport.YPosition = 20
	m.Viewport.SetContent(itemDetails)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.Viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	titleStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	title := titleStyle.Render("Select a file (" + fmt.Sprint(len(m.Files)) + "):")
	view := lipgloss.JoinVertical(lipgloss.Top, title, container)

	return view
}

func RenderNavBar(m Model) string {
	navbar := ""
	textStyle := lipgloss.NewStyle().Bold(true)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#000")).Padding(1, 2)
	if m.CurrentView == "companies" {
		navbar = textStyle.Render("Company selection")
	} else if m.CurrentView == "categories" {
		navbar = textStyle.Render(m.SelectedCompany.DisplayName + " > Category selection")
	} else if m.CurrentView == "details" {
		navbar = textStyle.Render(m.SelectedCompany.DisplayName + " > " + m.SelectedCategory)
	}

	return style.Render(navbar)
}
