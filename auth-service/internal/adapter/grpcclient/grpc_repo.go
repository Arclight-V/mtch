package grpcclient

import (
	"context"
	pb "proto"
)

type GrpcUserRepo struct {
	cli pb.UserInfoClient
}

func NewGRPCUserRepo(cli pb.UserInfoClient) *GrpcUserRepo {
	return &GrpcUserRepo{cli: cli}
}

func (r *GrpcUserRepo) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := r.cli.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
