package app

import (
	"fmt"
	"os"
	"os/exec"
)

type EKeyCommand struct{}

func (j EKeyCommand) Execute(m *Model) error {
	if m.CurrentView != "details" && !m.ItemDetailsFocus {
		return nil
	}

	filePath := getCurrentFilePath(m)
	editor := os.Getenv("EDITOR")
	fmt.Println(editor)
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

func getCurrentFilePath(m *Model) string {
	filePath := "/Users/eytananjel/Notes/" + m.SelectedCompany.DisplayName + "/" + m.SelectedCategory + "/" + m.Files[m.FilesCursor].Name
	return filePath
}
