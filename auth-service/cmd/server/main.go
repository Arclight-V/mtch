package main

import (
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/grpcclient"
	httpadapter "github.com/Arclight-V/mtch/auth-service/internal/adapter/http"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	pb "proto"
)

const (
	grpcAddr = "localhost:50051"
)

var handler *httpadapter.Handler

func main() {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	repo := grpcclient.NewGRPCUserRepo(pb.NewUserInfoClient(conn))
	userClient := auth.Interactor{UserRepo: repo}
	handler = httpadapter.NewHandler(&userClient)

	log.Printf("server listening at %v", 8000)
	http.ListenAndServe(":8000", httpadapter.NewRouter(handler))
}
