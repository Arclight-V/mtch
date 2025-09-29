package httpadapter

import (
	"context"
	"encoding/json"
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/http/dto"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	regUC    auth.RegisterUseCase
	loginUC  auth.LoginUseCase
	validate *validator.Validate
}

func NewHandler(regUC auth.RegisterUseCase, loginUC auth.LoginUseCase) *Handler {
	return &Handler{regUC: regUC, loginUC: loginUC, validate: validator.New()}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login called")
	//if r.Header.Get("Content-Type") != "application/json" {
	//	http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
	//	return
	//}
	//var req pb.LoginRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	//defer cancel()
	//resp, err := h.loginUC.Login(ctx, &req)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusUnauthorized)
	//}
	//http.SetCookie(w, &http.Cookie{
	//	Name:     "refresh_token",
	//	Value:    resp.RefreshToken,
	//	Path:     "/",
	//	HttpOnly: true,
	//	Secure:   true,
	//	SameSite: http.SameSiteLaxMode,
	//	Expires:  time.Now().Add(30 * 24 * time.Hour),
	//})
	//
	//out := struct {
	//	User        *pb.User `json:"user"`
	//	AccessToken string   `json:"access_token"`
	//	ExpiresIn   int64    `json:"expires_in"`
	//}{
	//	User:        resp.User,
	//	AccessToken: resp.AccessToken,
	//	ExpiresIn:   resp.ExpiresIn,
	//}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(&out)
}

// Register
// @Summary User Registration
// @Description Creates an account, returns an unverified user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterRequest  body  dto.RegisterRequest  true "registration payload"
// @Success      201 {object}     dto.RegisterResponse
// @Failure      400 {object} 	  dto.ErrorResponse
// @Failure      409 {object} 	  dto.ErrorResponse
// @Failure      415 {object} 	  dto.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register called")
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	var in dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	regInput := auth.RegisterInput{
		Email:    in.Email,
		Password: in.Password,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	regOutput, err := h.regUC.Register(ctx, regInput)
	if err != nil {
		writeJSONError(w, http.StatusConflict, err.Error())
		return
	}

	out := dto.RegisterResponse{
		User: dto.PendingUserDTO{
			UserID:      regOutput.UserID,
			Email:       regOutput.Email,
			VerifyToken: regOutput.VerifyToken,
			Verified:    false,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&out)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}
