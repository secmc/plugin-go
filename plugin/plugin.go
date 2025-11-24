package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/secmc/plugin-go/plugin/player"
	generated "github.com/secmc/plugin/proto/generated/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Plugin struct {
	id     string
	name   string
	stream generated.Plugin_EventStreamClient

	commands []*Command

	playerMu sync.Mutex
	players  map[string]*player.Player
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
	case *generated.EventEnvelope_PlayerJoin:
		pl := player.New(event.PlayerJoin.Name, event.PlayerJoin.PlayerUuid)
		p.addPlayer(pl)
	case *generated.EventEnvelope_PlayerQuit:
		pl, ok := p.playerFromUUID(event.PlayerQuit.PlayerUuid)
		if !ok {
			return
		}
		p.deletePlayer(pl)
	}
}

func (p *Plugin) addPlayer(pl *player.Player) {
	p.playerMu.Lock()
	p.players[p.id] = pl
	p.playerMu.Unlock()
}

func (p *Plugin) deletePlayer(pl *player.Player) {
	p.playerMu.Lock()
	delete(p.players, pl.UUID())
	p.playerMu.Unlock()
}

func (p *Plugin) playerFromUUID(uuid string) (*player.Player, bool) {
	p.playerMu.Lock()
	pl, ok := p.players[uuid]
	p.playerMu.Unlock()
	return pl, ok
}

func (p *Plugin) executeCommand(name string) {
	for _, cmd := range p.commands {
		if strings.EqualFold(cmd.name, name) {
			//
		}
	}
}
