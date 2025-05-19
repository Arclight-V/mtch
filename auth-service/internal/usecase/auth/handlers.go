package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
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

	access, err := uc.TokenSigner.Sign(domain.TokenClaims{UserId: resp.User.Uuid, Role: "user", Exp: time.Now().Add(time.Minute * 15)})
	if err != nil {
		return nil, err
	}
	refresh, err := uc.TokenSigner.Sign(domain.TokenClaims{UserId: resp.User.Uuid, Role: "user", Exp: time.Now().Add(time.Hour * 24)})
	if err != nil {
		return nil, err
	}

	resp.AccessToken = access
	resp.RefreshToken = refresh
	resp.ExpiresIn = int64(time.Minute * 15 / time.Second)
	return resp, nil
}
