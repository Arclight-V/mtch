package user

import (
	pb "proto"
	usecase "user-service/internal/usecase/user"
)

type server struct {
	pb.UnimplementedUserInfoServer
}

type usersService struct {
	userUC usecase.UserUseCase
	server
}

func NewUserServerGRPC(userUC usecase.UserUseCase) *usersService {
	return &usersService{userUC: userUC}
}
