package main

import (
	"context"
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
			Password:  "password",
			Email:     "email",
		},
		SessionId: "session_id",
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
