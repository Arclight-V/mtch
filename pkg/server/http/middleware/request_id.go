package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type ctxKey int

const reqIDKey = ctxKey(0)

func NewContextWithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, reqIDKey, rid)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(reqIDKey).(string)
	return rid, ok
}

func RequestID(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
			r.Header.Set("X-Request-ID", reqID)
		}
		ctx := NewContextWithRequestID(r.Context(), reqID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
