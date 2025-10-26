package app

import "github.com/charmbracelet/log"

type UppercaseSKeyCommand struct{}

func (j UppercaseSKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		if m.TaskManager.SelectedTask.Started {
			if err := m.TaskManager.UpdateTaskToScheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 1
			m.ViewManager.IsKanbanTaskUpdated = true
		} else {
			if err := m.TaskManager.UpdateTaskToUnscheduled(&m.FileManager, m.TaskManager.SelectedTask); err != nil {
				m.Errors = append(m.Errors, err.Error())
			}
			m.ViewManager.KanbanListCursor = 0
			m.ViewManager.IsKanbanTaskUpdated = true
			log.Info("Updating task to unscheduled ", m.ViewManager.IsKanbanTaskUpdated)
		}

		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
	}

	return nil
}

func (j UppercaseSKeyCommand) HelpText() string {
	return "UppercaseSKeyCommand help text"
}

func (j UppercaseSKeyCommand) AllowedStates() []string {
	return []string{}
}
