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
			level.Info(logger).Log(
				"path", r.URL.Path,
				"method", r.Method,
				"remote", r.RemoteAddr,
				"duration", time.Since(start),
				"request_id", rid,
			)
		})
	}
}
