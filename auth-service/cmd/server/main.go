package main

import (
	"context"
	"encoding/json"
	"fmt"
	goji "goji.io"
	"goji.io/pat"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	pb "proto"
	"time"
)

const (
	grpcAddr = "localhost:50051"
)

var userClient pb.UserInfoClient

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login called")
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	resp, err := userClient.Login(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	fmt.Fprintf(w, "Hello, %s %s!", resp.User.FirstName, resp.User.LastName)
	w.WriteHeader(http.StatusOK)

}

func main() {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not create grpc connection: %v", err)
	}
	defer conn.Close()

	userClient = pb.NewUserInfoClient(conn)

	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/login"), login)

	log.Printf("server listening at %v", 8000)
	http.ListenAndServe(":8000", mux)
}
