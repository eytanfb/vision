package app

type LKeyCommand struct{}

func (j LKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			m.GoToNextKanbanList()
			m.ViewManager.KanbanTaskCursor = 0

			return nil
		}
	}

	return EnterKeyCommand{}.Execute(m)
}

func (j LKeyCommand) HelpText() string {
	return "LKeyCommand help text"
}

func (j LKeyCommand) AllowedStates() []string {
	return []string{}
}
