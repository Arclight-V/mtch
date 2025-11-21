package usecase

import (
	"context"
	"time"

	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../usecase/mocks/ports_mock.go
type UserRepo interface {
	Login(ctx context.Context, request *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error)
	Register(ctx context.Context, request *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error)
	VerifyEmail(ctx context.Context, request *userservicepb.VerifyEmailRequest) (*userservicepb.VerifyEmailResponse, error)

	NotifyUserRegistered(ctx context.Context, request *notificationservicepb.NotificationUserContactsRequest) (*notificationservicepb.NotificationUserContactsResponse, error)
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

type Metrics interface {
	IncLoginAttempts(status string)
	ObserveLoginDuration(status string, s float64)
}
