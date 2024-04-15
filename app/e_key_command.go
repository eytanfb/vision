package app

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

type EKeyCommand struct{}

func (j EKeyCommand) Execute(m *Model) error {
	if m.IsDetailsView() || m.IsKanbanView() {
		log.Info("Opening file in vim", m.FileManager.SelectedFile.Name)
		filePath := m.FileManager.SelectedFile.FullPath
		cmd := exec.Command("vim", "-u", "~/.dotfiles/.vimrc", filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return nil
}

func (j EKeyCommand) HelpText() string {
	return "EKeyCommand help text"
}

func (j EKeyCommand) AllowedStates() []string {
	return []string{}
}
