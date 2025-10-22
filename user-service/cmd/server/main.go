package main

import (
	"context"
	"log"
	"os"
	"regexp"

	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"go.opentelemetry.io/otel/attribute"

	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/Arclight-V/mtch/pkg/prober"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
	"github.com/Arclight-V/mtch/pkg/signaler"
	"github.com/Arclight-V/mtch/pkg/tracing/otel"
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

	// Setup optional tracing.
	{
		var (
			baseCtx = context.Background()
		)

		otelShutdown, err := otel.SetupOTelSDK(
			baseCtx,
			// TODO:: config - taking values from config.yml
			otel.WithServiceName("auth-service"),
			otel.WithAttributes(
				attribute.String("env", "dev"),
				attribute.String("version", "1.0.0"),
			),
		)
		if err != nil {
			log.Fatalf("failed to setup OTel SDK: %v", err)
		}

		ctx, cancel := context.WithCancel(baseCtx)
		g.Add(func() error {
			<-ctx.Done()
			return ctx.Err()
		}, func(error) {
			if otelShutdown != nil {
				if err := otelShutdown(ctx); err != nil {
					level.Warn(logger).Log("msg", "OTel SDK shutdown failed", "err", err)
				}
			}
			cancel()
		})

	}

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
		//TODO Handle err
		grpcLogOpts, _ := logging.NewGRPCOption()

		s := grpcserver.NewServer(logger, metrics, grpcLogOpts, grpcProbe,
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
