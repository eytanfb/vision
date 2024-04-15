package app

type DKeyCommand struct{}

func (j DKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() && m.ViewManager.HideSidebar {
		m.TaskManager.UpdateTaskToCompleted(m.FileManager, m.TaskManager.SelectedTask)
		m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)
	}

	return nil
}

func (j DKeyCommand) HelpText() string {
	return "DKeyCommand help text"
}

func (j DKeyCommand) AllowedStates() []string {
	return []string{}
}
