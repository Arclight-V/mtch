package usecase

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	pb "proto"
)

type UserRepo interface {
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
}

type TokenSigner interface {
	Sign(claims domain.TokenClaims) (string, error)
}
