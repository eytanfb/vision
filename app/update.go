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
			factory := NewKeyCommandFactory()
			if key == "esc" {
				cmdResult := factory.CreateKeyCommand("esc").Execute(m)
				cmds = append(cmds, cmdResult)
			} else if key == "enter" {
				cmdResult := factory.CreateKeyCommand("enter").Execute(m)
				cmds = append(cmds, cmdResult)
			}

			if m.IsFilterView() {
				m.FilterInput, cmd = m.FilterInput.Update(msg)
				m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
				cmds = append(cmds, cmd)
			} else {
				if key == "[" {
					cmdResult := factory.CreateKeyCommand("[").Execute(m)
					cmds = append(cmds, cmdResult)
				} else if key == "]" {
					cmdResult := factory.CreateKeyCommand("]").Execute(m)
					cmds = append(cmds, cmdResult)
				}

				m.NewTaskInput, cmd = m.NewTaskInput.Update(msg)
				cmds = append(cmds, cmd)
			}

			if key == "tab" {
				cmdResult := factory.CreateKeyCommand("tab").Execute(m)
				cmds = append(cmds, cmdResult)
			} else if key == "shift+tab" {
				cmdResult := factory.CreateKeyCommand("shift+tab").Execute(m)
				cmds = append(cmds, cmdResult)
			}
		} else {
			keyCommandFactory := NewKeyCommandFactory()
			keyCommand := keyCommandFactory.CreateKeyCommand(key)
			cmdResult := keyCommand.Execute(m)
			cmds = append(cmds, cmdResult)
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
