package user

import (
	"context"
	"user-service/internal/models"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
}
