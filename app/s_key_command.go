package app

type SKeyCommand struct{}

func (j SKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() {
		m.GoToNextViewWithCategory("standups")
	}

	return nil
}

func (j SKeyCommand) HelpText() string {
	return "SKeyCommand help text"
}

func (j SKeyCommand) AllowedStates() []string {
	return []string{}
}
