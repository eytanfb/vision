package app

import tea "github.com/charmbracelet/bubbletea"

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
	registry.Register("shift+tab", ShiftTabKeyCommand{})

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

func (kcf KeyCommandFactory) CreateKeyCommand(key string) Command {
	cmd := kcf.registry.Get(key)
	if cmd != nil {
		return &ContextValidatingCommand{command: cmd}
	}
	return NilCommand{}
}

// ContextValidatingCommand wraps a Command and validates context before execution
type ContextValidatingCommand struct {
	command Command
}

func (cvc *ContextValidatingCommand) Execute(m *Model) tea.Cmd {
	contexts := cvc.command.Contexts()
	if len(contexts) > 0 && !m.MatchesContext(contexts) {
		return nil // Command not valid in current context
	}
	return cvc.command.Execute(m)
}

func (cvc *ContextValidatingCommand) Description() string {
	return cvc.command.Description()
}

func (cvc *ContextValidatingCommand) Contexts() []string {
	return cvc.command.Contexts()
}
