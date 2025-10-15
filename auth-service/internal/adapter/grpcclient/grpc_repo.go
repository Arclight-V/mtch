package grpcclient

import (
	"context"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

type GrpcUserRepo struct {
	cli userservicepb.UserServiceClient
}

func NewGRPCUserRepo(cli userservicepb.UserServiceClient) *GrpcUserRepo {
	return &GrpcUserRepo{cli: cli}
}

func (r *GrpcUserRepo) Login(ctx context.Context, request *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	resp, err := r.cli.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *GrpcUserRepo) Register(ctx context.Context, request *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error) {
	resp, err := r.cli.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *GrpcUserRepo) VerifyEmail(ctx context.Context, request *userservicepb.VerifyEmailRequest) (*userservicepb.VerifyEmailResponse, error) {
	//TODO implement me
	resp, err := r.cli.VerifyEmail(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
