package notification

import (
	"context"
)

type VerifyData struct {
	Email       string
	VerifyToken string
}

// EmailSender interface
//
//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/email_mock.go
type EmailSender interface {
	SendUserRegistered(ctx context.Context, vd VerifyData) error
}
