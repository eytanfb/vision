package app

type EnterKeyCommand struct{}

func (j EnterKeyCommand) Execute(m *Model) error {
	if m.IsAddTaskView() {
		company := m.GetCurrentCompanyName()
		input := m.NewTaskInput.Value()

		if m.ViewManager.IsTaskDetailsFocus() {
			//m.FileManager.CreateTask(company, input)
			EscKeyCommand{}.Execute(m)
		} else {
			m.FileManager.CreateTask(company, input)
			EscKeyCommand{}.Execute(m)
		}

		return nil
	} else if m.IsFilterView() {
		m.ViewManager.IsFilterView = false
		m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
		m.TaskManager.TasksCursor = 0

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
