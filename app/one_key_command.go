package app

type OneKeyCommand struct{}

func (j OneKeyCommand) Execute(m *Model) error {
	return UppercaseCKeyCommand{}.Execute(m)
}

func (j OneKeyCommand) HelpText() string {
	return "OneKeyCommand help text"
}

func (j OneKeyCommand) AllowedStates() []string {
	return []string{}
}
