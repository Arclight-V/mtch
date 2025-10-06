package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/notification"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
	"log"
	pb "proto"
	"time"
)

type Interactor struct {
	UserRepo          usecase.UserRepo
	TokenSigner       usecase.TokenSigner
	Hasher            security.PasswordHasher
	PasswordValidator security.PasswordValidator
	EmailSender       notification.EmailSender
	VerifyTokenRepo   usecase.VerifyTokenRepo
}

func (uc *Interactor) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	request := &pb.LoginRequest{
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

	pbRegReq := &pb.RegisterRequest{
		PersonalData: &pb.PersonalData{
			FirstName: in.FirstName,
			LastName:  in.LastName,
			Contact:   in.Contact,
			Password:  in.Password,
			BirthDate: &pb.Date{
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
	log.Printf("user registered to: %v", output)

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

	pbResp := &pb.VerifyEmailRequest{
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
