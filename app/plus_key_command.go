package app

type PlusKeyCommand struct{}

func (j PlusKeyCommand) Execute(m *Model) error {
	if !m.ViewManager.IsWeeklyView {
		m.TaskManager.ChangeDailySummaryDateToNextDay()
	} else {
		m.TaskManager.ChangeWeeklySummaryToNextWeek()
	}

	return nil
}

func (j PlusKeyCommand) HelpText() string {
	return "PlusKeyCommand help text"
}

func (j PlusKeyCommand) AllowedStates() []string {
	return []string{}
}
