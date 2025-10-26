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
	} else if m.IsCategoryView() {
		nc.moveDownInCategoryView(m)
	} else if m.IsDetailsView() {
		nc.moveDownInDetailsView(m)
	}
	return nil
}

func (nc NavigationCommands) moveDownInCategoryView(m *Model) {
	if !m.ViewManager.HideSidebar {
		m.GoToNextCategory()
	} else {
		m.GoToNextKanbanTask()
	}
}

func (nc NavigationCommands) moveDownInDetailsView(m *Model) {
	if m.IsItemDetailsFocus() {
		nc.moveDownInItemDetails(m)
	} else {
		m.GoToNextFile()
		m.Viewport.GotoTop()
	}
}

func (nc NavigationCommands) moveDownInItemDetails(m *Model) {
	if m.IsTaskDetailsFocus() {
		m.GoToNextTask()
	} else {
		m.Viewport.LineDown(10)
	}
}

// MoveUp handles k key - move cursor up
func (nc NavigationCommands) MoveUp(m *Model) tea.Cmd {
	if m.IsDetailsView() {
		if m.IsItemDetailsFocus() {
			if m.IsTaskDetailsFocus() {
				m.GoToPreviousTask()
			} else {
				m.Viewport.LineUp(10)
			}
		} else {
			m.GoToPreviousFile()
			m.Viewport.GotoTop()
		}
	} else if m.IsCategoryView() {
		if !m.ViewManager.HideSidebar {
			m.GoToPreviousCategory()
		} else {
			m.GoToPreviousKanbanTask()
		}
	} else if m.IsCompanyView() {
		m.GoToPreviousCompany()
	}
	return nil
}

// MoveLeft handles h key - move left or go back
func (nc NavigationCommands) MoveLeft(m *Model) tea.Cmd {
	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			m.GoToPreviousKanbanList()
			m.ViewManager.KanbanTaskCursor = 0
			return nil
		}
	}

	// Delegate to Esc behavior for going back
	return nc.goBack(m)
}

// MoveRight handles l key - move right or select
func (nc NavigationCommands) MoveRight(m *Model) tea.Cmd {
	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			m.GoToNextKanbanList()
			m.ViewManager.KanbanTaskCursor = 0
			return nil
		}
	}

	// Delegate to Enter behavior for selection
	return nc.selectItem(m)
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

func (nc NavigationCommands) goBack(m *Model) tea.Cmd {
	// Esc key behavior - go to previous view
	if m.IsDetailsView() {
		m.GoToPreviousView()
	}
	return nil
}

func (nc NavigationCommands) selectItem(m *Model) tea.Cmd {
	// Enter key behavior - select current item
	if !m.ViewManager.HideSidebar {
		m.Select()
	}
	return nil
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
