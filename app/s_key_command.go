package app

type SKeyCommand struct{}

func (j SKeyCommand) Execute(m *Model) error {
	if m.TaskManager.SelectedTask.Scheduled {
		m.TaskManager.UpdateTaskToStarted(m.FileManager, m.TaskManager.SelectedTask)
	} else {
		m.TaskManager.UpdateTaskToScheduled(m.FileManager, m.TaskManager.SelectedTask)
	}

	m.FileManager.FetchTasks(&m.DirectoryManager, &m.TaskManager)

	return nil
}

func (j SKeyCommand) HelpText() string {
	return "SKeyCommand help text"
}

func (j SKeyCommand) AllowedStates() []string {
	return []string{}
}
