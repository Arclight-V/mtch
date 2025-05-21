package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	pb "proto"
	"time"
)

type Interactor struct {
	UserRepo    usecase.UserRepo
	TokenSigner usecase.TokenSigner
}

func (uc *Interactor) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := uc.UserRepo.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	access, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return nil, err
	}
	refresh, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return nil, err
	}

	resp.AccessToken = access
	resp.RefreshToken = refresh
	resp.ExpiresIn = int64(time.Minute * 15 / time.Second)
	return resp, nil
}
