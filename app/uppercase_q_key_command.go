package app

type UppercaseQKeyCommand struct{}

func (j UppercaseQKeyCommand) Execute(m *Model) error {
	m.GoToCompany("qvest.us")

	return nil
}

func (j UppercaseQKeyCommand) HelpText() string {
	return "UppercaseQKeyCommand help text"
}

func (j UppercaseQKeyCommand) AllowedStates() []string {
	return []string{}
}
