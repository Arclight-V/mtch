package main

import (
	"log"
	"os"
	"regexp"

	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"

	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/Arclight-V/mtch/pkg/prober"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
	"github.com/Arclight-V/mtch/pkg/signaler"
	"github.com/Arclight-V/mtch/pkg/userservice"

	grpcuser "github.com/Arclight-V/mtch/user-service/internal/adapter/grpc/user"
	"github.com/Arclight-V/mtch/user-service/internal/infrastructure/user/repository"
	usecase "github.com/Arclight-V/mtch/user-service/internal/usecase/user"
)

const (
	port = ":50051"
)

func main() {
	cfg, err := config.GetConfig(os.Getenv("USER_CONFIG"))
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := logging.NewLogger(cfg.LogCfg.Level, cfg.LogCfg.Format, cfg.LogCfg.DebugName)

	metrics := prometheus.NewRegistry()
	metrics.MustRegister(
		versioncollector.NewCollector("mtch-user-service"),
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	prometheus.DefaultRegisterer = metrics

	var g run.Group

	// Listen for reload signals.
	{
		shutdown := make(chan struct{})
		g.Add(func() error {
			return signaler.WaitForInterrupt(shutdown)
		}, func(err error) {
			close(shutdown)
		})
	}

	grpcProbe := prober.NewGRPC()
	httpProbe := prober.NewHTTP()
	statusProber := prober.Combine(grpcProbe, httpProbe)

	level.Debug(logger).Log("msg", "starting HTTP server")
	{
		srv := httpserver.NewServer(logger, metrics, httpProbe,
			httpserver.WithListen(cfg.Http.MetricsListenAddr))

		g.Add(func() error {
			statusProber.Healthy()

			return srv.ListenAndServe()
		}, func(err error) {
			statusProber.NotReady(err)
			defer statusProber.NotHealthy(err)

			srv.Shutdown(err)
		})
	}

	userRepo := repository.NewUsersDBMemory()
	userUC := usecase.NewUserUseCase(userRepo)
	server := grpcuser.NewUserServiceServer(userUC)

	level.Debug(logger).Log("msg", "starting GRPC server")
	{
		s := grpcserver.NewServer(logger, metrics, grpcProbe,
			grpcserver.WithServer(userservice.RegisterUserServer(server)),
			grpcserver.WithListen(cfg.UserServiceServer.Port),
			grpcserver.WithGracePeriod(cfg.UserServiceServer.GracePeriod),
		)

		g.Add(func() error {
			statusProber.Ready()

			return s.ListenAndServe()
		}, func(err error) {
			statusProber.NotReady(err)

			s.Shutdown(err)
		})
	}

	if err := g.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
		os.Exit(1)
	}
	log.Println("Shutting down")
}
