package app

import (
	"os"
	"os/exec"
)

type EKeyCommand struct{}

func (j EKeyCommand) Execute(m *Model) error {
	if !m.IsDetailsView() {
		return nil
	}

	//editor := os.Getenv("EDITOR")
	//if editor != "" {
	//splitEditorArgs := strings.Split(editor, " ")
	//if len(splitEditorArgs) > 1 {
	//for _, arg := range splitEditorArgs[1:] {
	//args = append(args, arg)
	//}
	//}
	//}

	//if editor == "" {
	//editor = "vim" // Default to vim if $EDITOR is not set
	//}
	filePath := m.FileManager.SelectedFile.FullPath
	cmd := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)
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
