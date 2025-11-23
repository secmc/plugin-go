package main

import (
	"log"

	"github.com/secmc/plugin-go/plugin"
)

func main() {
	commands := []*plugin.Command{
		plugin.NewCommand("plugin", "", []string{}, testCommands{}),
	}

	p, err := plugin.NewPlugin("test", plugin.WithCommandsOpt(commands...))
	if err != nil {
		log.Fatalln(err)
	}

	p.Start()
}

type testCommands struct {
}

func (testCommands) Run() {}
