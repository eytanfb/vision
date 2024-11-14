package app

type TabKeyCommand struct{}

func (j TabKeyCommand) Execute(m *Model) error {
	if m.ViewManager.SuggestionsListsCursor == -1 {
		m.ViewManager.SuggestionsListsCursor = 0
	}

	m.ViewManager.NextSuggestion(&m.FileManager)

	return nil
}

func (j TabKeyCommand) HelpText() string {
	return "TabKeyCommand help text"
}

func (j TabKeyCommand) AllowedStates() []string {
	return []string{}
}
