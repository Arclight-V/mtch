package usecase

import (
	"context"
	pb "proto"
)

type UserRepo interface {
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
}
