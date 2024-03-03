package app

type JKeyCommand struct{}

func (j JKeyCommand) Execute(m *Model) error {
	moveDown(m)

	return nil
}

func (j JKeyCommand) HelpText() string {
	return "JKeyCommand help text"
}

func (j JKeyCommand) AllowedStates() []string {
	return []string{}
}

func moveDown(m *Model) {
	if m.IsCompanyView() {
		m.GoToNextCompany()
	} else if m.IsCategoryView() {
		m.GoToNextCategory()
	} else if m.IsDetailsView() {
		if m.IsItemDetailsFocus() {
			if m.IsTaskDetailsFocus() {
				m.GoToNextTask()
			} else {
				m.Viewport.LineDown(10)
			}
		} else {
			m.GoToNextFile()
			m.Viewport.GotoTop()
		}
	}
}
