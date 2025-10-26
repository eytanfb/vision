package app

import tea "github.com/charmbracelet/bubbletea"

type KeyCommand interface {
	Execute(m *Model) tea.Cmd
	HelpText() string
	AllowedStates() []string
}
