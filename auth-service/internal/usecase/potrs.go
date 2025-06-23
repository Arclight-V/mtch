package usecase

import (
	"context"
	pb "proto"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../usecase/mocks/ports_mock.go
type UserRepo interface {
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error)
}

type TokenSigner interface {
	SignAccess(uuid, sid string) (string, error)
	SignRefresh(uuid, sid string) (string, string, error)
}
