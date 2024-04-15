package app

type HKeyCommand struct{}

func (j HKeyCommand) Execute(m *Model) error {
	if m.IsCategoryView() {
		if m.ViewManager.HideSidebar {
			m.GoToPreviousKanbanList()
			m.ViewManager.KanbanTaskCursor = 0

			return nil
		}
	}

	return EscKeyCommand{}.Execute(m)
}

func (j HKeyCommand) HelpText() string {
	return "HKeyCommand help text"
}

func (j HKeyCommand) AllowedStates() []string {
	return []string{}
}
