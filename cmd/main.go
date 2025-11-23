package main

import (
	"log"

	"github.com/secmc/plugin-go/plugin"
)

func main() {
	commands := []*plugin.Command{
		plugin.NewCommand("plugin", "", nil, testCommands{}),
	}

	_, err := plugin.NewPlugin("test", plugin.WithCommandsOpt(commands...))
	if err != nil {
		log.Fatalln(err)
	}
	select {}
}

type testCommands struct {
}

func (testCommands) Run() {}
