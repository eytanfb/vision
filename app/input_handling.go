package app

import (
	"strings"
	"vision/utils"
)

// InputHandling handles input-related commands
type InputHandling struct{}

// HandleKey routes input handling keys to their appropriate handlers
func (ih InputHandling) HandleKey(key string, m *Model) error {
	switch key {
	case "enter":
		return ih.HandleEnter(m)
	case "esc":
		return ih.HandleEscape(m)
	case "/":
		return ih.StartFilter(m)
	case "t":
		return ih.ShowTasks(m)
	case "m":
		return ih.GoToMeetings(m)
	}
	return nil
}

// HandleEnter processes the enter key
func (ih InputHandling) HandleEnter(m *Model) error {
	if m.IsSuggestionsActive() {
		acceptedSuggestion := m.FileManager.GetActiveSuggestion(m.ViewManager.SuggestionsListsCursor, m.ViewManager.SuggestionCursor)

		if acceptedSuggestion != "" {
			currentValue := m.NewTaskInput.Value()
			filterValue := m.FileManager.SuggestionsFilterValue

			if filterValue == "" {
				filterValue = "[["
			}

			newValue := strings.Replace(currentValue, filterValue, acceptedSuggestion, 1)

			if filterValue == "[[" {
				newValue = "[[" + newValue
			}

			m.NewTaskInput.SetValue(newValue + "]]")
			m.NewTaskInput.SetCursor(len(m.NewTaskInput.Value()))
		}

		return ih.HandleEscape(m)
	}

	if m.IsAddTaskView() {
		company := m.GetCurrentCompanyName()
		input := m.NewTaskInput.Value()

		if err := m.FileManager.CreateTask(company, input); err != nil {
			m.Errors = append(m.Errors, err.Error())
		}
		return ih.HandleEscape(m)
	} else if m.IsAddSubTaskView() {
		company := m.GetCurrentCompanyName()
		input := m.NewTaskInput.Value()
		selectedFile := m.FileManager.SelectedFile

		// Parse hashtags to Obsidian date format before creating subtask
		processedInput := utils.ParseHashtagsToObsidianDates(input)

		if err := m.FileManager.CreateSubTask(company, selectedFile, processedInput); err != nil {
			m.Errors = append(m.Errors, err.Error())
		}
		return ih.HandleEscape(m)
	} else if m.IsFilterView() {
		m.ViewManager.IsFilterView = false
		m.TaskManager.TaskCollection.FilterValue = m.FilterInput.Value()
		m.TaskManager.TasksCursor = 0
		m.FilterInput.Blur()

		return nil
	}

	m.Select()

	return nil
}

// HandleEscape processes the escape key
func (ih InputHandling) HandleEscape(m *Model) error {
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
		// Toggle sidebar when in task details
		m.ViewManager.ToggleHideSidebar()
		goToPreviousView = false
	}

	m.FileManager.FetchFiles(&m.DirectoryManager, &m.TaskManager)

	if goToPreviousView {
		m.GoToPreviousView()
	}

	return nil
}

// StartFilter activates filter mode
func (ih InputHandling) StartFilter(m *Model) error {
	m.ViewManager.IsFilterView = true
	return nil
}

// ShowTasks shows the tasks view
func (ih InputHandling) ShowTasks(m *Model) error {
	if m.IsItemDetailsFocus() {
		m.TaskManager.TasksCursor = -1
		m.TaskManager.ChangeDailySummaryDateToToday()
		m.ShowTasks()
	}
	return nil
}

// GoToMeetings navigates to the meetings category
func (ih InputHandling) GoToMeetings(m *Model) error {
	if m.IsCategoryView() {
		m.GoToNextViewWithCategory("meetings")
	}
	return nil
}

// Command implementations for registry

type EnterKeyCommand struct{}

func (cmd EnterKeyCommand) Execute(m *Model) error {
	return InputHandling{}.HandleEnter(m)
}

func (cmd EnterKeyCommand) Description() string {
	return "Confirm selection or input"
}

func (cmd EnterKeyCommand) Contexts() []string {
	return []string{}
}

type EscKeyCommand struct{}

func (cmd EscKeyCommand) Execute(m *Model) error {
	return InputHandling{}.HandleEscape(m)
}

func (cmd EscKeyCommand) Description() string {
	return "Go back or cancel"
}

func (cmd EscKeyCommand) Contexts() []string {
	return []string{}
}

type SlashKeyCommand struct{}

func (cmd SlashKeyCommand) Execute(m *Model) error {
	return InputHandling{}.StartFilter(m)
}

func (cmd SlashKeyCommand) Description() string {
	return "Start filtering"
}

func (cmd SlashKeyCommand) Contexts() []string {
	return []string{}
}

type TKeyCommand struct{}

func (cmd TKeyCommand) Execute(m *Model) error {
	return InputHandling{}.ShowTasks(m)
}

func (cmd TKeyCommand) Description() string {
	return "Show tasks for today"
}

func (cmd TKeyCommand) Contexts() []string {
	return []string{"item_details"}
}

type MKeyCommand struct{}

func (cmd MKeyCommand) Execute(m *Model) error {
	return InputHandling{}.GoToMeetings(m)
}

func (cmd MKeyCommand) Description() string {
	return "Go to meetings"
}

func (cmd MKeyCommand) Contexts() []string {
	return []string{"category_view"}
}
