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
		}
	case tea.WindowSizeMsg:
		m.ViewManager.SetWidth(msg.Width)
		m.ViewManager.SetHeight(msg.Height)

		if !m.ViewManager.Ready {
			m.Viewport = viewport.New(msg.Width-m.ViewManager.SidebarWidth, msg.Height-28)
			m.ViewManager.Ready = true
		} else {
			m.Viewport.Width = msg.Width - m.ViewManager.SidebarWidth
			m.Viewport.Height = msg.Height - 28
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
