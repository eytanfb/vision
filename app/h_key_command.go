package app

type HKeyCommand struct{}

func (j HKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "details" {
		m.ItemDetailsFocus = false
	}

	return nil
}

func (j HKeyCommand) HelpText() string {
	return "HKeyCommand help text"
}

func (j HKeyCommand) AllowedStates() []string {
	return []string{}
}
