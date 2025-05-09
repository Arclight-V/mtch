package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	pb "proto"
)

type Interactor struct {
	UserRepo usecase.UserRepo
}

func (uc *Interactor) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := uc.UserRepo.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	hash_password_resp, err := bcrypt.GenerateFromPassword([]byte(resp.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(hash_password_resp, []byte(request.Password)); err != nil {
		return nil, err
	}
	return resp, nil
}
