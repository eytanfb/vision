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
	if m.IsSuggestionsActive() {
		m.ViewManager.NextSuggestion(&m.FileManager)
	} else if m.IsCategoryView() {
		categoryViewBehavior(m)
	} else if m.IsDetailsView() {
		detailsViewBehavior(m)
	}
}

func categoryViewBehavior(m *Model) {
	if !m.ViewManager.HideSidebar {
		m.GoToNextCategory()
	} else {
		m.GoToNextKanbanTask()
	}
}

func detailsViewBehavior(m *Model) {
	if m.IsItemDetailsFocus() {
		itemDetailsViewBehavior(m)
	} else {
		m.GoToNextFile()
		m.Viewport.GotoTop()
	}
}

func itemDetailsViewBehavior(m *Model) {
	if m.IsTaskDetailsFocus() {
		m.GoToNextTask()
	} else {
		m.Viewport.LineDown(10)
	}
}
