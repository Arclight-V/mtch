package usecase

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/user"
)

type userUseCase struct {
	userRepo user.Repository
}

func NewUserUseCase(userRepo user.Repository) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Register(ctx context.Context, in *models.RegisterInput) (*models.RegisterOutput, error) {
	existUser, err := u.userRepo.FindByContact(ctx, in.PersonalDate.Contact)

	// a User was not found
	if err != nil {
		usr, err2 := u.userRepo.Create(ctx, in)
		if err2 != nil {
			return nil, err2
		}
		return &models.RegisterOutput{UserID: usr.UserID, Status: models.CreatedUnverified}, nil
	}
	if existUser.Activated {
		return &models.RegisterOutput{UserID: existUser.UserID, Status: models.ExistsVerified}, errors.New("user is activated")
	}

	return &models.RegisterOutput{UserID: existUser.UserID, Status: models.ExistsUnverified}, errors.New("user exist, but not activated")

}

func (u *userUseCase) VerifyEmail(ctx context.Context, in *models.VerifyEmailInput) (*models.VerifyEmailOutput, error) {
	// TODO: implement me
	//existUser, err := u.userRepo.FindById(ctx, in.UserID)
	//if err != nil {
	//	return &models.VerifyEmailOutput{}, err
	//}

	return nil, nil
}
