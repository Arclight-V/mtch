package messagebroker

import (
	"context"
)

type Event struct {
	Topic   string
	Key     []byte
	Value   []byte
	Headers map[string][]byte
	Type    string
}

type Handler func(ctx context.Context, e Event) error

type Publisher interface {
	Publish(ctx context.Context, event *Event) error
	Close() error
}

type Consumer interface {
	Start(ctx context.Context, h Handler) error
	Close() error
}
