package user

import (
	"context"
	"user-service/internal/models"
)

// User User Case interface
type UserUseCase interface {
	Register(ctx context.Context, in *models.RegisterInput) (*models.RegisterOutput, error)
	VerifyEmail(ctx context.Context, in *models.VerifyEmailInput) (*models.VerifyEmailOutput, error)
}
