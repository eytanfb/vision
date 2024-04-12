package app

type TKeyCommand struct{}

func (j TKeyCommand) Execute(m *Model) error {
	if m.IsItemDetailsFocus() {
		m.TaskManager.ChangeDailySummaryDateToToday()
		m.ShowTasks()
	}

	return nil
}

func (j TKeyCommand) HelpText() string {
	return "TKeyCommand help text"
}

func (j TKeyCommand) AllowedStates() []string {
	return []string{}
}
