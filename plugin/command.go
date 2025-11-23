package plugin

type Command struct {
	runnables []Runnable

	name        string
	aliases     []string
	description string
}

func NewCommand(name, description string, aliases []string, runnables ...Runnable) *Command {
	return &Command{
		runnables:   runnables,
		name:        name,
		description: description,
		aliases:     aliases,
	}
}

type Runnable interface {
	Run()
}
