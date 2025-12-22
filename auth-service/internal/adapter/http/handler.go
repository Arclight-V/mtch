package httpadapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	goji "goji.io"
	"goji.io/pat"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	"github.com/Arclight-V/mtch/auth-service/internal/adapter/http/models"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/pkg/server/http/middleware"
)

const (
	apiBase = "/api/v1/"

	// rateFormatted - a temporary solution (read from config)
	rateFormatted = "5-M"
)

// Options for the web Handler
type Options struct {
	ListenAddress string
	TLSConfig     *tls.Config
	Registry      *prometheus.Registry

	FrontendPath string
}
type Handler struct {
	router  http.Handler
	options *Options
	httpSrv *http.Server
	regUC   auth.RegisterUseCase

	loginUC  auth.LoginUseCase
	verifyUC auth.VerifyUseCase
	validate *validator.Validate

	logger      log.Logger
	featureList *feature_list.FeatureList

	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandler(
	logger log.Logger,
	featureList *feature_list.FeatureList,
	o *Options,

	regUC auth.RegisterUseCase,
	loginUC auth.LoginUseCase,
	verifyUC auth.VerifyUseCase) *Handler {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	var registerer prometheus.Registerer = nil
	if o.Registry != nil {
		registerer = o.Registry
	}

	validate := validator.New()
	validate.RegisterAlias("contact", "email|e164")

	h := &Handler{
		logger:      logger,
		featureList: featureList,
		options:     o,

		regUC:    regUC,
		loginUC:  loginUC,
		verifyUC: verifyUC,
		validate: validate,

		requestsTotal: promauto.With(registerer).NewCounterVec(
			prometheus.CounterOpts{
				Name: "auth_http_requests_total",
				Help: "Total number of HTTP requests",
			}, []string{"code", "method"},
		),
		requestDuration: promauto.With(registerer).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			}, []string{"handler"},
		),
	}

	router := goji.NewMux()
	//router.Use(rateLimiter())
	//TODO:: name from config
	router.Use(otelhttp.NewMiddleware("auth-service"))
	router.Use(requestID)
	router.Use(logging(logger))

	router.Handle(pat.New("/swagger/*"), httpSwagger.WrapHandler)

	api := goji.SubMux()
	router.Handle(pat.New(apiBase+"*"), api)

	authMux := goji.SubMux()
	api.Handle(pat.New("/auth/*"), authMux)

	instrf := func(route string, next http.Handler) http.HandlerFunc {

		return promhttp.InstrumentHandlerDuration(
			h.requestDuration.MustCurryWith(prometheus.Labels{"handler": route}),
			promhttp.InstrumentHandlerCounter(h.requestsTotal, next),
		)
	}

	authMux.HandleFunc(
		pat.Post("/register"),
		instrf("POST "+apiBase+"/auth/register",
			http.HandlerFunc(h.Register),
		),
	)

	authMux.HandleFunc(
		pat.Post("/verify-code"),
		instrf("POST "+apiBase+"/auth/verify-code",
			http.HandlerFunc(h.VerifyCode),
		),
	)

	authMux.HandleFunc(pat.Post("/login"),
		instrf("POST "+apiBase+"/auth/verify-email",
			http.HandlerFunc(h.Login),
		),
	)

	static := http.StripPrefix("/app/", http.FileServer(http.Dir(resolveWebAppDir(logger, o.FrontendPath))))
	// Redirect /app -> /app/
	router.Handle(pat.Get("/app"), http.RedirectHandler("/app/", http.StatusMovedPermanently))
	// All files front: /app/*
	router.Handle(pat.Get("/app/*"), static)

	h.router = router

	h.httpSrv = &http.Server{
		Handler:   router,
		TLSConfig: h.options.TLSConfig,
	}

	return h
}

func (h *Handler) Run() error {
	level.Info(h.logger).Log("msg", "Start listening for connections", "address", h.options.ListenAddress)

	listener, err := net.Listen("tcp", h.options.ListenAddress)
	if err != nil {
		return err
	}

	//TODO: Add
	// Monitor incoming connections with conntrack.
	//listener = conntrack.NewListener(listener,
	//	conntrack.TrackWithName("http"),
	//	conntrack.TrackWithTracing())

	if h.options.TLSConfig != nil {
		level.Info(h.logger).Log("msg", "Serving HTTPS", "address", h.options.ListenAddress)
		return h.httpSrv.ServeTLS(listener, "", "")
	}

	level.Info(h.logger).Log("msg", "Serving plain HTTP", "address", h.options.ListenAddress)
	return h.httpSrv.Serve(listener)
}

func (h *Handler) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = h.httpSrv.Shutdown(ctx)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
	//	User        *pb.User `json:"userservice"`
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
// @Description Creates an account, returns an unverified userservice
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
	ctx, span := otel.Tracer("auth-service").Start(r.Context(), "Register")
	defer span.End()

	level.Info(h.logger).Log("msg", "Register called")

	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	var in models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	birthDay, err := strconv.Atoi(in.BirthDay)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	birthMonth, err := strconv.Atoi(in.BirthMonth)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	birthEarth, err := strconv.Atoi(in.BirthYear)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
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

func (h *Handler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	level.Debug(h.logger).Log("msg", "VerifyCode called")
	if !h.featureList.IsEnabled(feature_list.VerifyCodeEnabled) {
		writeJSONError(w, http.StatusNotImplemented, "Verify code is disabled")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}
	var in models.VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validate.Struct(in); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	verifyOut, err := h.verifyUC.VerifyCode(ctx, &auth.VerifyInput{Code: in.Code})
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	out := models.VerifyCodeResponse{
		UserID:     verifyOut.UserID,
		VerifiedAt: verifyOut.VerifiedAt,
		Verified:   verifyOut.Verified,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(&out)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func requestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func logging(logger log.Logger) func(http.Handler) http.Handler {
	return middleware.Logging(logger)
}

// rateLimiter() - move to config?
func rateLimiter() func(http.Handler) http.Handler {
	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted(rateFormatted)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	store := mem.NewStore()

	// Create a new middleware with the limiter instance.
	middleware := mhttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
	return middleware.Handler
}

func resolveWebAppDir(logger log.Logger, path string) string {
	if override := path; override != "" {
		if stat, err := os.Stat(override); err == nil && stat.IsDir() {
			return override
		} else {
			level.Warn(logger).Log("msg", "WEBWASM_DIR does not exist or is not a directory", "path", override, "err", err)
		}
	}

	candidates := []string{
		"./webwasm",
		filepath.Clean("../webwasm"),
		filepath.Clean("./../webwasm"),
	}

	for _, candidate := range candidates {
		if stat, err := os.Stat(candidate); err == nil && stat.IsDir() {
			return candidate
		}
	}

	level.Warn(logger).Log("msg", "webwasm assets directory not found, serving empty handler")
	return "./webwasm"
}
