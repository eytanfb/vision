package app

type MinusKeyCommand struct{}

func (j MinusKeyCommand) Execute(m *Model) error {
	m.TaskManager.ChangeWeeklySummaryToPreviousWeek()

	return nil
}

func (j MinusKeyCommand) HelpText() string {
	return "MinusKeyCommand help text"
}

func (j MinusKeyCommand) AllowedStates() []string {
	return []string{}
}
