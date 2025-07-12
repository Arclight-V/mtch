package usecase

import (
	"context"
	"user-service/internal/models"
	"user-service/internal/user"
)

type userUseCase struct {
	userRepo user.Repository
}

func NewUserUseCase(userRepo user.Repository) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	//TODO implement me

	existUser, err := u.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	return existUser, nil
}
