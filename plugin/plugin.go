package plugin

import (
	"context"
	"fmt"

	generated "github.com/secmc/plugin/proto/generated/go"
	"google.golang.org/grpc"
)

type Plugin struct {
	stream generated.Plugin_EventStreamClient
}

func NewPlugin(name string) (*Plugin, error) {
	conn, err := grpc.NewClient("unix:///tmp/dragonfly_plugin.sock")
	if err != nil {
		return nil, err
	}
	client := generated.NewPluginClient(conn)
	stream, err := client.EventStream(context.Background())
	if err != nil {
		return nil, err
	}
	pl := &Plugin{
		stream: stream,
	}

	pl.stream.Send(&generated.PluginToHost{
		Payload: &generated.PluginToHost_Hello{
			Hello: &generated.PluginHello{
				Name: name,
			},
		},
	})
	go pl.handleMessages()
	return pl, nil
}

func (p *Plugin) handleMessages() {
	for {
		msg, err := p.stream.Recv()
		if err != nil {
			return
		}

		payload := msg.GetPayload()
		if _, ok := payload.(*generated.HostToPlugin_Hello); ok {
			fmt.Println("plugin successfully registered")
		}
	}
}
