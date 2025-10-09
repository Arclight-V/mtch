package main

import (
	"config"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "proto"
	grpcuser "user-service/internal/adapter/grpc/user"
	repository "user-service/internal/infrastructure/user/repository"
	usecase "user-service/internal/usecase/user"
)

const (
	port = ":50051"
)

func main() {
	cfg, err := config.GetConfig(os.Getenv("user-config"))
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	userRepo := repository.NewUsersDBMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	server := grpcuser.NewUserServerGRPC(userUC)

	s := grpc.NewServer()
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
