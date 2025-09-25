package usecase

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/user"
)

type userUseCase struct {
	userRepo user.Repository
}

func NewUserUseCase(userRepo user.Repository) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Register(ctx context.Context, user *models.RegistrationData) (*models.User, error) {
	//TODO implement me

	existUser, err := u.userRepo.FindByEmail(ctx, user.Email)
	// a User was not found
	if err != nil {
		return u.userRepo.Create(ctx, user)
	}
	if existUser.Verified {
		return nil, errors.New("User already registered.")
	} else {
		// TODO: add logic for re-verification by email change pashHas, etc.
	}

	return u.userRepo.Create(ctx, user)
}
