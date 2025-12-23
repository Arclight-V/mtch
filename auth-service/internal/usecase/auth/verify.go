package auth

import (
	"context"
	"time"
)

type VerifyInput struct {
	Code string
}

type VerifyOutput struct {
	VerifiedAt time.Time
	UserID     string
	Verified   bool
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/verify_email_mock.go
type VerifyUseCase interface {
	VerifyCode(ctx context.Context, in *VerifyInput) (*VerifyOutput, error)
}
