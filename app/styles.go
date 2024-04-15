package app

import "github.com/charmbracelet/lipgloss"

var (
	white                          = lipgloss.Color("#FFFFFF")
	taskDateColor                  = lipgloss.Color("#9A9CCD")
	highlightedTaskBackgroundColor = lipgloss.Color("#474747")
	completedColor                 = lipgloss.Color("#4CD137")
	scheduledColor                 = lipgloss.Color("#F2D0A4")
	startedColor                   = lipgloss.Color("#0AAFC7")
	overdueColor                   = lipgloss.Color("#EC4E20")
	completedFileColor             = lipgloss.Color("#4DA165")
	inactiveFileColor              = lipgloss.Color("#A0A0A0")
	highlightedTextColor           = lipgloss.Color("#CB48B7")
	summaryTitleColor              = lipgloss.Color("#9A9CCD")
)

var (
	defaultTextStyle        = lipgloss.NewStyle().Foreground(white)
	companyTextStyle        = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	selectedCompanyStyle    = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).Foreground(completedColor).Bold(true)
	scheduledTextStyle      = lipgloss.NewStyle().Foreground(scheduledColor)
	startedTextStyle        = lipgloss.NewStyle().Foreground(startedColor)
	completedTextStyle      = lipgloss.NewStyle().Foreground(completedColor)
	overdueTextStyle        = lipgloss.NewStyle().Foreground(overdueColor)
	taskFileTitleStyle      = lipgloss.NewStyle().Foreground(white).Bold(true).Underline(true)
	completedFileStyle      = lipgloss.NewStyle().Foreground(completedFileColor)
	inactiveFileStyle       = lipgloss.NewStyle().Foreground(inactiveFileColor)
	highlightedTextStyle    = lipgloss.NewStyle().Foreground(highlightedTextColor).Bold(true)
	renderedStatusTextStyle = lipgloss.NewStyle().Width(15).Align(lipgloss.Right)
	iconStyle               = lipgloss.NewStyle().MarginRight(1)
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

func progressTextStyle(style lipgloss.Style) lipgloss.Style {
	return style.Copy().Width(35).Align(lipgloss.Right)
}

func taskStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(2).Width(width)
}

func kanbanTaskTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(summaryTitleColor)
}

func kanbanTaskStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).MarginBottom(1)
}

func highlightedKanbanTaskStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Background(highlightedTaskBackgroundColor).MarginBottom(1)
}

func kanbanBoardTitleStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Underline(true).Foreground(color)
}

func datesContainerStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(2).Width(width)
}

func dateStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(highlightedTaskBackgroundColor).Foreground(taskDateColor).PaddingRight(2)
}

func tasksStringStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Background(highlightedTaskBackgroundColor).Width(width).PaddingTop(1).PaddingBottom(1).MarginTop(1).MarginBottom(1)
}

func summaryTitleStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().Align(lipgloss.Left).Width(width - 40).Bold(true).Foreground(summaryTitleColor)
}

func inactiveTitleStyle() lipgloss.Style {
	return taskFileTitleStyle.Copy().Foreground(inactiveFileColor)
}

func completedTitleStyle() lipgloss.Style {
	return taskFileTitleStyle.Copy().Foreground(completedFileColor)
}

func renderedActiveListStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginTop(1).MarginBottom(2)
}

func renderedInactiveListStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginTop(1).MarginBottom(3)
}

func renderedCompletedListStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
}

func taskItemDetailsContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).MarginLeft(2).Border(lipgloss.NormalBorder())
}

func boardContainerStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).Border(lipgloss.NormalBorder())
}
