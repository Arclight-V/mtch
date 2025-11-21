package main

import (
	"context"
	"log"
	"os"
	"regexp"

	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"go.opentelemetry.io/otel/attribute"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/notificationservice"
	"github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/Arclight-V/mtch/pkg/prober"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
	"github.com/Arclight-V/mtch/pkg/signaler"
	"github.com/Arclight-V/mtch/pkg/tracing/otel"

	grpcnotification "github.com/Arclight-V/mtch/notification/internal/adapter/grpc/notification"
	"github.com/Arclight-V/mtch/notification/internal/features"
	"github.com/Arclight-V/mtch/notification/internal/infrastructure/email"
	usecase "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
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

	featureList, err := feature_list.NewFeatureList(provider, "mtch-notification", logger, features.Features)
	if err != nil {
		// If a FeatureList initialization error occurs, log it and exit
		level.Error(logger).Log("msg", "failed to create FeatureList", "err", err)
		os.Exit(1)
	}

	metrics := prometheus.NewRegistry()
	metrics.MustRegister(
		versioncollector.NewCollector("mtch-notification"),
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

		traceExporter, err := otel.NewTraceExporterGRPC(baseCtx /*TODO: add an endpoint*/)
		if err != nil {
			log.Fatalf("failed to create trace exporter: %v", err)
		}

		otelShutdown, err := otel.SetupOTelSDK(
			baseCtx,
			// TODO:: config - taking values from config.yml
			otel.WithServiceName("notification"),
			otel.WithAttributes(
				attribute.String("env", "dev"),
				attribute.String("version", "1.0.0"),
			),
			otel.WithExporter(traceExporter),
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

	emailSender := email.NewSMTPClient(cfg)
	notificationUC := usecase.NewNotificationUseCase(emailSender, logger, featureList)
	server := grpcnotification.NewNotificationServiceServer(notificationUC, logger, featureList)

	level.Debug(logger).Log("msg", "starting GRPC server")
	{
		//TODO Handle err
		grpcLogOpts, _ := logging.NewGRPCOption()

		s := grpcserver.NewServer(logger, metrics, grpcLogOpts, grpcProbe,
			grpcserver.WithServer(notificationservice.RegisterNotificationServer(server)),
			grpcserver.WithListen(cfg.NotificationServiceServer.Port),
			grpcserver.WithGracePeriod(cfg.NotificationServiceServer.GracePeriod),
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
