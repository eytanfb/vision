package command

import tea "github.com/charmbracelet/bubbletea"

func init() {
	RegisterCommand("exit_navigation", func() Command {
		return &ExitNavigationCommand{}
	})
}

type ExitNavigationCommand struct {
	Key string
}

func (c *ExitNavigationCommand) Execute(m *Model) {
	return m, tea.Quit
}
