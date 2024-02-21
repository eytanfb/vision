package app

type SKeyCommand struct{}

func (j SKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "categories" {
		m.SelectedCategory = "standups"
		m.CurrentView = "details"
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	}

	return nil
}

func (j SKeyCommand) HelpText() string {
	return "SKeyCommand help text"
}

func (j SKeyCommand) AllowedStates() []string {
	return []string{}
}
