package app

import "github.com/charmbracelet/log"

type ShiftTabKeyCommand struct{}

func (j ShiftTabKeyCommand) Execute(m *Model) error {
	log.Info("ShiftTabKeyCommand")
	m.ViewManager.PreviousSuggestion(&m.FileManager)

	return nil
}

func (j ShiftTabKeyCommand) HelpText() string {
	return "ShiftTabKeyCommand help text"
}

func (j ShiftTabKeyCommand) AllowedStates() []string {
	return []string{}
}
