package app

import "github.com/charmbracelet/log"

type SKeyCommand struct{}

func (j SKeyCommand) Execute(m *Model) error {
	if m.ViewManager.HideSidebar {
		log.Info("SKeyCommand: Show sidebar")
		if m.TaskManager.SelectedTask.Scheduled {
			log.Info("SKeyCommand: Update task to started")
			m.TaskManager.UpdateTaskToStarted(m.FileManager, m.TaskManager.SelectedTask)
		} else {
			log.Info("SKeyCommand: Update task to scheduled")
			m.TaskManager.UpdateTaskToScheduled(m.FileManager, m.TaskManager.SelectedTask)
		}

		m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)
	}

	return nil
}

func (j SKeyCommand) HelpText() string {
	return "SKeyCommand help text"
}

func (j SKeyCommand) AllowedStates() []string {
	return []string{}
}
