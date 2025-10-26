package app

import tea "github.com/charmbracelet/bubbletea"

type NilKeyCommand struct{}

func (j NilKeyCommand) Execute(m *Model) tea.Cmd {
	return nil
}

func (j NilKeyCommand) HelpText() string {
	return ""
}

func (j NilKeyCommand) AllowedStates() []string {
	return []string{}
}
