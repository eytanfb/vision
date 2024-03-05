package app

type UppercaseCKeyCommand struct{}

func (j UppercaseCKeyCommand) Execute(m *Model) error {
	m.GoToCompany("clerky")

	return nil
}

func (j UppercaseCKeyCommand) HelpText() string {
	return "UppercaseCKeyCommand help text"
}

func (j UppercaseCKeyCommand) AllowedStates() []string {
	return []string{}
}
