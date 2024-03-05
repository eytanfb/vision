package app

type UppercaseLKeyCommand struct{}

func (j UppercaseLKeyCommand) Execute(m *Model) error {
	m.GoToCompany("lifeplus")

	return nil
}

func (j UppercaseLKeyCommand) HelpText() string {
	return "UppercaseLKeyCommand help text"
}

func (j UppercaseLKeyCommand) AllowedStates() []string {
	return []string{}
}
