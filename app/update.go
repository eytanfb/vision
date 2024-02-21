package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		} else {
			keyCommandFactory := KeyCommandFactory{}
			keyCommand := keyCommandFactory.CreateKeyCommand(key)

			err := keyCommand.Execute(m)
			if err != nil {
				return m, tea.Quit
			}

			m.Viewport, cmd = m.Viewport.Update(msg)
		}
	case tea.WindowSizeMsg:
		if !m.Ready {
			m.Viewport = viewport.New(msg.Width-60, msg.Height-20)
			m.Viewport.YPosition = 20
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width - 60
			m.Viewport.Height = msg.Height - 20
		}
		m.Width = msg.Width
		m.Height = msg.Height
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
