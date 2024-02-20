package app

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e":
			// Assuming `getCurrentFilePath` is a function you'll implement to get the selected file path
			filePath := getCurrentFilePath(m)
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim" // Default to vim if $EDITOR is not set
			}
			cmd := exec.Command(editor, filePath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			go func() {
				_ = cmd.Run() // Consider handling this error
				// Refresh or re-render logic here if necessary
			}()
			return m, nil
		case "j":
			m.Cursor++
			if m.CurrentView == "Companies" {
				if m.Cursor >= len(m.Companies) {
					m.Cursor = len(m.Companies) - 1
				}
			} else if m.CurrentView == "categories" {
				if m.Cursor >= len(m.Categories) {
					m.Cursor = len(m.Categories) - 1
				}
			} else if m.CurrentView == "details" {
				if m.ItemDetailsFocus {
					m.Viewport.LineDown(10)
				} else {
					if m.Cursor >= len(m.Files) {
						m.Cursor = len(m.Files) - 1
					}
				}
			}

		case "k":
			if m.CurrentView == "details" {
				if m.ItemDetailsFocus {
					m.Viewport.LineUp(10)
				} else {
					m.Cursor--
					if m.Cursor < 0 {
						m.Cursor = 0
					}
				}
			} else {
				m.Cursor--
				if m.Cursor < 0 {
					m.Cursor = 0
				}
			}
		case "enter":
			if m.CurrentView == "companies" {
				for _, company := range m.Companies {
					if company.DisplayName == m.Companies[m.Cursor].DisplayName {
						m.SelectedCompany = company
						break
					}
				}
				m.CurrentView = "categories"
				m.Cursor = 0
			} else if m.CurrentView == "categories" {
				m.SelectedCategory = m.Categories[m.Cursor]
				m.CurrentView = "details"
				m.Cursor = 0
				m.Files = m.FetchFiles()
			}
		case "esc":
			if m.CurrentView == "categories" {
				m.CurrentView = "companies"
				m.Cursor = 0
			} else if m.CurrentView == "details" {
				m.CurrentView = "categories"
				m.Cursor = 0
			}
		case "l":
			if m.CurrentView == "details" {
				m.ItemDetailsFocus = true
			}
		case "h":
			if m.CurrentView == "details" {
				m.ItemDetailsFocus = false
			}
		case "t":
			if m.CurrentView == "categories" {
				m.SelectedCategory = "tasks"
				m.CurrentView = "details"
				m.Cursor = 0
				m.Files = m.FetchFiles()
			}
		case "m":
			if m.CurrentView == "categories" {
				m.SelectedCategory = "meetings"
				m.CurrentView = "details"
				m.Cursor = 0
				m.Files = m.FetchFiles()
			}
		case "s":
			if m.CurrentView == "categories" {
				m.SelectedCategory = "standups"
				m.CurrentView = "details"
				m.Cursor = 0
				m.Files = m.FetchFiles()
			}
		}
	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width - 40
		m.Viewport.Height = msg.Height - 40
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func getCurrentFilePath(m Model) string {
	return ""
}
