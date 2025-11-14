package main

import (
	"context"
	"github.com/Arclight-V/mtch/notification-service/internal/features"
	"github.com/Arclight-V/mtch/pkg/notificationservice"
	"github.com/Arclight-V/mtch/pkg/prober"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
	"github.com/Arclight-V/mtch/pkg/signaler"

	grpcnotification "github.com/Arclight-V/mtch/notification-service/internal/adapter/grpc/notification"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"log"
	"os"
	"regexp"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/go-kit/log/level"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
)

func main() {
	cfg, err := config.GetConfig(os.Getenv("NOTIFICATION_CONFIG"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logging.NewLogger(cfg.LogCfg.Level, cfg.LogCfg.Format, cfg.LogCfg.DebugName)

	// Use flagd as the OpenFeature provider
	provider, err := flagd.NewProvider(
		flagd.WithFileResolver(),
		flagd.WithOfflineFilePath(cfg.FlagD.FlagsPath),
	)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize flagd", "err", err)
		os.Exit(1)
	}

	featureList, err := feature_list.NewFeatureList(provider, "mtch-auth-service", logger, features.Features)
	if err != nil {
		// If a FeatureList initialization error occurs, log it and exit
		level.Error(logger).Log("msg", "failed to create FeatureList", "err", err)
		os.Exit(1)
	}

	_ = featureList

	metrics := prometheus.NewRegistry()
	metrics.MustRegister(
		versioncollector.NewCollector("mtch-notification-service"),
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
		),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	prometheus.DefaultRegisterer = metrics

	var (
		g       run.Group
		baseCtx = context.Background()
	)

	ctx, cancel := context.WithCancel(baseCtx)

	// Listen for reload signals
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
			statusProber.Ready()

			return srv.ListenAndServe()

		}, func(err error) {
			statusProber.NotReady(err)
			defer statusProber.NotHealthy(err)

			srv.Shutdown(err)
		})

	}

	server := grpcnotification.NewNotificationServiceServer()

	level.Debug(logger).Log("msg", "starting GRPC server")
	{
		//TODO Handle err
		grpcLogOpts, _ := logging.NewGRPCOption()

		s := grpcserver.NewServer(logger, metrics, grpcLogOpts, grpcProbe,
			grpcserver.WithServer(notificationservice.RegisterNotificationUserServer(server)),
			grpcserver.WithListen(cfg.NotificationServer.Port),
			grpcserver.WithGracePeriod(cfg.NotificationServer.GracePeriod),
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
