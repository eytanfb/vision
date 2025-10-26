package app

import tea "github.com/charmbracelet/bubbletea"

// Command represents a keyboard command that can be executed
type Command interface {
	// Execute runs the command on the given model and returns a tea.Cmd
	// The command can return nil (no further action) or a tea.Cmd that will
	// generate messages to be processed by Update()
	Execute(model *Model) tea.Cmd

	// Description returns a human-readable description of what this command does
	Description() string

	// Contexts returns which views this command is available in
	// Empty slice means available in all contexts
	Contexts() []string
}
