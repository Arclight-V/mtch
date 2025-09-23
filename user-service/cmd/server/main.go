package main

import (
	"config"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	pb "proto"
	userServerGRPC "user-service/internal/user/delivery/grpc/service"
	"user-service/internal/user/repository"
	userUserCase "user-service/internal/user/usecase"
)

const (
	port = ":50051"
)

func main() {
	cfg, err := config.GetConfig(os.Getenv("user-config"))
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	s := grpc.NewServer()
	userRepo := repository.NewUserRepository()
	userUC := userUserCase.NewUserUseCase(userRepo)
	server := userServerGRPC.NewUserServerGRPC(userUC)
	pb.RegisterUserInfoServer(s, server)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
