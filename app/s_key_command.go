package app

import "github.com/charmbracelet/log"

type SKeyCommand struct{}

func (j SKeyCommand) Execute(m *Model) error {
	if m.ViewManager.HideSidebar {
		log.Info("SKeyCommand: Show sidebar")
		if m.TaskManager.SelectedTask.Scheduled {
			log.Info("SKeyCommand: Update task to started")
			if err := m.TaskManager.UpdateTaskToStarted(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
		} else {
			log.Info("SKeyCommand: Update task to scheduled")
			if err := m.TaskManager.UpdateTaskToScheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
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
