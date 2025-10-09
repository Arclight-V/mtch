package usecase

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	pb "proto"
	"time"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../usecase/mocks/ports_mock.go
type UserRepo interface {
	Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error)
	VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error)
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../usecase/mocks/token_signer_mock.go
type TokenSigner interface {
	SignAccess(uuid, sid string) (string, error)
	SignRefresh(uuid, sid string) (string, string, error)
	SignVerifyToken(uuid string, ttl time.Duration) (domain.VerifyTokenIssue, string, error)
	ParseVerifyToken(tokenStr string) (domain.VerifyEmailToken, error)
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../usecase/mocks/verify_token_repo_mock.go
type VerifyTokenRepo interface {
	InsertIssue(ctx context.Context, v domain.VerifyTokenIssue) error
	TryConsumeJTI(ctx context.Context, v domain.VerifyEmailToken) error
}
