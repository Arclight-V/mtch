package main

import (
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure"
	"github.com/Arclight-V/mtch/auth-service/internal/infrastructure/crypto"
	passwd "github.com/Arclight-V/mtch/auth-service/internal/infrastructure/password_validator"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
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
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	repo := grpcclient.NewGRPCUserRepo(pb.NewUserInfoClient(conn))
	log.Println(repo)
	signer := infrastructure.NewJWTSigner(secretAccessKey, secretRefreshKey, secretVerifyKey)
	userClient := auth.Interactor{UserRepo: repo, TokenSigner: signer}
	hasher := crypto.NewBcryptHasher(bcrypt.DefaultCost)
	passwordValidator := passwd.NewUserPasswordValidator()
	handler = httpadapter.NewHandler(&userClient, &userClient, hasher, passwordValidator)

	log.Printf("server listening at %v", 8000)
	http.ListenAndServe(":8000", httpadapter.NewRouter(handler))
}
