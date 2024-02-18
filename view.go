package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := "Something is wrong"

	if m.currentView == "companies" {
		companyNames := make([]string, len(m.companies))
		for i, company := range m.companies {
			companyNames[i] = company.DisplayName
		}

		content = m.RenderList(companyNames, "company")
	} else if m.currentView == "categories" {
		content = m.RenderList(m.categories, "category")
	} else if m.currentView == "details" {
		content = m.RenderFiles()
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.RenderNavBar(), content)
}

func (m Model) RenderList(items []string, title string) string {
	var style lipgloss.Style
	list := ""
	containerStyle := lipgloss.NewStyle().Width(30)

	for index, item := range items {
		line := ""
		style = lipgloss.NewStyle()

		if index == m.cursor {
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

func (m Model) RenderFiles() string {
	var style lipgloss.Style
	containerStyle := lipgloss.NewStyle()
	listContainerStyle := lipgloss.NewStyle()
	itemDetailsContainerStyle := lipgloss.NewStyle().MarginLeft(2)
	list := ""
	itemDetails := ""

	if m.itemDetailsFocus {
		itemDetailsContainerStyle = itemDetailsContainerStyle.Border(lipgloss.RoundedBorder())
	} else {
		listContainerStyle = listContainerStyle.Border(lipgloss.RoundedBorder())
	}

	for index, file := range m.files {
		line := ""
		style = lipgloss.NewStyle()

		if index == m.cursor {
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

	m.viewport.YPosition = 20
	m.viewport.SetContent(itemDetails)

	itemDetailsContainer := itemDetailsContainerStyle.Render(m.viewport.View())
	container := containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, listContainer, itemDetailsContainer))

	titleStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	title := titleStyle.Render("Select a file (" + fmt.Sprint(len(m.files)) + "):")
	view := lipgloss.JoinVertical(lipgloss.Top, title, container)

	return view
}

func (m Model) RenderNavBar() string {
	navbar := ""
	textStyle := lipgloss.NewStyle().Bold(true)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#000")).Padding(1, 2)
	if m.currentView == "companies" {
		navbar = textStyle.Render("Company selection")
	} else if m.currentView == "categories" {
		navbar = textStyle.Render(m.selectedCompany.DisplayName + " > Category selection")
	} else if m.currentView == "details" {
		navbar = textStyle.Render(m.selectedCompany.DisplayName + " > " + m.selectedCategory)
	}

	return style.Render(navbar)
}
