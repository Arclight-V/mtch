package user

import (
	"context"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

// UserUserCase interface
//
//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/ports_mock.go
type UserUseCase interface {
	Register(ctx context.Context, in *domain.RegisterInput) (*domain.RegisterOutput, error)
	VerifyEmail(ctx context.Context, in *domain.VerifyInput) (*domain.VerifyOutput, error)
}
