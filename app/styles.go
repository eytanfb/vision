package app

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

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
	activeFileColor                = lipgloss.Color("#FFFFFF")
	highlightedTextColor           = lipgloss.Color("#CB48B7")
	summaryTitleColor              = lipgloss.Color("#9A9CCD")
	suggestionTitleColor           = lipgloss.Color("#9A9CCD")
	suggestionTextColor            = lipgloss.Color("#9A9CCD")
	selectedSuggestionTextColor    = lipgloss.Color("#CB48B7")
	priorityTextColor              = lipgloss.Color("#EC4E20")
)

// StyleCache caches lipgloss.Style instances by their dimensions to avoid
// recreating styles on every render. Cache is invalidated on window resize.
type StyleCache struct {
	kanbanTaskStyles      map[int]lipgloss.Style
	highlightedKanbanTask map[int]lipgloss.Style
	taskStyles            map[int]lipgloss.Style
	containerStyles       map[string]lipgloss.Style // key: "width-height"
	boardStyles           map[string]lipgloss.Style // key: "width-height-selected"
	mu                    sync.RWMutex
}

var globalStyleCache = &StyleCache{
	kanbanTaskStyles:      make(map[int]lipgloss.Style),
	highlightedKanbanTask: make(map[int]lipgloss.Style),
	taskStyles:            make(map[int]lipgloss.Style),
	containerStyles:       make(map[string]lipgloss.Style),
	boardStyles:           make(map[string]lipgloss.Style),
}

// InvalidateAll clears all cached styles. Called on window resize to prevent
// unbounded cache growth and ensure styles reflect new dimensions.
func (sc *StyleCache) InvalidateAll() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.kanbanTaskStyles = make(map[int]lipgloss.Style)
	sc.highlightedKanbanTask = make(map[int]lipgloss.Style)
	sc.taskStyles = make(map[int]lipgloss.Style)
	sc.containerStyles = make(map[string]lipgloss.Style)
	sc.boardStyles = make(map[string]lipgloss.Style)
}

var (
	defaultTextStyle            = lipgloss.NewStyle().Foreground(white)
	companyTextStyle            = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	selectedCompanyStyle        = lipgloss.NewStyle().MarginLeft(2).MarginRight(2).Foreground(completedColor).Bold(true)
	scheduledTextStyle          = lipgloss.NewStyle().Foreground(scheduledColor)
	startedTextStyle            = lipgloss.NewStyle().Foreground(startedColor)
	completedTextStyle          = lipgloss.NewStyle().Foreground(completedColor)
	overdueTextStyle            = lipgloss.NewStyle().Foreground(overdueColor)
	priorityTextStyle           = lipgloss.NewStyle().Foreground(priorityTextColor)
	taskFileTitleStyle          = lipgloss.NewStyle().Foreground(white).Bold(true).Underline(true)
	completedFileStyle          = lipgloss.NewStyle().Foreground(completedFileColor)
	inactiveFileStyle           = lipgloss.NewStyle().Foreground(inactiveFileColor)
	highlightedTextStyle        = lipgloss.NewStyle().Foreground(highlightedTextColor).Bold(true)
	renderedStatusTextStyle     = lipgloss.NewStyle().Width(15).Align(lipgloss.Right)
	iconStyle                   = lipgloss.NewStyle().MarginRight(1)
	kanbanTaskTitleStyle        = lipgloss.NewStyle().Bold(true).Foreground(summaryTitleColor)
	suggestionTitleStyle        = lipgloss.NewStyle().Foreground(suggestionTitleColor).Bold(true)
	suggestionTextStyle         = lipgloss.NewStyle().Foreground(suggestionTextColor)
	selectedSuggestionTextStyle = lipgloss.NewStyle().Foreground(selectedSuggestionTextColor)
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

func summaryContainerStyle(width, height int) lipgloss.Style {
	return contentContainerStyle(width, height)
}

func summaryContainerStyleNoBorder(width, height int) lipgloss.Style {
	return contentContainerStyleNoBorder(width, height)
}

func contentContainerStyle(width, height int) lipgloss.Style {
	key := fmt.Sprintf("container-%d-%d", width, height)

	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.containerStyles[key]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.containerStyles[key]; ok {
		return style
	}

	style := lipgloss.NewStyle().Width(width).Height(height).Border(lipgloss.NormalBorder()).MarginLeft(2)
	globalStyleCache.containerStyles[key] = style
	return style
}

func contentContainerStyleNoBorder(width, height int) lipgloss.Style {
	key := fmt.Sprintf("container-noborder-%d-%d", width, height)

	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.containerStyles[key]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.containerStyles[key]; ok {
		return style
	}

	style := lipgloss.NewStyle().Width(width).Height(height).MarginLeft(2)
	globalStyleCache.containerStyles[key] = style
	return style
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
		style = style.BorderForeground(lipgloss.Color("63"))
	}

	return style
}

func filesItemDetailsContainerStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width).Height(height).MarginLeft(2).Border(lipgloss.NormalBorder())
}

func progressTextStyle(style lipgloss.Style) lipgloss.Style {
	return style.Width(35).Align(lipgloss.Right)
}

func taskStyle(width int) lipgloss.Style {
	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.taskStyles[width]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.taskStyles[width]; ok {
		return style
	}

	style := lipgloss.NewStyle().MarginLeft(2).Width(width)
	globalStyleCache.taskStyles[width] = style
	return style
}

func kanbanTaskStyle(width int) lipgloss.Style {
	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.kanbanTaskStyles[width]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.kanbanTaskStyles[width]; ok {
		return style
	}

	style := lipgloss.NewStyle().Width(width).MarginBottom(1)
	globalStyleCache.kanbanTaskStyles[width] = style
	return style
}

func highlightedKanbanTaskStyle(width int) lipgloss.Style {
	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.highlightedKanbanTask[width]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.highlightedKanbanTask[width]; ok {
		return style
	}

	style := lipgloss.NewStyle().Width(width).Background(highlightedTaskBackgroundColor).MarginBottom(1)
	globalStyleCache.highlightedKanbanTask[width] = style
	return style
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
	return taskFileTitleStyle.Foreground(inactiveFileColor)
}

func completedTitleStyle() lipgloss.Style {
	return taskFileTitleStyle.Foreground(completedFileColor)
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

func boardContainerStyle(width int, height int, selected bool) lipgloss.Style {
	key := fmt.Sprintf("board-%d-%d-%t", width, height, selected)

	globalStyleCache.mu.RLock()
	if style, ok := globalStyleCache.boardStyles[key]; ok {
		globalStyleCache.mu.RUnlock()
		return style
	}
	globalStyleCache.mu.RUnlock()

	globalStyleCache.mu.Lock()
	defer globalStyleCache.mu.Unlock()

	// Double-check after acquiring write lock
	if style, ok := globalStyleCache.boardStyles[key]; ok {
		return style
	}

	style := lipgloss.NewStyle().Width(width).Height(height).Border(lipgloss.NormalBorder())

	if selected {
		style = style.BorderForeground(lipgloss.Color("63"))
	}

	globalStyleCache.boardStyles[key] = style
	return style
}

func filterInputStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color(color))
}

func taskOptionSuggestionsContainerStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(2)
}
