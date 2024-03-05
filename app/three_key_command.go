package app

type ThreeKeyCommand struct{}

func (j ThreeKeyCommand) Execute(m *Model) error {
	return UppercaseLKeyCommand{}.Execute(m)
}

func (j ThreeKeyCommand) HelpText() string {
	return "ThreeKeyCommand help text"
}

func (j ThreeKeyCommand) AllowedStates() []string {
	return []string{}
}
