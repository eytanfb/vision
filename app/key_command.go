package app

type KeyCommand interface {
	Execute(m *Model) error
	HelpText() string
	AllowedStates() []string
}
