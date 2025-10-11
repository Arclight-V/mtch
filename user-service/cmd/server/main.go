package main

import (
	"config"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/oklog/run"

	"github.com/Arclight-V/mtch/pkg/prober"
	"github.com/Arclight-V/mtch/pkg/signaler"

	pb "proto"
	grpcuser "user-service/internal/adapter/grpc/user"
	"user-service/internal/infrastructure/user/repository"
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

	var g run.Group

	userRepo := repository.NewUsersDBMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	server := grpcuser.NewUserServerGRPC(userUC)

	s := grpc.NewServer()
	grpcProbe := prober.NewGRPC()
	statusProber := prober.Combine(grpcProbe)

	pb.RegisterUserInfoServer(s, server)

	// Listen for reload signals.
	{
		shutdown := make(chan struct{})
		g.Add(func() error {
			return WaitForInterrupt(shutdown)
		}, func(err error) {
			close(shutdown)
		})
	}

	g.Add(func() error {
		statusProber.Ready()
		lis, err := net.Listen("tcp", cfg.Server.Port)
		if err != nil {
			log.Printf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		return s.Serve(lis)

	}, func(err error) {
		statusProber.NotReady(err)
		s.GracefulStop()
	})

	if err := g.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
		os.Exit(1)
	}
	log.Println("Shutting down")
}

func WaitForInterrupt(cancel <-chan struct{}) error {
	interrupt := signaler.WaitForInterrupt()
	select {
	case s := <-interrupt:
		log.Printf("received signal: %v", s)
		return nil
	case <-cancel:
		return fmt.Errorf("Captured %v, shutdown requested.\n", interrupt)
	}
}
