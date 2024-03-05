package app

type TwoKeyCommand struct{}

func (j TwoKeyCommand) Execute(m *Model) error {
	return UppercaseQKeyCommand{}.Execute(m)
}

func (j TwoKeyCommand) HelpText() string {
	return "TwoKeyCommand help text"
}

func (j TwoKeyCommand) AllowedStates() []string {
	return []string{}
}
