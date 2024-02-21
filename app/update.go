package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		} else {
			keyCommandFactory := KeyCommandFactory{}
			keyCommand := keyCommandFactory.CreateKeyCommand(key)

			if keyCommand != nil {
				err := keyCommand.Execute(&m)
				if err != nil {
					return m, tea.Quit
				}
			}
		}
	case tea.WindowSizeMsg:
		m.Viewport.Width = msg.Width - 60
		m.Viewport.Height = msg.Height - 20
		m.Width = msg.Width
		m.Height = msg.Height
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
