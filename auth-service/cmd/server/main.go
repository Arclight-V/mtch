package main

import (
	"config"
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/email"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
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
	conn, err := grpc.NewClient(cfg.Client.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	repo := grpcclient.NewGRPCUserRepo(pb.NewUserInfoClient(conn))
	signer := infrastructure.NewJWTSigner(secretAccessKey, secretRefreshKey, secretVerifyKey)
	hasher := crypto.NewBcryptHasher(bcrypt.DefaultCost)
	passwordValidator := passwd.NewUserPasswordValidator()
	emailSender := email.NewNoopSender()
	userClient := auth.Interactor{UserRepo: repo, TokenSigner: signer, Hasher: hasher, PasswordValidator: passwordValidator, EmailSender: emailSender}
	handler = httpadapter.NewHandler(&userClient, &userClient)

	log.Printf("server listening at %v", cfg.Http.HTTPAddr)
	http.ListenAndServe(cfg.Http.HTTPAddr, httpadapter.NewRouter(handler))
}
