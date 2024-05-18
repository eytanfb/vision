package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Errors = []string{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		log.Info("Key pressed: ", key)
		if key == "ctrl+c" {
			return m, tea.Quit
		} else if key == "q" {
			if m.IsAddTaskView() {
				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				return m, cmd
			}
			return m, tea.Quit
		} else if m.IsAddTaskView() || m.IsFilterView() || m.IsAddSubTaskView() {
			if key == "esc" {
				KeyCommandFactory{}.CreateKeyCommand("esc").Execute(m)
			} else if key == "enter" {
				KeyCommandFactory{}.CreateKeyCommand("enter").Execute(m)
			}

			if m.IsFilterView() {
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
			} else {
				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
			}
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
			m.Viewport = viewport.New(m.ViewManager.DetailsViewWidth, m.ViewManager.DetailsViewHeight)
			m.ViewManager.Ready = true
		} else {
			m.Viewport.Width = m.ViewManager.DetailsViewWidth
			m.Viewport.Height = m.ViewManager.DetailsViewHeight

			if m.ViewManager.DetailsViewHeight > 65 {
				m.ViewManager.KanbanViewLineDownFactor = 5
			} else if m.ViewManager.DetailsViewHeight > 45 {
				m.ViewManager.KanbanViewLineDownFactor = 10
			} else {
				m.ViewManager.KanbanViewLineDownFactor = 15
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
