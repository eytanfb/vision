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
			} else if m.IsAddSubTaskView() {
				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				return m, cmd
			} else if m.IsFilterView() {
				m.FilterInput, cmd = m.FilterInput.Update(msg)
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
				if key == "[" {
					KeyCommandFactory{}.CreateKeyCommand("[").Execute(m)
				} else if key == "]" {
					KeyCommandFactory{}.CreateKeyCommand("]").Execute(m)
				}

				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
			}

			if key == "tab" {
				KeyCommandFactory{}.CreateKeyCommand("tab").Execute(m)
			} else if key == "shift+tab" {
				ShiftTabKeyCommand{}.Execute(m)
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

			// Adjust scroll speed based on terminal height
			switch {
			case m.ViewManager.DetailsViewHeight > largeTerminalHeight:
				m.ViewManager.KanbanViewLineDownFactor = largeTerminalScrollFactor
			case m.ViewManager.DetailsViewHeight > mediumTerminalHeight:
				m.ViewManager.KanbanViewLineDownFactor = mediumTerminalScrollFactor
			default:
				m.ViewManager.KanbanViewLineDownFactor = smallTerminalScrollFactor
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
