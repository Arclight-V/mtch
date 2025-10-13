package prober

import (
	"io"
	"net/http"
	"sync/atomic"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type check func() bool

type HTTPProbe struct {
	ready   atomic.Uint32
	healthy atomic.Uint32
}

func NewHTTP() *HTTPProbe { return &HTTPProbe{} }

func (p *HTTPProbe) HealthyHandler(logger log.Logger) http.Handler {
	return p.handler(logger, p.isHeathy)
}

func (p *HTTPProbe) ReadyHandler(logger log.Logger) http.Handler {
	return p.handler(logger, p.IsReady)
}

func (p *HTTPProbe) handler(logger log.Logger, c check) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if !c() {
			http.Error(w, "NOT OK", http.StatusServiceUnavailable)
			return
		}
		if _, err := io.WriteString(w, "OK"); err != nil {
			level.Error(logger).Log("msg", "failed to write response", "err", err)
		}
	}
}

func (p *HTTPProbe) IsReady() bool {
	ready := p.ready.Load()
	return ready > 0
}

func (p *HTTPProbe) isHeathy() bool {
	healthy := p.healthy.Load()
	return healthy > 0
}

func (p *HTTPProbe) Ready() {
	p.ready.Swap(1)
}

func (p *HTTPProbe) NotReady(err error) {
	p.ready.Swap(0)
}

func (p *HTTPProbe) Healthy() {
	p.healthy.Swap(1)
}

func (p *HTTPProbe) NotHealthy(err error) {
	p.healthy.Swap(0)
}
