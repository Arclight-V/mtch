package user

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kit/log/level"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

func (s *usersServiceServer) Register(ctx context.Context, req *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error) {
	level.Debug(s.logger).Log("method", "Register", "req", req)

	pd := &domain.PersonalData{
		FirstName: req.PersonalData.FirstName,
		LastName:  req.PersonalData.LastName,
		Contact:   req.PersonalData.Contact,
		Password:  req.PersonalData.Password,
	}
	pd.SetDateBirthday(
		int(req.PersonalData.BirthDate.BirthYear),
		int(req.PersonalData.BirthDate.BirthMonth),
		int(req.PersonalData.BirthDate.BirthDay))

	in := &domain.RegisterInput{PersonalDate: pd}

	user, err := s.userUC.Register(ctx, in)
	if err != nil {
		return nil, err
	}
	resp := &userservicepb.RegisterResponse{UserId: user.UserID.String(), Status: userservicepb.CreateUserStatus(user.Status)}

	return resp, nil
}

func (s *usersServiceServer) VerifyEmail(ctx context.Context, req *userservicepb.VerifyEmailRequest) (*userservicepb.VerifyEmailResponse, error) {
	fmt.Println("VerifyEmail called")

	in := &domain.VerifyEmailInput{
		UserID: req.Uuid,
	}

	out, err := s.userUC.VerifyEmail(ctx, in)
	if err != nil {
		return &userservicepb.VerifyEmailResponse{}, err
	}

	response := &userservicepb.VerifyEmailResponse{VerifiedAt: timestamppb.New(out.VerifiedAt), Verified: out.Verified}
	return response, nil
}

func (s *usersServiceServer) Login(ctx context.Context, req *userservicepb.LoginRequest) (*userservicepb.LoginResponse, error) {
	log.Println("Login called", req.Email, req.Password)

	resp := &userservicepb.LoginResponse{
		User: &userservicepb.User{
			Uuid:      "uuid",
			FirstName: "first_name",
			LastName:  "last_name",
			Email:     "email",
		},
		SessionId: "session_id",
	}

	hash_password_resp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// TODO: temporary password, replace with a hash from the database
	tmp_password := "password"
	if err := bcrypt.CompareHashAndPassword(hash_password_resp, []byte(tmp_password)); err != nil {
		return nil, err
	}

	return resp, nil
}
