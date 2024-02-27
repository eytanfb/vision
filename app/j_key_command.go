package app

type JKeyCommand struct{}

func (j JKeyCommand) Execute(m *Model) error {
	m.MoveDown()

	return nil
}

func (j JKeyCommand) HelpText() string {
	return "JKeyCommand help text"
}

func (j JKeyCommand) AllowedStates() []string {
	return []string{}
}
