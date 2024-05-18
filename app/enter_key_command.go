package app

type EnterKeyCommand struct{}

func (j EnterKeyCommand) Execute(m *Model) error {
	if m.IsAddTaskView() {
		company := m.GetCurrentCompanyName()
		input := m.NewTaskInput.Value()

		m.FileManager.CreateTask(company, input)
		return EscKeyCommand{}.Execute(m)
	} else if m.IsAddSubTaskView() {
		company := m.GetCurrentCompanyName()
		input := m.NewTaskInput.Value()
		selectedFile := m.FileManager.SelectedFile

		m.FileManager.CreateSubTask(company, selectedFile, input)
		return EscKeyCommand{}.Execute(m)
	} else if m.IsFilterView() {
		m.ViewManager.IsFilterView = false
		m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
		m.TaskManager.TasksCursor = 0
		m.FilterInput.Blur()

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
