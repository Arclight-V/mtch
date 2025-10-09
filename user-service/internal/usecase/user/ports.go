package user

import (
	"context"

	domain "user-service/internal/domain/user"
)

// User User Case interface
//
//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/ports_mock.go
type UserUseCase interface {
	Register(ctx context.Context, in *domain.RegisterInput) (*domain.RegisterOutput, error)
	VerifyEmail(ctx context.Context, in *domain.VerifyEmailInput) (*domain.VerifyEmailOutput, error)
}
