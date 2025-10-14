package httpadapter

import (
	"fmt"
	goji "goji.io"
	"goji.io/pat"
	"mime"
	"net/http"

	"github.com/go-kit/log"
	httpSwagger "github.com/swaggo/http-swagger"
	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/Arclight-V/mtch/pkg/server/http/middleware"

	_ "github.com/Arclight-V/mtch/auth-service/docs"
)

const (
	apiBase = "/api/v1/"

	// rateFormatted - a temporary solution (read from config)
	rateFormatted = "5-M"
)

// TODO: Transform to server and add logger
func NewRouter(h *Handler) http.Handler {
	_ = mime.AddExtensionType(".wasm", "application/wasm")

	root := goji.NewMux()
	root.Use(rateLimiter())
	root.Use(requestID)
	root.Use(logging(h.logger))

	api := goji.SubMux()
	root.Handle(pat.New("/swagger/*"), httpSwagger.WrapHandler)
	root.Handle(pat.New(apiBase+"*"), api)

	auth := goji.SubMux()
	api.Handle(pat.New("/auth/*"), auth)
	auth.HandleFunc(pat.Post("/register"), h.Register)
	auth.HandleFunc(pat.Get("/verify-email"), h.VerifyEmail)
	auth.HandleFunc(pat.Post("/login"), h.Login)

	static := http.StripPrefix("/app/", http.FileServer(http.Dir("./../webwasm")))
	// Редирект /app -> /app/
	root.Handle(pat.Get("/app"), http.RedirectHandler("/app/", http.StatusMovedPermanently))
	// Все файлы фронта: /app/*
	root.Handle(pat.Get("/app/*"), static)

	return root
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

func requestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

func logging(logger log.Logger) func(http.Handler) http.Handler {
	return middleware.Logging(logger)
}
