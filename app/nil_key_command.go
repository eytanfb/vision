package app

type NilKeyCommand struct{}

func (j NilKeyCommand) Execute(m *Model) error {
	return nil
}

func (j NilKeyCommand) HelpText() string {
	return ""
}

func (j NilKeyCommand) AllowedStates() []string {
	return []string{}
}
