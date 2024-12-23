package app

type CKeyCommand struct{}

func (c CKeyCommand) Execute(m *Model) error {
	m.ViewManager.IsCalendarView = !m.ViewManager.IsCalendarView
	return nil
}

func (c CKeyCommand) HelpText() string {
	return "CKeyCommand help text"
}

func (c CKeyCommand) AllowedStates() []string {
	return []string{}
}
