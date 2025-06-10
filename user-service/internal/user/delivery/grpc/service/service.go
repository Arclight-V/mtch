package service

import (
	"user-service/internal/user"
)

type usersService struct {
	userUC user.UserUseCase
}

func NewUsersService(userUC user.UserUseCase) *usersService {
	return &usersService{userUC: userUC}
}
