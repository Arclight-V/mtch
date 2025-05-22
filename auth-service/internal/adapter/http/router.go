package httpadapter

import (
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/login"), h.Login)
	mux.HandleFunc(pat.Post("/register"), h.Register)
	return mux
}
