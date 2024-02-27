package app

type MKeyCommand struct{}

func (j MKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() {
		m.GoToNextViewWithCategory("meetings")
	}

	return nil
}

func (j MKeyCommand) HelpText() string {
	return "MKeyCommand help text"
}

func (j MKeyCommand) AllowedStates() []string {
	return []string{}
}
