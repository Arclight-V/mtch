package main

import (
	"context"
	"fmt"
	"github.com/oklog/run"
	"log"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Arclight-V/mtch/pkg/prober"
	"github.com/Arclight-V/mtch/pkg/signaler"

	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/email"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/jwt_signer"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	config "github.com/Arclight-V/mtch/pkg/platform/config"
	grpcserver "github.com/Arclight-V/mtch/pkg/server/grpc"
	pb "proto"
)

const (
	grpcAddr = "localhost:50051"
)

// move to Vault
var secretAccessKey = []byte("secret-access-key")
var secretRefreshKey = []byte("secret-refresh-key")
var secretVerifyKey = []byte("secret-verify-key")

var handler *httpadapter.Handler

func main() {
	cfg, err := config.GetConfig(os.Getenv("auth-config"))
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	var g run.Group

	conn, err := grpc.NewClient(
		cfg.Client.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(grpcserver.NewUnaryClientRequestIDInterceptor()),
	)
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	grpcProbe := prober.NewGRPC()
	httpProbe := prober.NewHTTP()
	statusProber := prober.Combine(grpcProbe, httpProbe)

	// Listen for reload signals
	{
		shutdown := make(chan struct{})
		g.Add(func() error {
			return WaitForInterrupt(shutdown)
		}, func(err error) {
			close(shutdown)
		})
	}

	repo := grpcclient.NewGRPCUserRepo(pb.NewUserInfoClient(conn))
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
	handler = httpadapter.NewHandler(&userClient, &userClient, &userClient)
	router := httpadapter.NewRouter(handler, httpProbe)

	_ = mime.AddExtensionType(".wasm", "application/wasm")

	// TODO: change this
	srv := &http.Server{
		Addr:    cfg.Http.HTTPAddr,
		Handler: router,
	}

	g.Add(func() error {
		statusProber.Healthy()
		statusProber.Ready()

		log.Printf("server listening at %v", cfg.Http.HTTPAddr)
		return srv.ListenAndServe()

	}, func(err error) {
		statusProber.NotReady(err)
		defer statusProber.NotHealthy(err)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
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
