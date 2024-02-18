package main

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
			m.cursor++
			if m.currentView == "companies" {
				if m.cursor >= len(m.companies) {
					m.cursor = len(m.companies) - 1
				}
			} else if m.currentView == "categories" {
				if m.cursor >= len(m.categories) {
					m.cursor = len(m.categories) - 1
				}
			} else if m.currentView == "details" {
				if m.itemDetailsFocus {
					m.viewport.LineDown(10)
				} else {
					if m.cursor >= len(m.files) {
						m.cursor = len(m.files) - 1
					}
				}
			}

		case "k":
			if m.currentView == "details" {
				if m.itemDetailsFocus {
					m.viewport.LineUp(10)
				}
			} else {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = 0
				}
			}
		case "enter":
			if m.currentView == "companies" {
				for _, company := range m.companies {
					if company.DisplayName == m.companies[m.cursor].DisplayName {
						m.selectedCompany = company
						break
					}
				}
				m.currentView = "categories"
				m.cursor = 0
			} else if m.currentView == "categories" {
				m.selectedCategory = m.categories[m.cursor]
				m.currentView = "details"
				m.cursor = 0
				m.files = m.FetchFiles()
			}
		case "esc":
			if m.currentView == "categories" {
				m.currentView = "companies"
				m.cursor = 0
			} else if m.currentView == "details" {
				m.currentView = "categories"
				m.cursor = 0
			}
		case "l":
			if m.currentView == "details" {
				m.itemDetailsFocus = true
			}
		case "h":
			if m.currentView == "details" {
				m.itemDetailsFocus = false
			}
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 40
		m.viewport.Height = msg.Height - 40
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func getCurrentFilePath(m Model) string {
	return ""
}
