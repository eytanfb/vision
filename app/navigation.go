package app

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// NavigationCommands handles all cursor movement and navigation
type NavigationCommands struct{}

// HandleKey routes navigation keys to their appropriate handlers
func (nc NavigationCommands) HandleKey(key string, m *Model) tea.Cmd {
	switch key {
	case "j":
		return nc.MoveDown(m)
	case "k":
		return nc.MoveUp(m)
	case "h":
		return nc.MoveLeft(m)
	case "l":
		return nc.MoveRight(m)
	case "g":
		return nc.OpenGitHubDash(m)
	case "tab":
		return nc.NextSuggestion(m)
	case "shift+tab":
		return nc.PreviousSuggestion(m)
	}
	return nil
}

// MoveDown handles j key - move cursor down
func (nc NavigationCommands) MoveDown(m *Model) tea.Cmd {
	if m.IsSuggestionsActive() {
		m.ViewManager.NextSuggestion(&m.FileManager)
		return nil
	}

	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			return m.navigateKanbanCmd("down")
		}
		nc.moveDownInCategoryView(m)
		return nil
	}

	if m.IsDetailsView() {
		return nc.moveDownInDetailsView(m)
	}

	return nil
}

func (nc NavigationCommands) moveDownInCategoryView(m *Model) {
	m.GoToNextCategory()
}

func (nc NavigationCommands) moveDownInDetailsView(m *Model) tea.Cmd {
	if m.IsItemDetailsFocus() {
		if m.IsTaskDetailsFocus() {
			m.GoToNextTask()
			nc.updateSelectedTask(m)
			return nil
		}
		m.Viewport.LineDown(ViewportLineScroll)
		return nil
	}
	// File navigation - use message pattern
	return m.navigateFileCmd("down")
}

// MoveUp handles k key - move cursor up
func (nc NavigationCommands) MoveUp(m *Model) tea.Cmd {
	if m.IsDetailsView() {
		return nc.moveUpInDetailsView(m)
	}

	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			return m.navigateKanbanCmd("up")
		}
		m.GoToPreviousCategory()
		return nil
	}

	if m.IsCompanyView() {
		m.GoToPreviousCompany()
	}
	return nil
}

func (nc NavigationCommands) moveUpInDetailsView(m *Model) tea.Cmd {
	if m.IsItemDetailsFocus() {
		if m.IsTaskDetailsFocus() {
			m.GoToPreviousTask()
			nc.updateSelectedTask(m)
			return nil
		}
		m.Viewport.LineUp(ViewportLineScroll)
		return nil
	}
	// File navigation - use message pattern
	return m.navigateFileCmd("up")
}

// MoveLeft handles h key - move left or go back
func (nc NavigationCommands) MoveLeft(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		return m.navigateKanbanCmd("left")
	}
	nc.goBack(m)
	return nil
}

// MoveRight handles l key - move right or select
func (nc NavigationCommands) MoveRight(m *Model) tea.Cmd {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		return m.navigateKanbanCmd("right")
	}
	nc.goForward(m)
	return nil
}

// OpenGitHubDash handles g key - opens GitHub dashboard using non-blocking tea.ExecProcess
func (nc NavigationCommands) OpenGitHubDash(m *Model) tea.Cmd {
	c := exec.Command("gh", "dash")

	// Use tea.ExecProcess for non-blocking execution
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return ErrorOccurredMsg{
				Err:     err,
				Context: "opening GitHub dashboard",
			}
		}
		// Dashboard closed successfully
		return nil
	})
}

// NextSuggestion handles tab key - move to next suggestion
func (nc NavigationCommands) NextSuggestion(m *Model) tea.Cmd {
	if m.ViewManager.SuggestionsListsCursor == -1 {
		m.ViewManager.SuggestionsListsCursor = 0
	}
	m.ViewManager.NextSuggestion(&m.FileManager)
	return nil
}

// PreviousSuggestion handles shift+tab key - move to previous suggestion
func (nc NavigationCommands) PreviousSuggestion(m *Model) tea.Cmd {
	log.Info("PreviousSuggestion")
	m.ViewManager.PreviousSuggestion(&m.FileManager)
	return nil
}

// Helper methods for h and l keys that delegate to other command behavior

func (nc NavigationCommands) goBack(m *Model) {
	// Esc key behavior - go to previous view
	if m.IsDetailsView() {
		m.GoToPreviousView()
	}
}

func (nc NavigationCommands) goForward(m *Model) {
	// Enter key behavior - select current item
	if !m.ViewManager.HideSidebar {
		m.Select()
	}
}

// Command implementations for registry

type JKeyCommand struct{}

func (cmd JKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.MoveDown(m)
}

func (cmd JKeyCommand) Description() string {
	return "Move cursor down"
}

func (cmd JKeyCommand) Contexts() []string {
	return []string{}
}

type KKeyCommand struct{}

func (cmd KKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.MoveUp(m)
}

func (cmd KKeyCommand) Description() string {
	return "Move cursor up"
}

func (cmd KKeyCommand) Contexts() []string {
	return []string{}
}

type HKeyCommand struct{}

func (cmd HKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.MoveLeft(m)
}

func (cmd HKeyCommand) Description() string {
	return "Move left or go back"
}

func (cmd HKeyCommand) Contexts() []string {
	return []string{}
}

type LKeyCommand struct{}

func (cmd LKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.MoveRight(m)
}

func (cmd LKeyCommand) Description() string {
	return "Move right or select"
}

func (cmd LKeyCommand) Contexts() []string {
	return []string{}
}

type GKeyCommand struct{}

func (cmd GKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.OpenGitHubDash(m)
}

func (cmd GKeyCommand) Description() string {
	return "Open GitHub dashboard"
}

func (cmd GKeyCommand) Contexts() []string {
	return []string{}
}

type TabKeyCommand struct{}

func (cmd TabKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.NextSuggestion(m)
}

func (cmd TabKeyCommand) Description() string {
	return "Next suggestion"
}

func (cmd TabKeyCommand) Contexts() []string {
	return []string{}
}

type ShiftTabKeyCommand struct{}

func (cmd ShiftTabKeyCommand) Execute(m *Model) tea.Cmd {
	return NavigationCommands{}.PreviousSuggestion(m)
}

func (cmd ShiftTabKeyCommand) Description() string {
	return "Previous suggestion"
}

func (cmd ShiftTabKeyCommand) Contexts() []string {
	return []string{}
}

// Helper functions to update selected items after navigation
// These are called from Update() handlers, not from View()

func (nc NavigationCommands) updateSelectedTask(m *Model) {
	tasks := m.TaskManager.TaskCollection.GetTasks(m.FileManager.currentFileName())
	if m.TaskManager.TasksCursor >= 0 && m.TaskManager.TasksCursor < len(tasks) {
		task := tasks[m.TaskManager.TasksCursor]
		m.TaskManager.SelectedTask = task
		m.FileManager.SelectFile(task.FileName)
	}
}

func (nc NavigationCommands) updateSelectedTaskInKanban(m *Model) {
	// Get ALL tasks from ALL files (not just current file)
	allTasksMap := m.TaskManager.TaskCollection.TasksByFile

	// Use the same sorting as the view rendering (viewSort) to ensure cursor position
	// maps to the same task visually and in selection logic
	keys := m.getSortedKanbanKeys(allTasksMap)

	// Build flat list of tasks matching the current kanban list
	// This must match the same order as makeKanbanLists in view_builder.go
	var filteredTasks []Task

	for _, key := range keys {
		tasks := allTasksMap[key]
		for _, task := range tasks {
			// Match the same filtering logic as makeKanbanLists
			if m.ViewManager.KanbanListCursor == 0 {
				// Inactive: future scheduled or unscheduled (no scheduled date)
				if task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) || task.ScheduledDate == "" {
					filteredTasks = append(filteredTasks, task)
				}
			} else if m.ViewManager.KanbanListCursor == 1 {
				// Active: started or scheduled (not completed, not future)
				if !task.Completed && !task.IsScheduledForFuture(m.TaskManager.DailySummaryDate) && (task.Started || task.Scheduled) {
					filteredTasks = append(filteredTasks, task)
				}
			} else if m.ViewManager.KanbanListCursor == 2 {
				// Completed
				if task.Completed {
					filteredTasks = append(filteredTasks, task)
				}
			}
		}
	}

	// Clamp cursor to valid bounds if task list changed (e.g., after date filter change)
	if len(filteredTasks) > 0 {
		if m.ViewManager.KanbanTaskCursor >= len(filteredTasks) {
			m.ViewManager.KanbanTaskCursor = len(filteredTasks) - 1
		}
		if m.ViewManager.KanbanTaskCursor < 0 {
			m.ViewManager.KanbanTaskCursor = 0
		}
		m.TaskManager.SelectedTask = filteredTasks[m.ViewManager.KanbanTaskCursor]
	} else {
		// No tasks in this list - clear selection
		m.TaskManager.SelectedTask = Task{}
	}
}
