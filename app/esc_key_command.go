package app

type EscKeyCommand struct{}

func (j EscKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "categories" {
		m.CurrentView = "companies"
	} else if m.CurrentView == "details" {
		m.CurrentView = "categories"
	}

	return nil
}

func (j EscKeyCommand) HelpText() string {
	return "EscKeyCommand help text"
}

func (j EscKeyCommand) AllowedStates() []string {
	return []string{}
}
