package app

type HKeyCommand struct{}

func (j HKeyCommand) Execute(m *Model) error {
	return EscKeyCommand{}.Execute(m)
}

func (j HKeyCommand) HelpText() string {
	return "HKeyCommand help text"
}

func (j HKeyCommand) AllowedStates() []string {
	return []string{}
}
