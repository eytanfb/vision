package app

type DKeyCommand struct{}

func (j DKeyCommand) Execute(m *Model) error {


	return nil
}

func (j DKeyCommand) HelpText() string {
	return "DKeyCommand help text"
}

func (j DKeyCommand) AllowedStates() []string {
	return []string{}
}
