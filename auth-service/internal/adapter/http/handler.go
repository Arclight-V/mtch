package httpadapter

import (
	"context"
	"encoding/json"
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/http/dto"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	pb "proto"
	"time"
)

type Handler struct {
	uc       *auth.Interactor
	validate *validator.Validate
}

func NewHandler(uc *auth.Interactor) *Handler {
	return &Handler{uc: uc, validate: validator.New()}
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
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	out := struct {
		User        *pb.User `json:"user"`
		AccessToken string   `json:"access_token"`
		ExpiresIn   int64    `json:"expires_in"`
	}{
		User:        resp.User,
		AccessToken: resp.AccessToken,
		ExpiresIn:   resp.ExpiresIn,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&out)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register called")
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var in dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	protoReq := &pb.RegisterRequest{
		Email:    in.Email,
		Password: in.Password,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	protoResp, err := h.uc.Register(ctx, protoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// TODO:: add logic for calling the email service

	out := dto.RegisterResponse{
		User: dto.PendingUserDTO{
			UserID:   protoResp.User.Uuid,
			Email:    protoResp.User.Email,
			CreateAt: protoResp.User.CreatedAt.AsTime(),
			Verified: protoResp.User.Verified,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&out)
}
