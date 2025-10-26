package app

// Command represents a keyboard command that can be executed
type Command interface {
	// Execute runs the command on the given model
	Execute(model *Model) error

	// Description returns a human-readable description of what this command does
	Description() string

	// Contexts returns which views this command is available in
	// Empty slice means available in all contexts
	Contexts() []string
}
