package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/Arclight-V/mtch/pkg/prober"
)

type Server struct {
	logger log.Logger

	mux *http.ServeMux
	srv *http.Server

	opts options
}

func NewServer(logger log.Logger, probe *prober.HTTPProbe, opts ...Option) *Server {
	options := options{}

	for _, o := range opts {
		o.apply(&options)
	}

	mux := http.NewServeMux()
	if options.mux != nil {
		mux = options.mux
	}

	registerProber(mux, probe)

	if options.handler != nil {
		mux.Handle("/", options.handler)
	}

	var h http.Handler = mux

	return &Server{
		logger: logger,
		mux:    mux,
		srv:    &http.Server{Addr: options.listen, Handler: h},
		opts:   options,
	}

}

func (s *Server) ListenAndServe() error {
	level.Info(s.logger).Log("msg", "listen for request and metrics", "addr", s.opts.listen)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(err error) {
	level.Info(s.logger).Log("msg", "internal shutting down server", "err", err.Error())
	if err == http.ErrServerClosed {
		level.Warn(s.logger).Log("msg", "internal server closed unexpectedly")
		return
	}

	//TODO:: add grace period for shot-down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		level.Error(s.logger).Log("msg", "failed to shutdown http server", "err", err.Error())
		return
	}
	level.Info(s.logger).Log("msg", "http server shutdown gracefully")
}

func registerProber(mux *http.ServeMux, p *prober.HTTPProbe) {
	mux.Handle("/-/healthy", p.HealthyHandler(nil))
	mux.Handle("/-/readyz", p.ReadyHandler(nil))
}
