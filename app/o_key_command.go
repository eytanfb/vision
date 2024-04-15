package app

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type OKeyCommand struct{}

func (j OKeyCommand) Execute(m *Model) error {
	if m.IsDetailsView() || m.IsKanbanView() {
		filePath := m.FileManager.SelectedFile.FullPath
		homeDir, _ := os.UserHomeDir()
		notesPath := homeDir + "/Notes"
		//obsidian: //open?vault=Disk-X&file={file$}
		obsidianPath := constructObsidianURL(filePath, notesPath)

		cmd := exec.Command("open", "-a", "Obsidian", obsidianPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return nil
}

func (j OKeyCommand) HelpText() string {
	return "OKeyCommand help text"
}

func (j OKeyCommand) AllowedStates() []string {
	return []string{}
}

func constructObsidianURL(fullPath, notesPath string) string {
	relativePath := strings.Replace(fullPath, notesPath, "", 1)
	urlEncodedPath := url.PathEscape(relativePath) // This handles spaces and other characters.
	obsidianPath := "obsidian://open?vault=Notes&file=" + urlEncodedPath
	return obsidianPath
}
