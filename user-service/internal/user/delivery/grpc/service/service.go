package service

import (
	pb "proto"
	"user-service/internal/user"
)

type server struct {
	pb.UnimplementedUserInfoServer
}

type usersService struct {
	userUC user.UserUseCase
	server
}

func NewUserServerGRPC(userUC user.UserUseCase) *usersService {
	return &usersService{userUC: userUC}
}
