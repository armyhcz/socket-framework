package socket_framework

import (
	"context"
	"testing"

	"socket-framework/packet"

	"socket-framework/processor"
)

func TestStart(t *testing.T) {
	p := processor.New()
	p.Add("Ping", Ping)
	p.Add("ping", Ping)
	err := Start("/ws", ":8033", p)
	if err != nil {
		t.Error(err)
	}
}

func Ping(ctx context.Context, packet *packet.Packet) (*packet.Packet, error) {
	return packet.Build("pong", nil), nil
}
