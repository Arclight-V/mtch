package tech

import (
	"context"
	apphealth "github.com/Arclight-V/mtch/pkg/health"
	"net/http"
	"time"
)

type Handler struct {
	live  apphealth.LivenessChecker
	ready apphealth.ReadinessChecker

	perCheckTimout time.Duration
}

func NewHandler(l apphealth.LivenessChecker, r apphealth.ReadinessChecker) *Handler {
	return &Handler{live: l, ready: r, perCheckTimout: 200 * time.Millisecond}
}

func (h *Handler) Livez(w http.ResponseWriter, r *http.Request) {
	if err := h.live.Alive(r.Context()); err != nil {
		http.Error(w, "not alive", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.perCheckTimout)
	defer cancel()

	if err := h.ready.Ready(ctx); err != nil {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
