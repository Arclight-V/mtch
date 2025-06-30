package auth

import (
	"context"
	"time"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID    string
	Email     string
	CreatedAt time.Time
	Verified  bool
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/register_mock.go
type RegisterUseCase interface {
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)
}
