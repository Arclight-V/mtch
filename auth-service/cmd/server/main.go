package main

import (
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"config"
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/email"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/jwt_signer"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	pb "proto"

	signaler "github.com/Arclight-V/mtch/pkg/signaler"
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
	conn, err := grpc.NewClient(cfg.Client.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	shutdown := make(chan struct{})
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

	_ = mime.AddExtensionType(".wasm", "application/wasm")

	go func() {
		log.Printf("server listening at %v", cfg.Http.HTTPAddr)
		if err := http.ListenAndServe(cfg.Http.HTTPAddr, httpadapter.NewRouter(handler)); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	go waitForInterrupt(shutdown)
	<-shutdown

}

func waitForInterrupt(waiter chan<- struct{}) {
	interrupt := signaler.WaitForInterrupt()
	log.Printf("Captured %v, shutdown requested.\n", interrupt)
	waiter <- struct{}{}
}
