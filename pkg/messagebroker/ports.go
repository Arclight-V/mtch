package messagebroker

import (
	"context"
)

// Event defines the structure of a message exchanged through the message broker.
// It contains metadata such as Topic, Key, and Headers, as well as the message payload in Value.
type Event struct {
	Topic   string
	Key     []byte
	Value   []byte
	Headers map[string][]byte
}

// Handler represents a function that processes an incoming Event.
// It receives a context for cancellation and deadline control.
type Handler func(ctx context.Context, e Event) error

// Publisher defines the interface for producing messages to a broker.
// Implementations should handle connection management and message delivery.
type Publisher interface {
	Publish(ctx context.Context, event *Event) error
	Close() error
}

// Consumer defines the interface for consuming messages from a broker.
// Implementations are expected to manage subscription and message handling loops.
type Consumer interface {
	Start(ctx context.Context, h Handler) error
	Close() error
}
