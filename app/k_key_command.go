package app

type KKeyCommand struct{}

func (j KKeyCommand) Execute(m *Model) error {
	m.MoveUp()

	return nil
}

func (j KKeyCommand) HelpText() string {
	return "KKeyCommand help text"
}

func (j KKeyCommand) AllowedStates() []string {
	return []string{}
}
