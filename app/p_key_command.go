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
			m.TaskManager.UpdateTaskToPriority(m.FileManager, selectedTask)
			m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
		} else {
			log.Info("Removing priority marker from task")
			m.TaskManager.UpdateTaskToUnpriority(m.FileManager, selectedTask)
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
