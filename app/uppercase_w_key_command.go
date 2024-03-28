package app

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/log"
)

type UppercaseWKeyCommand struct{}

func (j UppercaseWKeyCommand) Execute(m *Model) error {
	if !m.IsCategoryView() || !m.ViewManager.IsWeeklyView {
		return nil
	}

	slackMessage := m.TaskManager.WeeklySummaryForSlack(m.DirectoryManager.CurrentCompanyName())
	err := clipboard.WriteAll(slackMessage)
	if err != nil {
		log.Error("Failed to copy to clipboard", err)
	}

	return nil
}

func (j UppercaseWKeyCommand) HelpText() string {
	return "UppercaseWKeyCommand help text"
}

func (j UppercaseWKeyCommand) AllowedStates() []string {
	return []string{}
}
