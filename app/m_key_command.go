package app

import "fmt"

type MKeyCommand struct{}

func (j MKeyCommand) Execute(m *Model) error {
	if m.CurrentView == "categories" {
		fmt.Println("MKeyCommand")
		m.SelectedCategory = "meetings"
		m.CurrentView = "details"
		m.FilesCursor = 0
		m.Files = m.FetchFiles()
	}

	return nil
}

func (j MKeyCommand) HelpText() string {
	return "MKeyCommand help text"
}

func (j MKeyCommand) AllowedStates() []string {
	return []string{}
}
