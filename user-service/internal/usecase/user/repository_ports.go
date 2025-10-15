package user

import (
	"context"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/repository_mock.go
type Repository interface {
	Create(ctx context.Context, regData *domain.RegisterInput) (*domain.User, error)
	FindByContact(ctx context.Context, contact string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
}
