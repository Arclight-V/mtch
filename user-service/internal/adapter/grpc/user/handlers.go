package user

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"log"

	"github.com/go-kit/log/level"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	usecase "github.com/Arclight-V/mtch/user-service/internal/usecase/user"
)

func (s *usersServiceServer) Register(ctx context.Context, req *userservicepb.RegisterRequest) (*userservicepb.RegisterResponse, error) {
	level.Debug(s.logger).Log("method", "Register", "req", req)

	pd := domain.NewPersonalDataFromRegisterRequest(req)
	in := &domain.RegisterInput{PersonalDate: pd}

	user, err := s.userUC.Register(ctx, in)
	resp := &userservicepb.RegisterResponse{UserId: user.UserID.String(), Status: userservicepb.CreateUserStatus(user.Status)}
	if err != nil {
		if errors.Is(err, usecase.ErrUserIsExistUnverified) || errors.Is(err, usecase.ErrUserIsExist) {
			err = status.Error(codes.AlreadyExists, err.Error())
		} else {
			err = status.Error(codes.Internal, err.Error())
		}
		return resp, err
	}

	return resp, nil
}

func (s *usersServiceServer) VerifyEmail(ctx context.Context, req *userservicepb.VerifyRequest) (*userservicepb.VerifyResponse, error) {
	in := &domain.VerifyInput{
		UserID: req.UserId,
		Code:   req.Code,
	}

	out, err := s.userUC.VerifyEmail(ctx, in)
	if err != nil {
		return nil, err
	}

	response := &userservicepb.VerifyResponse{
		UserId:     out.UserID,
		VerifiedAt: timestamppb.New(out.VerifiedAt),
		Verified:   out.Verified,
	}
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
