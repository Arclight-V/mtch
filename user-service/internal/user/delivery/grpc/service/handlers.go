package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	pb "proto"
	"user-service/internal/models"
)

func (s *usersService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Println("Register called:", req.Email)
	pendingUser, err := models.NewPendingUser(req.GetEmail(), req.GetPassword())
	if err != nil {
		return &pb.RegisterResponse{}, err
	}
	user, err := s.userUC.Register(ctx, pendingUser)
	if err != nil {
		return &pb.RegisterResponse{}, err
	}
	pbUser := userModelToProto(user)
	return &pb.RegisterResponse{User: pbUser}, nil
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
