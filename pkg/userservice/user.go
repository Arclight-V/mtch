package userservice

import (
	"google.golang.org/grpc"

	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

func RegisterUserServer(userSrv userservicepb.UserServiceServer) func(*grpc.Server) {
	return func(s *grpc.Server) {
		userservicepb.RegisterUserServiceServer(s, userSrv)
	}
}
