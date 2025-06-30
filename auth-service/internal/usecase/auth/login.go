package auth

import "context"

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	UserID   string
	Verified bool
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/login_mock.go
type LoginUseCase interface {
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
}
