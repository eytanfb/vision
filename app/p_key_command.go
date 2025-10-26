package app

import (
	"strings"

	"github.com/charmbracelet/log"
)

type PKeyCommand struct{}

func (j PKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		selectedTask := m.TaskManager.SelectedTask

		if !strings.Contains(selectedTask.Text, "🔺") {
			log.Info("Adding priority marker to task")
			if err := m.TaskManager.UpdateTaskToPriority(&m.FileManager, selectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		} else {
			log.Info("Removing priority marker from task")
			if err := m.TaskManager.UpdateTaskToUnpriority(&m.FileManager, selectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		}
	}

	return nil
}

func (j PKeyCommand) HelpText() string {
	return "Add priority marker to task"
}

func (j PKeyCommand) AllowedStates() []string {
	return []string{}
}
