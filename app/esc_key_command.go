package app

type EscKeyCommand struct{}

func (j EscKeyCommand) Execute(m *Model) error {
	m.GoToPreviousView()

	return nil
}

func (j EscKeyCommand) HelpText() string {
	return "EscKeyCommand help text"
}

func (j EscKeyCommand) AllowedStates() []string {
	return []string{}
}
