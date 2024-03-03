package app

type NKeyCommand struct{}

func (j NKeyCommand) Execute(m *Model) error {
	m.GoToNextCompany()
	m.GotoCategoryView()

	return nil
}

func (j NKeyCommand) HelpText() string {
	return "NKeyCommand help text"
}

func (j NKeyCommand) AllowedStates() []string {
	return []string{}
}
