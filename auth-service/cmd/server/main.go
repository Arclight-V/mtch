package main

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/features"
	"github.com/Arclight-V/mtch/pkg/messagebroker"
	"log"
	"mime"
	"os"
	"regexp"

	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/messagebroker/kafka/producer"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	config "github.com/Arclight-V/mtch/pkg/platform/config"
	"github.com/Arclight-V/mtch/pkg/prober"
	"github.com/Arclight-V/mtch/pkg/signaler"
	"github.com/Arclight-V/mtch/pkg/tracing/otel"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/jwt_signer"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/repository"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	httpserver "github.com/Arclight-V/mtch/pkg/server/http"
)

const (
	grpcAddr = "localhost:50051"
)

// move to Vault
var secretAccessKey = []byte("secret-access-key")
var secretRefreshKey = []byte("secret-refresh-key")
var secretVerifyKey = []byte("secret-verify-key")

func main() {
	cfg, err := config.GetConfig(os.Getenv("AUTH_CONFIG"))
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

	metrics := prometheus.NewRegistry()
	metrics.MustRegister(
		versioncollector.NewCollector("mtch-auth-service"),
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
			otel.WithServiceName("auth-service"),
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

	level.Debug(logger).Log("msg", "setting up receive HTTP handler")
	{
		level.Debug(logger).Log("msg", "creating gRPC user-service client", "addr", cfg.UserServiceClient.GRPCAddr)
		connUserService, err := grpc.NewClient(
			cfg.UserServiceClient.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(grpcserver.NewUnaryClientRequestIDInterceptor()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			level.Error(logger).Log("msg", errors.Wrapf(err, "failed to create gRPC client: %v", err))
		}
		defer connUserService.Close()

		level.Debug(logger).Log("msg", "creating gRPC notification-service client", "addr", cfg.NotificationServiceClient.GRPCAddr)
		connNotificationService, err := grpc.NewClient(
			cfg.NotificationServiceClient.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(grpcserver.NewUnaryClientRequestIDInterceptor()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			level.Error(logger).Log("msg", errors.Wrapf(err, "failed to create gRPC client: %v", err))
		}
		defer connNotificationService.Close()

		repo := grpcclient.NewGRPCUserRepo(
			userservicepb.NewUserServiceClient(connUserService),
			notificationservicepb.NewNotificationServiceClient(connNotificationService),
		)
		//authMetrics := httpmetrics.NewAuthMetrics(metrics)
		signer := jwt_signer.NewJWTSigner(secretAccessKey, secretRefreshKey, secretVerifyKey)
		hasher := crypto.NewBcryptHasher(bcrypt.DefaultCost)
		passwordValidator := passwd.NewUserPasswordValidator()
		verifyTokenRepo := repository.NewVerifyTokensMem()

		var publisher messagebroker.Publisher

		kafkaEnable := featureList.IsEnabled(feature_list.FeatureKafka)
		if kafkaEnable {
			p, err := producer.New(cfg.Kafka.Producer, logger,
				producer.WithCompressionType(cfg.Kafka.Producer.CompressionType),
				producer.WithAcks(cfg.Kafka.Producer.Acks),
				producer.WithLingerMS(cfg.Kafka.Producer.LingerMS),
				producer.WithFlushTimeoutMS(cfg.Kafka.Producer.FlushTimeoutMS),
				producer.WithEnableIdempotence(cfg.Kafka.Producer.EnableIdempotence),
			)
			if err != nil {
				level.Error(logger).Log("msg", "failed to create kafka publisher", "err", err)
				os.Exit(1)
			}
			publisher = p
		}

		userClient := auth.Interactor{
			UserRepo:          repo,
			TokenSigner:       signer,
			Hasher:            hasher,
			PasswordValidator: passwordValidator,
			VerifyTokenRepo:   verifyTokenRepo,
			Publisher:         publisher,
			FeatureList:       featureList,
		}

		webHandler := httpadapter.NewHandler(logger,
			&httpadapter.Options{
				ListenAddress: cfg.Http.ListenAddr,
				Registry:      metrics,
				FrontendPath:  cfg.FrontEnd.FrontendPath,
			},

			&userClient,
			&userClient,
			&userClient)

		_ = mime.AddExtensionType(".wasm", "application/wasm")

		g.Add(func() error {
			return errors.Wrap(webHandler.Run(), "error starting web server")
		}, func(err error) {
			webHandler.Shutdown()
			if kafkaEnable = featureList.IsEnabled(feature_list.FeatureKafka); kafkaEnable {
				_ = publisher.Close()
			}
			openfeature.Shutdown()
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "exiting")

}
