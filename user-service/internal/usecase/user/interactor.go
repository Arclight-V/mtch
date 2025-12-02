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

var (
	ErrUserIsExist           = errors.New("user is exist")
	ErrUserIsExistUnverified = errors.New("user is exist but is unverified")
	ErrUserNotCreated        = errors.New("user isn't created")
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

	if u.featureList.IsEnabled(features.StoreUsersInDB) {
		existUser, err = u.userRepoDB.Create(ctx, in)
	} else {
		existUser, err = u.userRepoMem.Create(ctx, in)
	}

	if err != nil && existUser != nil {
		switch existUser.Activated {
		case true:
			return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsVerified}, ErrUserIsExistUnverified
		case false:
			return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.ExistsUnverified}, ErrUserIsExist
		}
	} else if err != nil {
		return nil, errors.Wrap(err, ErrUserNotCreated.Error())
	}

	return &domain.RegisterOutput{UserID: existUser.UserID, Status: domain.CreatedUnverified}, nil

}

func (u *userUseCase) VerifyEmail(ctx context.Context, in *domain.VerifyEmailInput) (*domain.VerifyEmailOutput, error) {
	// TODO: implement me
	//existUser, err := u.userRepo.FindById(ctx, in.UserID)
	//if err != nil {
	//	return &models.VerifyEmailOutput{}, err
	//}

	return nil, nil
}
