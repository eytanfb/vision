package app

type EscKeyCommand struct{}

func (j EscKeyCommand) Execute(m *Model) error {
	if m.IsAddTaskView() {
		m.ViewManager.IsAddTaskView = false
		m.NewTaskInput.Blur()
	}
	m.GoToPreviousView()

	return nil
}

func (j EscKeyCommand) HelpText() string {
	return "EscKeyCommand help text"
}

func (j EscKeyCommand) AllowedStates() []string {
	return []string{}
}
