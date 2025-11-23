package plugin

import (
	"context"
	"fmt"
	"log"
	"os"

	generated "github.com/secmc/plugin/proto/generated/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Plugin struct {
	stream generated.Plugin_EventStreamClient
}

func NewPlugin(name string) (*Plugin, error) {
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
		stream: stream,
	}

	pluginID := os.Getenv("DF_PLUGIN_ID")
	err = pl.stream.Send(&generated.PluginToHost{
		PluginId: pluginID,
		Payload: &generated.PluginToHost_Hello{
			Hello: &generated.PluginHello{
				Name: name,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	go pl.handleMessages()
	return pl, nil
}

func (p *Plugin) handleMessages() {
	for {
		msg, err := p.stream.Recv()
		if err != nil {
			log.Fatalln(err)
			return
		}

		payload := msg.GetPayload()
		if _, ok := payload.(*generated.HostToPlugin_Hello); ok {
			fmt.Println("Golang plugin successfully registered")
		}
	}
}
