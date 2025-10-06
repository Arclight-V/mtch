package user

import (
	"context"
	"user-service/internal/models"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/repository_mock.go
type Repository interface {
	Create(ctx context.Context, regData *models.RegisterInput) (*models.User, error)
	FindByContact(ctx context.Context, contact string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
}
