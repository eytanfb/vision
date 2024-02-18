package command

type Command interface {
	Execute() (*main.Model, error)
}

type CommandFactory struct{}

func (cf *CommandFactory) CreateCommand() Command {
	return Command{}
}
