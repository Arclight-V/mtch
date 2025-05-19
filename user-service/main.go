package main

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "proto"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedUserInfoServer
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Println("Register called:", req.FirstName)
	return nil, nil
}
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserInfoServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
