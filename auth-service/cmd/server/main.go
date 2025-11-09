package main

import (
	"context"
	"github.com/Arclight-V/mtch/pkg/messagebroker/kafka/producer"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"mime"
	"os"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/prober"
	"github.com/Arclight-V/mtch/pkg/signaler"
	"github.com/Arclight-V/mtch/pkg/tracing/otel"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/email"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/jwt_signer"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/repository"
	config "github.com/Arclight-V/mtch/pkg/platform/config"
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
		conn, err := grpc.NewClient(
			cfg.Client.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(grpcserver.NewUnaryClientRequestIDInterceptor()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			level.Error(logger).Log("msg", errors.Wrapf(err, "failed to create gRPC client: %v", err))
		}
		defer conn.Close()

		repo := grpcclient.NewGRPCUserRepo(userservicepb.NewUserServiceClient(conn))
		//authMetrics := httpmetrics.NewAuthMetrics(metrics)
		signer := jwt_signer.NewJWTSigner(secretAccessKey, secretRefreshKey, secretVerifyKey)
		hasher := crypto.NewBcryptHasher(bcrypt.DefaultCost)
		passwordValidator := passwd.NewUserPasswordValidator()
		emailSender := email.NewSMTPClient(cfg)
		verifyTokenRepo := repository.NewVerifyTokensMem()
		publisher, err := producer.New(cfg.Kafka.Producer, logger,
			producer.WithCompressionType(cfg.Kafka.Producer.CompressionType),
			producer.WithAcks(cfg.Kafka.Producer.Acks),
			producer.WithLingerMS(cfg.Kafka.Producer.LingerMS),
			producer.WithFlushTimeoutMS(cfg.Kafka.Producer.FlushTimeoutMS),
			producer.WithEnableIdempotence(cfg.Kafka.Producer.EnableIdempotence),
		)
		_ = publisher
		if err != nil {
			level.Error(logger).Log("msg", "failed to create kafka publisher", "err", err)
			os.Exit(1)
		}
		userClient := auth.Interactor{
			UserRepo:          repo,
			TokenSigner:       signer,
			Hasher:            hasher,
			PasswordValidator: passwordValidator,
			EmailSender:       emailSender,
			VerifyTokenRepo:   verifyTokenRepo,
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
			_ = publisher.Close()
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "exiting")

}
