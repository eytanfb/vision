package app

type EnterKeyCommand struct{}

func (j EnterKeyCommand) Execute(m *Model) error {
	m.Select()

	return nil
}

func (j EnterKeyCommand) HelpText() string {
	return "EnterKeyCommand help text"
}

func (j EnterKeyCommand) AllowedStates() []string {
	return []string{}
}
