package usecase

import (
	"context"
	"log"
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
	log.Println("//TODO implement me")
	log.Println(user)
	return user, nil
}
