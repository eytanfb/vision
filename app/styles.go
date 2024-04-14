package app

import "github.com/charmbracelet/lipgloss"

var (
	white = lipgloss.Color("#FFF")
)

var (
	companyTextStyle     = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	selectedCompanyStyle = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).Foreground(lipgloss.Color("#4CD137")).Bold(true)
	scheduledTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2D0A4"))
	startedTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0AAFC7"))
	completedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#4CD137"))
	overdueTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#EC4E20"))
)

func navbarTextStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))
}

func navbarContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Padding(1).Border(lipgloss.NormalBorder())
}

func navbarStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(white).Width(width).Align(lipgloss.Center)
}

func sidebarStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Padding(1).Border(lipgloss.NormalBorder())
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
	style := lipgloss.NewStyle().Foreground(white).Width(width - 5).Align(lipgloss.Center)
	return style
}

func listContainerStyle(width int, height int, isItemDetailsFocus bool) lipgloss.Style {
	style := sidebarStyle(width, height)
	if !isItemDetailsFocus {
		style = style.Copy().BorderForeground(lipgloss.Color("63"))
	}

	return style
}

func filesItemDetailsContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).MarginLeft(2).Border(lipgloss.NormalBorder())
}
