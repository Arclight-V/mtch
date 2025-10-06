package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Arclight-V/mtch/auth-service/internal/adapter/http/models"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	regUC         auth.RegisterUseCase
	loginUC       auth.LoginUseCase
	verifyEmailUC auth.VerifyEmailUseCase
	validate      *validator.Validate
}

func NewHandler(
	regUC auth.RegisterUseCase,
	loginUC auth.LoginUseCase,
	verifyEmailUC auth.VerifyEmailUseCase) *Handler {

	validate := validator.New()
	validate.RegisterAlias("contact", "email|e164")

	return &Handler{
		regUC:         regUC,
		loginUC:       loginUC,
		verifyEmailUC: verifyEmailUC,
		validate:      validate,
	}
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
// @Param        RegisterRequest  body  models.RegisterRequest  true "registration payload"
// @Success      201 {object}     models.RegisterResponse
// @Failure      400 {object} 	  models.ErrorResponse
// @Failure      409 {object} 	  models.ErrorResponse
// @Failure      415 {object} 	  models.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register called")

	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	var in models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(in)
	if err := h.validate.Struct(in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	birthDay, err := strconv.Atoi(in.BirthDay)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	birthMonth, err := strconv.Atoi(in.BirthMonth)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	birthEarth, err := strconv.Atoi(in.BirthYear)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	regInput := &auth.RegisterInput{
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Contact:   in.Contact,
		Password:  in.Password,
		Date: &auth.Date{
			BirthDay:   int32(birthDay),
			BirthMonth: int32(birthMonth),
			BirthYear:  int32(birthEarth),
		},
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	regOutput, err := h.regUC.Register(ctx, regInput)
	if err != nil {
		writeJSONError(w, http.StatusConflict, err.Error())
		return
	}

	out := models.RegisterResponse{
		User: models.PendingUserDTO{
			UserID:   regOutput.UserID,
			Email:    regOutput.Email,
			Verified: false,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&out)
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	log.Println("Verify email called")

	//token := pat.Param(r, "token")
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSONError(w, http.StatusBadRequest, "token is required")
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	verifyOut, err := h.verifyEmailUC.VerifyEmail(ctx, auth.VerifyEmailInput{Token: token})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	out := models.VerifyEmailResponse{
		User: models.VerifiedEmailUserDTO{
			UserID:     verifyOut.UserID,
			VerifiedAt: verifyOut.VerifiedAt,
			Verified:   verifyOut.Verified,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&out)
	fmt.Println(token)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
