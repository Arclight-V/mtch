package main

import (
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"log"
	"mime"
	"os"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Arclight-V/mtch/pkg/logging"
	"github.com/Arclight-V/mtch/pkg/prober"
	"github.com/Arclight-V/mtch/pkg/signaler"
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
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "exiting")

}
