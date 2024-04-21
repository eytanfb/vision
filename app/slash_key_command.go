package app

type SlashKeyCommand struct{}

func (j SlashKeyCommand) Execute(m *Model) error {
	m.ViewManager.IsFilterView = true

	return nil
}

func (j SlashKeyCommand) HelpText() string {
	return "SlashKeyCommand help text"
}

func (j SlashKeyCommand) AllowedStates() []string {
	return []string{}
}
