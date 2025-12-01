package user

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	"github.com/Arclight-V/mtch/user-service/internal/features"
)

type userUseCase struct {
	userRepoMem Repository
	userRepoDB  Repository

	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewUserUseCase(
	logger log.Logger,
	featureList *feature_list.FeatureList,
	userRepoMem Repository,
	userRepoDB Repository,
) *userUseCase {
	logger = log.With(logger, "component", "UserUseCase")

	return &userUseCase{
		logger:      logger,
		featureList: featureList,
		userRepoMem: userRepoMem,
		userRepoDB:  userRepoDB,
	}
}

func (u *userUseCase) Register(ctx context.Context, in *domain.RegisterInput) (*domain.RegisterOutput, error) {
	level.Debug(u.logger).Log("msg", "Register", "input", in)

	var (
		existUser *domain.User
		err       error
	)

	if !u.featureList.IsEnabled(features.StoreUsersInDB) {
		existUser, err = u.userRepoMem.FindByContact(ctx, in.PersonalDate.Contact)
	} else {
		existUser, err = u.userRepoDB.FindByContact(ctx, in.PersonalDate.Contact)
	}

	// a User was not found
	if err != nil {
		var (
			usr       *domain.User
			createErr error
		)

		if !u.featureList.IsEnabled(features.StoreUsersInDB) {
			usr, createErr = u.userRepoMem.Create(ctx, in)
		} else {
			usr, createErr = u.userRepoDB.Create(ctx, in)
		}
		if createErr != nil {
			return nil, errors.Wrap(createErr, "user isn't created")
		}

		return &domain.RegisterOutput{UserID: usr.UserID, Status: domain.CreatedUnverified}, nil
	}
	if existUser.Activated {
		return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsVerified}, errors.Wrap(err, "user isn't activated")
	}

	return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsUnverified}, errors.New("user is exist, but not activated")

}

func (u *userUseCase) VerifyEmail(ctx context.Context, in *domain.VerifyEmailInput) (*domain.VerifyEmailOutput, error) {
	// TODO: implement me
	//existUser, err := u.userRepo.FindById(ctx, in.UserID)
	//if err != nil {
	//	return &models.VerifyEmailOutput{}, err
	//}

	return nil, nil
}
