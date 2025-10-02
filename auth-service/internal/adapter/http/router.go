package httpadapter

import (
	limiter "github.com/ulule/limiter/v3"
	goji "goji.io"
	"goji.io/pat"
	"log"
	"net/http"

	_ "github.com/Arclight-V/mtch/auth-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"
)

const (
	apiBase = "/api/v1/"

	// rateFormatted - a temporary solution (read from config)
	rateFormatted = "5-M"
)

func NewRouter(h *Handler) http.Handler {

	root := goji.NewMux()
	root.Use(rateLimiter())

	api := goji.SubMux()
	root.Handle(pat.New("/swagger/*"), httpSwagger.WrapHandler)
	root.Handle(pat.New(apiBase+"*"), api)

	auth := goji.SubMux()
	api.Handle(pat.New("/auth/*"), auth)

	auth.HandleFunc(pat.Post("/register"), h.Register)
	auth.HandleFunc(pat.Get("/verify-email"), h.VerifyEmail)
	auth.HandleFunc(pat.Post("/login"), h.Login)

	return root
}

// rateLimiter() - move to config?
func rateLimiter() func(http.Handler) http.Handler {
	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted(rateFormatted)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	store := mem.NewStore()

	// Create a new middleware with the limiter instance.
	middleware := mhttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
	return middleware.Handler
}
