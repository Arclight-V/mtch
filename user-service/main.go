package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "proto"
	userServerGRPC "user-service/internal/user/delivery/grpc/service"
	"user-service/internal/user/repository"
	userUserCase "user-service/internal/user/usecase"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	userRepo := repository.NewUserRepository()
	userUC := userUserCase.NewUserUseCase(userRepo)

	server := userServerGRPC.NewUserServerGRPC(userUC)

	pb.RegisterUserInfoServer(s, server)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
