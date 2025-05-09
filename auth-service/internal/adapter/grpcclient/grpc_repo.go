package grpcclient

import (
	"context"
	"log"
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
	log.Println(resp, "19")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
