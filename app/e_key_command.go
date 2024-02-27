package app

import (
	"os"
	"os/exec"
)

type EKeyCommand struct{}

func (j EKeyCommand) Execute(m *Model) error {
	if !m.IsItemDetailsFocus() {
		return nil
	}

	filePath := m.GetCurrentFilePath()
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vim" // Default to vim if $EDITOR is not set
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func (j EKeyCommand) HelpText() string {
	return "EKeyCommand help text"
}

func (j EKeyCommand) AllowedStates() []string {
	return []string{}
}
