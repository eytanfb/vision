package app

type EnterKeyCommand struct{}

func (j EnterKeyCommand) Execute(m *Model) error {
	if m.IsAddTaskView() {
		company := m.GetCurrentCompanyName()
		taskName := m.NewTaskInput.Value()
		m.FileManager.CreateTask(company, taskName)
		EscKeyCommand{}.Execute(m)

		return nil
	}
	m.Select()

	return nil
}

func (j EnterKeyCommand) HelpText() string {
	return "EnterKeyCommand help text"
}

func (j EnterKeyCommand) AllowedStates() []string {
	return []string{}
}
