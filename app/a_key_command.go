package app

type AKeyCommand struct{}

func (j AKeyCommand) Execute(m *Model) error {
	m.ViewManager.IsAddTaskView = true
	m.NewTaskInput.Focus()

	return nil
}

func (j AKeyCommand) HelpText() string {
	return "AKeyCommand help text"
}

func (j AKeyCommand) AllowedStates() []string {
	return []string{}
}
