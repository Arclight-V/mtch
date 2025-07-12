package email

import (
	"context"
	"log"
)

type NoopSender struct {
}

func NewNoopSender() *NoopSender { return new(NoopSender) }

func (n *NoopSender) SendUserRegistered(ctx context.Context, to string) error {
	log.Printf("sending user registered to: %v", to)
	return nil
}
