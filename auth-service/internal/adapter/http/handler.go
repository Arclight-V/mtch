package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"log"
	"net/http"
	pb "proto"
	"time"
)

type Handler struct {
	uc *auth.Interactor
}

func NewHandler(uc *auth.Interactor) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
	resp, err := h.uc.Login(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	if resp != nil {
		fmt.Fprintf(w, "Hello, %s %s!", resp.User.FirstName, resp.User.LastName)

	}
	w.WriteHeader(http.StatusOK)
}
