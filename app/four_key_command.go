package app

type FourKeyCommand struct{}

func (j FourKeyCommand) Execute(m *Model) error {
	m.GoToCompany("personal")

	return nil
}

func (j FourKeyCommand) HelpText() string {
	return "FourKeyCommand help text"
}

func (j FourKeyCommand) AllowedStates() []string {
	return []string{}
}
