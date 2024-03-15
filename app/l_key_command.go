package app

type LKeyCommand struct{}

func (j LKeyCommand) Execute(m *Model) error {
	return EnterKeyCommand{}.Execute(m)
}

func (j LKeyCommand) HelpText() string {
	return "LKeyCommand help text"
}

func (j LKeyCommand) AllowedStates() []string {
	return []string{}
}
