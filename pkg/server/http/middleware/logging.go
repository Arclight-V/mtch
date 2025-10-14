package middleware

import (
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func Logging(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			rid, _ := RequestIDFromContext(r.Context())
			level.Info(logger).Log("http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"request_id", rid,
				"duration", time.Since(start),
				"remote", r.RemoteAddr,
			)
		})
	}
}
