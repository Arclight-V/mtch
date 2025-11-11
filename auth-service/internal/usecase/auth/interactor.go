package auth

import (
	"context"
	"log"
	"time"

	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/notification"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
	"github.com/Arclight-V/mtch/pkg/messagebroker"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

type Interactor struct {
	UserRepo usecase.UserRepo
	//Metrics           *prometheus.Registry
	TokenSigner       usecase.TokenSigner
	Hasher            security.PasswordHasher
	PasswordValidator security.PasswordValidator
	EmailSender       notification.EmailSender
	VerifyTokenRepo   usecase.VerifyTokenRepo
	Publisher         messagebroker.Publisher
}

func (uc *Interactor) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	request := &userservicepb.LoginRequest{
		Email:    input.Email,
		Password: input.Password,
	}
	resp, err := uc.UserRepo.Login(ctx, request)
	if err != nil {
		return LoginOutput{}, err
	}

	access, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}
	refresh, err := uc.TokenSigner.SignAccess(resp.User.Uuid, resp.SessionId)
	if err != nil {
		return LoginOutput{}, err
	}

	resp.AccessToken = access
	resp.RefreshToken = refresh
	resp.ExpiresIn = int64(time.Minute * 15 / time.Second)
	return LoginOutput{}, nil
}

func (uc *Interactor) Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error) {
	if err := uc.PasswordValidator.Validate(in.Password); err != nil {
		return nil, err
	}

	if err := in.SetPassword(in.Password, uc.Hasher); err != nil {
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

	resp, err := uc.UserRepo.Register(ctx, pbRegReq)
	if err != nil {
		return nil, err
	}

	output := &RegisterOutput{
		UserID: resp.UserId,
		Email:  in.Contact,
	}
	log.Printf("userservice registered to: %v", output)

	verifyTokenIssue, token, err := uc.TokenSigner.SignVerifyToken(output.UserID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	vd := notification.VerifyData{
		Email:       in.Contact,
		VerifyToken: token,
	}

	if err := uc.EmailSender.SendUserRegistered(ctx, vd); err != nil {
		return nil, err
	}
	event := messagebroker.Event{
		Topic: "notifications.request.v1",
		// TODO: when there is a practical need for orderliness or idempotence.
		// Key:
		Value: []byte(vd.Email),
		Headers: map[string][]byte{
			"event-type":    []byte("verification"),
			"event-version": []byte("1"),
			"content-type":  []byte("application/json"),
		},
	}
	if err := uc.Publisher.Publish(context.TODO(), &event); err != nil {
		return nil, err
	}
	if err := uc.VerifyTokenRepo.InsertIssue(ctx, verifyTokenIssue); err != nil {
		return nil, err
	}

	return output, nil
}

func (uc *Interactor) VerifyEmail(ctx context.Context, in VerifyEmailInput) (VerifyEmailOutput, error) {
	v, err := uc.TokenSigner.ParseVerifyToken(in.Token)
	if err != nil {
		return VerifyEmailOutput{}, err
	}

	if err := uc.VerifyTokenRepo.TryConsumeJTI(ctx, v); err != nil {
		return VerifyEmailOutput{}, err
	}

	pbResp := &userservicepb.VerifyEmailRequest{
		Uuid: v.UserID,
	}

	resp, err := uc.UserRepo.VerifyEmail(ctx, pbResp)
	if err != nil {
		return VerifyEmailOutput{}, err
	}

	output := VerifyEmailOutput{
		UserID:     v.UserID,
		VerifiedAt: resp.VerifiedAt.AsTime(),
		Verified:   resp.Verified,
	}
	return output, nil
}
