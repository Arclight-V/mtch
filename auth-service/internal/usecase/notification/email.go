package notification

import "context"

type EmailSender interface {
	SendUserRegistered(ctx context.Context, to string) error
}
