package app

type UppercaseAKeyCommand struct{}

func (j UppercaseAKeyCommand) Execute(m *Model) error {
	if m.FileManager.SelectedFile.Name == "" {
		m.Errors = append(m.Errors, "No file selected")
		return nil
	}
	m.ViewManager.IsAddSubTaskView = true
	m.NewTaskInput.Reset()

	prompt := m.FileManager.SelectedFile.FileNameWithoutExtension() + "\n"

	m.NewTaskInput.Prompt = prompt
	m.NewTaskInput.Placeholder = "> Add a subtask..."
	m.NewTaskInput.Focus()

	return nil
}

func (j UppercaseAKeyCommand) HelpText() string {
	return "UppercaseAKeyCommand help text"
}

func (j UppercaseAKeyCommand) AllowedStates() []string {
	return []string{}
}
