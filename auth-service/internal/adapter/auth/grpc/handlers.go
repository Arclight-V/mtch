package grpc

import (
	"context"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

// Register user registration
func (r *AuthClient) Register(ctx context.Context, request *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error) {
	resp, err := r.userClient.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NotifyUserRegistered Sending information about registration in the system to the user.
//
// It can be used as a kafka alternative.
func (r *AuthClient) NotifyUserRegistered(ctx context.Context, request *notificationservicepb.NotificationUserContactsRequest) (*notificationservicepb.NotificationUserContactsResponse, error) {
	resp, err := r.notificationClient.NotifyUserRegistered(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyCode verification of user codes. New users, password recovery
func (r *AuthClient) VerifyCode(ctx context.Context, request *userservicepb.VerifyRequest) (*userservicepb.VerifyResponse, error) {
	//TODO implement me
	resp, err := r.userClient.VerifyCode(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Login user's login
func (r *AuthClient) Login(ctx context.Context, request *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	resp, err := r.userClient.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
