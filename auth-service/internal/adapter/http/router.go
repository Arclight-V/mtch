package httpadapter

import (
	goji "goji.io"
	"goji.io/pat"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	root := goji.NewMux()

	api := goji.SubMux()
	root.Handle(pat.New("/api/v1/*"), api)

	auth := goji.SubMux()
	api.Handle(pat.New("/auth/*"), auth)

	auth.HandleFunc(pat.Post("/register"), h.Register)
	auth.HandleFunc(pat.Post("/login"), h.Login)

	return root
}
