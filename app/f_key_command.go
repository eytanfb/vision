package app

type FKeyCommand struct{}

func (j FKeyCommand) Execute(m *Model) error {
	m.ViewManager.ToggleHideSidebar()

	return nil
}

func (j FKeyCommand) HelpText() string {
	return "FKeyCommand help text"
}

func (j FKeyCommand) AllowedStates() []string {
	return []string{}
}
