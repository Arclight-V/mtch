package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/notification"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
	pb "proto"
	"time"
)

type Interactor struct {
	UserRepo          usecase.UserRepo
	TokenSigner       usecase.TokenSigner
	Hasher            security.PasswordHasher
	PasswordValidator security.PasswordValidator
	EmailSender       notification.EmailSender
}

func (uc *Interactor) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	request := &pb.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	}
	resp, err := uc.UserRepo.Login(ctx, request)
	if err != nil {
		return LoginOutput{}, err
	}

	access, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}
	refresh, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}

	resp.AccessToken = access
	resp.RefreshToken = refresh
	resp.ExpiresIn = int64(time.Minute * 15 / time.Second)
	return LoginOutput{}, nil
}

func (uc *Interactor) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	if err := uc.PasswordValidator.Validate(input.Password); err != nil {
		return RegisterOutput{}, err
	}

	if err := input.SetPassword(input.Password, uc.Hasher); err != nil {
		return RegisterOutput{}, err
	}

	pbRegReq := &pb.RegisterRequest{
		Email:    input.Email,
		Password: input.Password,
	}
	resp, err := uc.UserRepo.Register(ctx, pbRegReq)
	if err != nil {
		return RegisterOutput{}, err
	}

	output := RegisterOutput{
		UserID:    resp.User.Uuid,
		Email:     resp.User.Email,
		CreatedAt: resp.User.CreatedAt.AsTime(),
		Verified:  resp.User.Verified,
	}

	access, err := uc.TokenSigner.SignVerifyToken(output.UserID, 24*time.Hour)
	if err != nil {
		return RegisterOutput{}, err
	}
	output.VerifyToken = access

	if err := uc.EmailSender.SendUserRegistered(ctx, output.Email); err != nil {
		return RegisterOutput{}, err
	}
	return output, nil
}
