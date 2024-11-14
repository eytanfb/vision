package app

type EscKeyCommand struct{}

func (j EscKeyCommand) Execute(m *Model) error {
	goToPreviousView := true

	if m.IsSuggestionsActive() {
		m.ViewManager.IsSuggestionsActive = false
		m.ViewManager.SuggestionsListsCursor = -1
		m.ViewManager.SuggestionCursor = -1
		return nil
	}

	if m.IsAddTaskView() {
		m.ViewManager.IsAddTaskView = false
		m.NewTaskInput.Blur()
		goToPreviousView = false
	} else if m.IsAddSubTaskView() {
		m.ViewManager.IsAddSubTaskView = false
		m.NewTaskInput.Blur()
		goToPreviousView = false
	} else if m.IsFilterView() {
		m.ViewManager.IsFilterView = false
		m.FilterInput.Blur()
		goToPreviousView = false
	} else if m.ViewManager.IsTaskDetailsFocus() {
		FKeyCommand{}.Execute(m)
		goToPreviousView = false
	}

	m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)

	if goToPreviousView {
		m.GoToPreviousView()
	}

	return nil
}

func (j EscKeyCommand) HelpText() string {
	return "EscKeyCommand help text"
}

func (j EscKeyCommand) AllowedStates() []string {
	return []string{}
}
