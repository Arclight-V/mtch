package user

import (
	"context"
	"errors"

	domain "user-service/internal/domain/user"
)

type userUseCase struct {
	userRepo Repository
}

func NewUserUseCase(userRepo Repository) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Register(ctx context.Context, in *domain.RegisterInput) (*domain.RegisterOutput, error) {
	existUser, err := u.userRepo.FindByContact(ctx, in.PersonalDate.Contact)

	// a User was not found
	if err != nil {
		usr, err2 := u.userRepo.Create(ctx, in)
		if err2 != nil {
			return nil, err2
		}
		return &domain.RegisterOutput{UserID: usr.UserID, Status: domain.CreatedUnverified}, nil
	}
	if existUser.Activated {
		return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsVerified}, errors.New("userservice is activated")
	}

	return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsUnverified}, errors.New("userservice exist, but not activated")

}

func (u *userUseCase) VerifyEmail(ctx context.Context, in *domain.VerifyEmailInput) (*domain.VerifyEmailOutput, error) {
	// TODO: implement me
	//existUser, err := u.userRepo.FindById(ctx, in.UserID)
	//if err != nil {
	//	return &models.VerifyEmailOutput{}, err
	//}

	return nil, nil
}
