package app

// CommandRegistry maps keyboard keys to their corresponding commands
type CommandRegistry struct {
	commands map[string]Command
}

// NewRegistry creates a new command registry
func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

// Register adds a command for a specific key
func (r *CommandRegistry) Register(key string, cmd Command) {
	r.commands[key] = cmd
}

// Get retrieves a command for a specific key
// Returns the command if found, nil otherwise
func (r *CommandRegistry) Get(key string) Command {
	cmd, ok := r.commands[key]
	if !ok {
		return nil
	}
	return cmd
}

// AllKeys returns all registered keys
func (r *CommandRegistry) AllKeys() []string {
	keys := make([]string, 0, len(r.commands))
	for k := range r.commands {
		keys = append(keys, k)
	}
	return keys
}

// AllCommands returns all registered commands with their keys
func (r *CommandRegistry) AllCommands() map[string]Command {
	result := make(map[string]Command)
	for k, v := range r.commands {
		result[k] = v
	}
	return result
}
