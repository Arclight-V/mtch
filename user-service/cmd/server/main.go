package main

import (
	"fmt"
	"github.com/go-kit/log/level"
	"log"
	"net"
	"os"

	"github.com/oklog/run"
	"google.golang.org/grpc"

	"github.com/Arclight-V/mtch/pkg/logging"
	config "github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/Arclight-V/mtch/pkg/prober"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
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

	logger := logging.NewLogger(cfg.LogCfg.Level, cfg.LogCfg.Format, cfg.LogCfg.DebugName)

	var g run.Group

	// Listen for reload signals.
	{
		shutdown := make(chan struct{})
		g.Add(func() error {
			return WaitForInterrupt(shutdown)
		}, func(err error) {
			close(shutdown)
		})
	}

	grpcProbe := prober.NewGRPC()
	httpProbe := prober.NewHTTP()
	statusProber := prober.Combine(grpcProbe, httpProbe)

	level.Debug(logger).Log("msg", "starting HTTP server")
	{
		srv := httpserver.NewServer(logger, httpProbe,
			httpserver.WithListen(cfg.Http.MetricsListenAddr))

		g.Add(func() error {
			statusProber.Healthy()
			statusProber.Ready()

			return srv.ListenAndServe()

		}, func(err error) {
			statusProber.NotReady(err)
			defer statusProber.NotHealthy(err)

			srv.Shutdown(err)
		})
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcserver.NewUnaryServerRequestIDInterceptor()),
	)
	userRepo := repository.NewUsersDBMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	server := grpcuser.NewUserServerGRPC(userUC)
	pb.RegisterUserInfoServer(s, server)

	g.Add(func() error {
		statusProber.Healthy()
		statusProber.Ready()
		lis, err := net.Listen("tcp", cfg.Server.Port)
		if err != nil {
			log.Printf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		return s.Serve(lis)

	}, func(err error) {
		statusProber.NotReady(err)
		defer statusProber.NotHealthy(err)

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
