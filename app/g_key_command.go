package app

import (
	"os"
	"os/exec"
)

type GKeyCommand struct{}

func (j GKeyCommand) Execute(m *Model) error {
	cmd := exec.Command("gh", "dash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (j GKeyCommand) HelpText() string {
	return "GKeyCommand help text"
}

func (j GKeyCommand) AllowedStates() []string {
	return []string{}
}
