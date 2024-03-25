package app

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/log"
)

type UppercaseDKeyCommand struct{}

func (j UppercaseDKeyCommand) Execute(m *Model) error {
	if !m.IsCategoryView() || m.ViewManager.IsWeeklyView {
		return nil
	}

	slackMessage := m.TaskManager.SummaryForSlack(m.DirectoryManager.CurrentCompanyName())
	err := clipboard.WriteAll(slackMessage)
	if err != nil {
		log.Error("Failed to copy to clipboard", err)
	}

	return nil
}

func (j UppercaseDKeyCommand) HelpText() string {
	return "UppercaseDKeyCommand help text"
}

func (j UppercaseDKeyCommand) AllowedStates() []string {
	return []string{}
}
