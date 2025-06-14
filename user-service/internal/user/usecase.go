package user

import (
	"context"
	"user-service/internal/models"
)

// User User Case interface
type UserUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
}
