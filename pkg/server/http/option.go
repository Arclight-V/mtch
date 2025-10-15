package http

import "net/http"

type options struct {
	listen string
	mux    *http.ServeMux

	handler http.Handler
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithListen(listen string) Option {
	return optionFunc(func(o *options) {
		o.listen = listen
	})
}

func WithMux(mux *http.ServeMux) Option {
	return optionFunc(func(o *options) {
		o.mux = mux
	})
}
