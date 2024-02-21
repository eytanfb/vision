package app

type TKeyCommand struct{}

func (j TKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "categories" {
		m.SelectedCategory = "tasks"
		m.CurrentView = "details"
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	}

	return nil
}

func (j TKeyCommand) HelpText() string {
	return "TKeyCommand help text"
}

func (j TKeyCommand) AllowedStates() []string {
	return []string{}
}
