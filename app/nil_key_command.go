package app

import tea "github.com/charmbracelet/bubbletea"

type NilCommand struct{}

func (j NilCommand) Execute(m *Model) tea.Cmd {
	return nil
}

func (j NilCommand) Description() string {
	return ""
}

func (j NilCommand) Contexts() []string {
	return []string{}
}
