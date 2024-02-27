package app

type LKeyCommand struct{}

func (j LKeyCommand) Execute(m *Model) error {
	if m.IsDetailsView() {
		m.GainDetailsFocus()
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
