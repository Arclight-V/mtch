package httpmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type HTTPMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewAuthMetrics(reg prometheus.Registerer) *HTTPMetrics {
	m := new(HTTPMetrics)

	factory := promauto.With(reg)
	m.requestsTotal = factory.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "status"})
	m.requestDuration = factory.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"route"})

	return m
}

func (m *HTTPMetrics) Instrument(route string, h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(
		m.requestDuration.MustCurryWith(prometheus.Labels{"route": route}),
		promhttp.InstrumentHandlerCounter(m.requestsTotal, h),
	)
}

func (m *HTTPMetrics) InstrumentFunc(route string, hf http.HandlerFunc) http.Handler {
	return m.Instrument(route, http.HandlerFunc(hf))
}
