package auth

import (
	"context"
	"time"
)

type VerifyEmailInput struct {
	Token string
}

type VerifyEmailOutput struct {
	UserID     string
	VerifiedAt time.Time
	Verified   bool
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/verify_email_mock.go
type VerifyEmailUseCase interface {
	VerifyEmail(ctx context.Context, in VerifyEmailInput) (VerifyEmailOutput, error)
}
