package user

import (
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
	usecase "user-service/internal/usecase/user"
)

type usersServiceServer struct {
	userservicepb.UnimplementedUserServiceServer

	userUC usecase.UserUseCase
}

func NewUserServiceServer(userUC usecase.UserUseCase) *usersServiceServer {
	return &usersServiceServer{userUC: userUC}
}
