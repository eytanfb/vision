package app

type UppercaseSKeyCommand struct{}

func (j UppercaseSKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		if m.TaskManager.SelectedTask.Started {
			m.TaskManager.UpdateTaskToScheduled(m.FileManager, m.TaskManager.SelectedTask)
			m.ViewManager.KanbanListCursor = 1
		} else {
			m.TaskManager.UpdateTaskToUnscheduled(m.FileManager, m.TaskManager.SelectedTask)
			m.ViewManager.KanbanListCursor = 0
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
