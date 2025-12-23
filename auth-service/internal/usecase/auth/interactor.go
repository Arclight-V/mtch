package auth

import (
	"context"
	"github.com/go-kit/log/level"
	"time"

	"github.com/go-kit/log"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/messagebroker"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
)

type AuthUseCase struct {
	userRepo usecase.UserRepo
	//Metrics           *prometheus.Registry
	tokenSigner       usecase.TokenSigner
	hasher            security.PasswordHasher
	passwordValidator security.PasswordValidator
	publisher         messagebroker.Publisher

	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewAuthUseCase(
	logger log.Logger,
	featureList *feature_list.FeatureList,
	userRepo usecase.UserRepo,
	tokenSigner usecase.TokenSigner,
	hasher security.PasswordHasher,
	passwordValidator security.PasswordValidator,
	publisher messagebroker.Publisher,
) *AuthUseCase {
	return &AuthUseCase{
		logger:            logger,
		featureList:       featureList,
		userRepo:          userRepo,
		tokenSigner:       tokenSigner,
		hasher:            hasher,
		passwordValidator: passwordValidator,
		publisher:         publisher,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	request := &userservicepb.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	}
	resp, err := uc.userRepo.Login(ctx, request)
	if err != nil {
		return LoginOutput{}, err
	}

	access, err := uc.tokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}
	refresh, err := uc.tokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}

	resp.AccessToken = access
	resp.RefreshToken = refresh
	resp.ExpiresIn = int64(time.Minute * 15 / time.Second)
	return LoginOutput{}, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	if err := uc.passwordValidator.Validate(in.Password); err != nil {
		return nil, err
	}

	if err := in.SetPassword(in.Password, uc.hasher); err != nil {
		return nil, err
	}
	if err := in.SetEmailOrPhone(); err != nil {
		return nil, err
	}

	pbRegReq := &userservicepb.RegisterRequest{
		PersonalData: &userservicepb.PersonalData{
			FirstName: in.FirstName,
			LastName:  in.LastName,
			Contact:   in.Contact,
			Password:  in.Password,
			BirthDate: &userservicepb.Date{
				BirthDay:   in.Date.BirthDay,
				BirthMonth: in.Date.BirthMonth,
				BirthYear:  in.Date.BirthYear,
			},
		},
	}

	resp, err := uc.userRepo.Register(ctx, pbRegReq)
	if err != nil {
		return nil, err
	}

	output := &RegisterOutput{
		UserID: resp.UserId,
		Email:  in.Contact,
	}

	// Evaluate kafka-enable feature flag
	kafkaEnable := uc.featureList.IsEnabled(feature_list.FeatureKafka)
	if kafkaEnable {
		event := messagebroker.Event{
			Topic: "notifications.request.v1",
			// TODO: when there is a practical need for orderliness or idempotence.
			// Key:
			Value: []byte(in.Contact),
			Headers: map[string][]byte{
				"event-type":    []byte("verification"),
				"event-version": []byte("1"),
				"content-type":  []byte("application/json"),
			},
		}
		if err := uc.publisher.Publish(context.TODO(), &event); err != nil {
			return nil, err
		}
		return output, nil
	}

	contacts := []*notificationservicepb.Contact{{Chanel: notificationservicepb.Channel_ChannelEmail, Value: output.Email}}
	notifyReq := &notificationservicepb.NotificationUserContactsRequest{UserID: output.UserID, Contacts: contacts}

	if _, err := uc.userRepo.NotifyUserRegistered(ctx, notifyReq); err != nil {
		level.Error(uc.logger).Log("msg", "failed to notify user", "err", err)
	}

	return output, nil
}

func (uc *AuthUseCase) VerifyCode(ctx context.Context, in *VerifyInput) (*VerifyOutput, error) {

	pbResp := &userservicepb.VerifyRequest{
		Code: in.Code,
	}

	resp, err := uc.userRepo.VerifyCode(ctx, pbResp)
	if err != nil {
		return nil, err
	}

	output := &VerifyOutput{
		VerifiedAt: resp.VerifiedAt.AsTime(),
		Verified:   resp.Verified,
	}
	return output, nil
}
