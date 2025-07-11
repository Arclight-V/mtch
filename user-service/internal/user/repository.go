package user

import (
	"context"
	"user-service/internal/models"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/repository_mock.go
type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
