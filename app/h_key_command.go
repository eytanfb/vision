package app

type HKeyCommand struct{}

func (j HKeyCommand) Execute(m *Model) error {
	if m.IsItemDetailsFocus() {
		m.LoseDetailsFocus()
	} else {
		return EscKeyCommand{}.Execute(m)
	}

	return nil
}

func (j HKeyCommand) HelpText() string {
	return "HKeyCommand help text"
}

func (j HKeyCommand) AllowedStates() []string {
	return []string{}
}
