package usecase

import (
	"context"
	"user-service/internal/models"
	"user-service/internal/user"
)

type userUseCase struct {
	userRepo user.Repository
}

func NewUserUseCase(userRepo user.Repository) *userUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) Register(ctx context.Context, user *models.RegistrationData) (*models.RegistrationOutput, error) {
	existUser, err := u.userRepo.FindByEmail(ctx, user.Email)

	// a User was not found
	if err != nil {
		usr, err2 := u.userRepo.Create(ctx, user)
		if err2 != nil {
			return nil, err2
		}
		return &models.RegistrationOutput{UserID: usr.UserID, Status: models.CreatedUnverified}, nil
	}
	if existUser.Verified {
		return &models.RegistrationOutput{UserID: existUser.UserID, Status: models.ExistsVerified}, nil
	}

	return &models.RegistrationOutput{UserID: existUser.UserID, Status: models.ExistsUnverified}, nil

}

func (u *userUseCase) VerifyEmail(ctx context.Context, in *models.VerifyEmailInput) (*models.VerifyEmailOutput, error) {
	// TODO: implement me
	//existUser, err := u.userRepo.FindById(ctx, in.UserID)
	//if err != nil {
	//	return &models.VerifyEmailOutput{}, err
	//}

	return nil, nil
}
