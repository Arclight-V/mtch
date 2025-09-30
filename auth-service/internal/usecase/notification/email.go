package notification

import (
	"context"
)

type VerifyData struct {
	Email       string
	VerifyToken string
}

type EmailSender interface {
	SendUserRegistered(ctx context.Context, vd VerifyData) error
}
