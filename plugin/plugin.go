package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	generated "github.com/secmc/plugin/proto/generated/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Plugin struct {
	id     string
	name   string
	stream generated.Plugin_EventStreamClient

	commands []*Command
}

type Opt func(p *Plugin)

func WithCommandsOpt(commands ...*Command) func(p *Plugin) {
	return func(p *Plugin) {
		p.commands = commands
	}
}

func NewPlugin(name string, opts ...Opt) (*Plugin, error) {
	serverHost := os.Getenv("DF_PLUGIN_SERVER_ADDRESS")
	conn, err := grpc.NewClient(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := generated.NewPluginClient(conn)
	stream, err := client.EventStream(context.Background())
	if err != nil {
		return nil, err
	}
	pl := &Plugin{
		id:     os.Getenv("DF_PLUGIN_ID"),
		name:   name,
		stream: stream,
	}

	return pl, nil
}

func (p *Plugin) Start() error {
	var commandSpecs []*generated.CommandSpec
	for _, cmd := range p.commands {
		commandSpecs = append(commandSpecs, &generated.CommandSpec{
			Name:        cmd.name,
			Aliases:     cmd.aliases,
			Description: cmd.description,
		})
	}

	err := p.stream.Send(&generated.PluginToHost{
		PluginId: p.id,
		Payload: &generated.PluginToHost_Hello{
			Hello: &generated.PluginHello{
				Name:     p.name,
				Commands: commandSpecs,
			},
		},
	})
	if err != nil {
		return err
	}

	return p.handleMessages()
}

func (p *Plugin) handleMessages() error {
	for {
		msg, err := p.stream.Recv()
		if err != nil {
			return err
		}

		if event := msg.GetEvent(); event != nil {
			p.handleEvent(event)
			continue
		}

		switch payload := msg.GetPayload().(type) {
		case *generated.HostToPlugin_Hello:
			fmt.Println("Golang plugin successfully registered")
		case *generated.HostToPlugin_ServerInfo:
			fmt.Println(payload.ServerInfo.Plugins)
		}
	}
}

func (p *Plugin) handleEvent(e *generated.EventEnvelope) {
	switch event := e.Payload.(type) {
	case *generated.EventEnvelope_Command:
		p.executeCommand(event.Command.Name)
	}
}

func (p *Plugin) executeCommand(name string) {
	for _, cmd := range p.commands {
		if strings.EqualFold(cmd.name, name) {
		}
	}
}
