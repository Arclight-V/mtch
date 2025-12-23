package grpcclient

import (
	"context"

	"github.com/go-kit/log"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

type GrpcUserRepo struct {
	userClient         userservicepb.UserServiceClient
	notificationClient notificationservicepb.NotificationServiceClient

	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewGRPCUserRepo(
	userClient userservicepb.UserServiceClient,
	notificationClient notificationservicepb.NotificationServiceClient,
	logger log.Logger,
	featureList *feature_list.FeatureList,
) *GrpcUserRepo {
	logger = log.With(logger, "component", "gRPC/GrpcUserClient")

	return &GrpcUserRepo{
		userClient:         userClient,
		notificationClient: notificationClient,
		logger:             logger,
		featureList:        featureList,
	}
}

func (r *GrpcUserRepo) Login(ctx context.Context, request *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	resp, err := r.userClient.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *GrpcUserRepo) Register(ctx context.Context, request *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error) {
	resp, err := r.userClient.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *GrpcUserRepo) VerifyCode(ctx context.Context, request *userservicepb.VerifyRequest) (*userservicepb.VerifyResponse, error) {
	//TODO implement me
	resp, err := r.userClient.VerifyCode(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *GrpcUserRepo) NotifyUserRegistered(ctx context.Context, request *notificationservicepb.NotificationUserContactsRequest) (*notificationservicepb.NotificationUserContactsResponse, error) {
	resp, err := r.notificationClient.NotifyUserRegistered(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
