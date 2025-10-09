package user

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "proto"
	domain "user-service/internal/domain/user"
)

func (s *usersService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Println("Register called:", req.PersonalData.Contact)

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
		return &pb.RegisterResponse{}, err
	}
	resp := &pb.RegisterResponse{UserId: user.UserID.String(), Status: pb.CreateUserStatus(user.Status)}

	return resp, nil
}

func (s *usersService) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	fmt.Println("VerifyEmail called")

	in := &domain.VerifyEmailInput{
		UserID: req.Uuid,
	}

	out, err := s.userUC.VerifyEmail(ctx, in)
	if err != nil {
		return &pb.VerifyEmailResponse{}, err
	}

	response := &pb.VerifyEmailResponse{VerifiedAt: timestamppb.New(out.VerifiedAt), Verified: out.Verified}
	return response, nil
}

func (s *usersService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Println("Login called", req.Email, req.Password)

	resp := &pb.LoginResponse{
		User: &pb.User{
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
