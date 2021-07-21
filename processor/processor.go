package processor

import (
	"context"
	"errors"

	"socket-framework/packet"
)

const (
	TypeCloseEvent = "closed"
)

type PacketHandler func(ctx context.Context, packet *packet.Packet) (*packet.Packet, error)

type Consumers struct {
	Handlers map[string]PacketHandler
}

func (c *Consumers) Add(t string, handler PacketHandler) {
	if c.Handlers == nil {
		h := make(map[string]PacketHandler)
		h[t] = handler
		c.Handlers = h
		return
	}
	c.Handlers[t] = handler
}

func (c *Consumers) Start(ctx context.Context, bs []byte) (*packet.Packet, error) {
	pack, err := packet.Parse(bs)
	if err != nil {
		return nil, err
	}

	if h, ok := c.Handlers[pack.Type]; ok {
		return h(ctx, pack)
	}

	return nil, errors.New("无效的请求")
}

func New() *Consumers {
	return &Consumers{}
}
