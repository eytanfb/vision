package app

type KeyCommandFactory struct {
	registry *CommandRegistry
}

func NewKeyCommandFactory() *KeyCommandFactory {
	registry := NewRegistry()

	// Navigation commands
	registry.Register("j", JKeyCommand{})
	registry.Register("k", KKeyCommand{})
	registry.Register("h", HKeyCommand{})
	registry.Register("l", LKeyCommand{})
	registry.Register("g", GKeyCommand{})
	registry.Register("tab", TabKeyCommand{})
	registry.Register("shift + tab", ShiftTabKeyCommand{})

	// File operations
	registry.Register("e", EKeyCommand{})
	registry.Register("o", OKeyCommand{})
	registry.Register("n", NKeyCommand{})
	registry.Register("f", FKeyCommand{})

	// Task operations
	registry.Register("d", DKeyCommand{})
	registry.Register("s", SKeyCommand{})
	registry.Register("p", PKeyCommand{})
	registry.Register("D", UppercaseDKeyCommand{})
	registry.Register("S", UppercaseSKeyCommand{})
	registry.Register("a", AKeyCommand{})
	registry.Register("A", UppercaseAKeyCommand{})

	// View control
	registry.Register("c", CKeyCommand{})
	registry.Register("w", WKeyCommand{})
	registry.Register("W", UppercaseWKeyCommand{})
	registry.Register("1", OneKeyCommand{})
	registry.Register("2", TwoKeyCommand{})
	registry.Register("3", ThreeKeyCommand{})
	registry.Register("+", PlusKeyCommand{})
	registry.Register("-", MinusKeyCommand{})
	registry.Register("C", UppercaseCKeyCommand{})
	registry.Register("Q", UppercaseQKeyCommand{})
	registry.Register("L", UppercaseLKeyCommand{})

	// Input handling
	registry.Register("enter", EnterKeyCommand{})
	registry.Register("esc", EscKeyCommand{})
	registry.Register("/", SlashKeyCommand{})
	registry.Register("t", TKeyCommand{})
	registry.Register("m", MKeyCommand{})

	return &KeyCommandFactory{
		registry: registry,
	}
}

func (kcf KeyCommandFactory) CreateKeyCommand(key string) KeyCommand {
	cmd := kcf.registry.Get(key)
	if cmd != nil {
		// Wrap the command to implement the old KeyCommand interface
		return &CommandAdapter{command: cmd}
	}
	return NilKeyCommand{}
}

// CommandAdapter adapts the new Command interface to the old KeyCommand interface
type CommandAdapter struct {
	command Command
}

func (ca *CommandAdapter) Execute(m *Model) error {
	return ca.command.Execute(m)
}

func (ca *CommandAdapter) HelpText() string {
	return ca.command.Description()
}

func (ca *CommandAdapter) AllowedStates() []string {
	return ca.command.Contexts()
}
