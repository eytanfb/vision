package app

import (
	"net/url"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// FileOperations handles file-related commands
type FileOperations struct{}

// HandleKey routes file operation keys to their appropriate handlers
func (fo FileOperations) HandleKey(key string, m *Model) tea.Cmd {
	switch key {
	case "e":
		return fo.OpenInVim(m)
	case "o":
		return fo.OpenInObsidian(m)
	case "n":
		return fo.NextCompany(m)
	case "f":
		return fo.ToggleSidebar(m)
	}
	return nil
}

// OpenInVim opens the current file in Vim using non-blocking tea.ExecProcess
func (fo FileOperations) OpenInVim(m *Model) tea.Cmd {
	if m.IsDetailsView() || m.IsKanbanView() {
		log.Info("Opening file in vim", m.FileManager.SelectedFile.Name)
		filePath := m.FileManager.SelectedFile.FullPath

		c := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)

		// Use tea.ExecProcess for non-blocking execution
		return tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				return EditorClosedMsg{Err: err}
			}
			return EditorClosedMsg{}
		})
	}
	return nil
}

// OpenInObsidian opens the current file in Obsidian using non-blocking tea.ExecProcess
func (fo FileOperations) OpenInObsidian(m *Model) tea.Cmd {
	if m.IsDetailsView() || m.IsKanbanView() {
		filePath := m.FileManager.SelectedFile.FullPath
		homeDir, _ := os.UserHomeDir()
		notesPath := homeDir + "/Notes"
		obsidianPath := constructObsidianURL(filePath, notesPath)

		c := exec.Command("open", "-a", "Obsidian", obsidianPath)

		// Use tea.ExecProcess for non-blocking execution
		return tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				return ErrorOccurredMsg{
					Err:     err,
					Context: "opening in Obsidian",
				}
			}
			// Obsidian opened successfully, no need to reload file
			return nil
		})
	}
	return nil
}

// NextCompany switches to the next company
func (fo FileOperations) NextCompany(m *Model) tea.Cmd {
	m.GoToNextCompany()
	m.GotoCategoryView()
	return nil
}

// ToggleSidebar toggles the visibility of the sidebar
func (fo FileOperations) ToggleSidebar(m *Model) tea.Cmd {
	m.ViewManager.ToggleHideSidebar()
	return nil
}

// Helper function for constructing Obsidian URLs
func constructObsidianURL(fullPath, notesPath string) string {
	relativePath := strings.Replace(fullPath, notesPath, "", 1)
	urlEncodedPath := url.PathEscape(relativePath)
	obsidianPath := "obsidian://open?vault=Notes&file=" + urlEncodedPath
	return obsidianPath
}

// Command implementations for registry

type EKeyCommand struct{}

func (cmd EKeyCommand) Execute(m *Model) tea.Cmd {
	return FileOperations{}.OpenInVim(m)
}

func (cmd EKeyCommand) Description() string {
	return "Open current file in Vim"
}

func (cmd EKeyCommand) Contexts() []string {
	return []string{"details", "kanban"}
}

type OKeyCommand struct{}

func (cmd OKeyCommand) Execute(m *Model) tea.Cmd {
	return FileOperations{}.OpenInObsidian(m)
}

func (cmd OKeyCommand) Description() string {
	return "Open current file in Obsidian"
}

func (cmd OKeyCommand) Contexts() []string {
	return []string{"details", "kanban"}
}

type NKeyCommand struct{}

func (cmd NKeyCommand) Execute(m *Model) tea.Cmd {
	return FileOperations{}.NextCompany(m)
}

func (cmd NKeyCommand) Description() string {
	return "Switch to next company"
}

func (cmd NKeyCommand) Contexts() []string {
	return []string{}
}

type FKeyCommand struct{}

func (cmd FKeyCommand) Execute(m *Model) tea.Cmd {
	return FileOperations{}.ToggleSidebar(m)
}

func (cmd FKeyCommand) Description() string {
	return "Toggle sidebar visibility"
}

func (cmd FKeyCommand) Contexts() []string {
	return []string{}
}
