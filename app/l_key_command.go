package app

type LKeyCommand struct{}

func (j LKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "details" {
		m.ItemDetailsFocus = true
	} else {
		return EnterKeyCommand{}.Execute(m)
	}

	return nil
}

func (j LKeyCommand) HelpText() string {
	return "LKeyCommand help text"
}

func (j LKeyCommand) AllowedStates() []string {
	return []string{}
}
