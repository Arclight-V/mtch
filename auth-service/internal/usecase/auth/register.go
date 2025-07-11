package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
	"time"
)

type RegisterInput struct {
	Email    string
	Password string
}

func (ri *RegisterInput) SetPassword(plain string, h security.PasswordHasher) error {
	hash, err := h.Hash(plain)
	if err != nil {
		return err
	}
	ri.Password = hash
	return nil
}

type RegisterOutput struct {
	UserID      string
	Email       string
	CreatedAt   time.Time
	Verified    bool
	VerifyToken string
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/register_mock.go
type RegisterUseCase interface {
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)
}
