package app

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

// ViewControl handles view-related commands
type ViewControl struct{}

// HandleKey routes view control keys to their appropriate handlers
func (vc ViewControl) HandleKey(key string, m *Model) tea.Cmd {
	switch key {
	case "c":
		return vc.ToggleCalendarView(m)
	case "w":
		return vc.ToggleWeeklyView(m)
	case "W":
		return vc.CopyWeeklySummary(m)
	case "1":
		return vc.GoToClerky(m)
	case "2":
		return vc.GoToQvest(m)
	case "3":
		return vc.GoToLifeplus(m)
	case "+":
		return vc.NextDay(m)
	case "-":
		return vc.PreviousDay(m)
	case "C":
		return vc.GoToClerky(m)
	case "Q":
		return vc.GoToQvest(m)
	case "L":
		return vc.GoToLifeplus(m)
	}
	return nil
}

// ToggleCalendarView toggles the calendar view on/off
func (vc ViewControl) ToggleCalendarView(m *Model) tea.Cmd {
	m.ViewManager.IsCalendarView = !m.ViewManager.IsCalendarView
	return nil
}

// ToggleWeeklyView toggles the weekly view on/off
func (vc ViewControl) ToggleWeeklyView(m *Model) tea.Cmd {
	m.ViewManager.ToggleWeeklyView()
	return nil
}

// CopyWeeklySummary copies the weekly summary to clipboard
func (vc ViewControl) CopyWeeklySummary(m *Model) tea.Cmd {
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

// GoToClerky switches to the Clerky company
func (vc ViewControl) GoToClerky(m *Model) tea.Cmd {
	m.GoToCompany("clerky")
	return nil
}

// GoToQvest switches to the Qvest company
func (vc ViewControl) GoToQvest(m *Model) tea.Cmd {
	m.GoToCompany("qvest.us")
	return nil
}

// GoToLifeplus switches to the Lifeplus company
func (vc ViewControl) GoToLifeplus(m *Model) tea.Cmd {
	m.GoToCompany("lifeplus")
	return nil
}

// NextDay moves to the next day or week
func (vc ViewControl) NextDay(m *Model) tea.Cmd {
	if !m.ViewManager.IsTaskDetailsFocus() {
		if !m.ViewManager.IsWeeklyView {
			m.TaskManager.ChangeDailySummaryDateToNextDay()
		} else {
			m.TaskManager.ChangeWeeklySummaryToNextWeek()
		}
	}
	return nil
}

// PreviousDay moves to the previous day or week
func (vc ViewControl) PreviousDay(m *Model) tea.Cmd {
	if !m.ViewManager.IsTaskDetailsFocus() {
		if !m.ViewManager.IsWeeklyView {
			m.TaskManager.ChangeDailySummaryDateToPreviousDay()
		} else {
			m.TaskManager.ChangeWeeklySummaryToPreviousWeek()
		}
	}
	return nil
}

// Command implementations for registry

type CKeyCommand struct{}

func (cmd CKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.ToggleCalendarView(m)
}

func (cmd CKeyCommand) Description() string {
	return "Toggle calendar view"
}

func (cmd CKeyCommand) Contexts() []string {
	return []string{}
}

type WKeyCommand struct{}

func (cmd WKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.ToggleWeeklyView(m)
}

func (cmd WKeyCommand) Description() string {
	return "Toggle weekly view"
}

func (cmd WKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseWKeyCommand struct{}

func (cmd UppercaseWKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.CopyWeeklySummary(m)
}

func (cmd UppercaseWKeyCommand) Description() string {
	return "Copy weekly summary to clipboard"
}

func (cmd UppercaseWKeyCommand) Contexts() []string {
	return []string{"weekly_view"}
}

type OneKeyCommand struct{}

func (cmd OneKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToClerky(m)
}

func (cmd OneKeyCommand) Description() string {
	return "Switch to Clerky"
}

func (cmd OneKeyCommand) Contexts() []string {
	return []string{}
}

type TwoKeyCommand struct{}

func (cmd TwoKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToQvest(m)
}

func (cmd TwoKeyCommand) Description() string {
	return "Switch to Qvest"
}

func (cmd TwoKeyCommand) Contexts() []string {
	return []string{}
}

type ThreeKeyCommand struct{}

func (cmd ThreeKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToLifeplus(m)
}

func (cmd ThreeKeyCommand) Description() string {
	return "Switch to Lifeplus"
}

func (cmd ThreeKeyCommand) Contexts() []string {
	return []string{}
}

type PlusKeyCommand struct{}

func (cmd PlusKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.NextDay(m)
}

func (cmd PlusKeyCommand) Description() string {
	return "Next day/week"
}

func (cmd PlusKeyCommand) Contexts() []string {
	return []string{}
}

type MinusKeyCommand struct{}

func (cmd MinusKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.PreviousDay(m)
}

func (cmd MinusKeyCommand) Description() string {
	return "Previous day/week"
}

func (cmd MinusKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseCKeyCommand struct{}

func (cmd UppercaseCKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToClerky(m)
}

func (cmd UppercaseCKeyCommand) Description() string {
	return "Switch to Clerky"
}

func (cmd UppercaseCKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseQKeyCommand struct{}

func (cmd UppercaseQKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToQvest(m)
}

func (cmd UppercaseQKeyCommand) Description() string {
	return "Switch to Qvest"
}

func (cmd UppercaseQKeyCommand) Contexts() []string {
	return []string{}
}

type UppercaseLKeyCommand struct{}

func (cmd UppercaseLKeyCommand) Execute(m *Model) tea.Cmd {
	return ViewControl{}.GoToLifeplus(m)
}

func (cmd UppercaseLKeyCommand) Description() string {
	return "Switch to Lifeplus"
}

func (cmd UppercaseLKeyCommand) Contexts() []string {
	return []string{}
}
