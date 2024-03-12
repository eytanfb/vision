package app

type WKeyCommand struct{}

func (j WKeyCommand) Execute(m *Model) error {
	m.ViewManager.ToggleWeeklyView()

	return nil
}

func (j WKeyCommand) HelpText() string {
	return "WKeyCommand help text"
}

func (j WKeyCommand) AllowedStates() []string {
	return []string{}
}
