package app

type KKeyCommand struct{}

func (j KKeyCommand) Execute(m *Model) error {
	moveUp(m)

	return nil
}

func (j KKeyCommand) HelpText() string {
	return "KKeyCommand help text"
}

func (j KKeyCommand) AllowedStates() []string {
	return []string{}
}

func moveUp(m *Model) {
	if m.IsDetailsView() {
		if m.IsItemDetailsFocus() {
			if m.IsTaskDetailsFocus() {
				m.GoToPreviousTask()
			} else {
				m.Viewport.LineUp(10)
			}
		} else {
			m.GoToPreviousFile()
			m.Viewport.GotoTop()
		}
	} else if m.IsCategoryView() {
		if !m.ViewManager.HideSidebar {
			m.GoToPreviousCategory()
		} else {
			m.GoToPreviousKanbanTask()
		}
	} else if m.IsCompanyView() {
		m.GoToPreviousCompany()
	}
}
