package prober

import (
	"google.golang.org/grpc/health"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCProbe struct {
	h *health.Server
}

func NewGRPC() *GRPCProbe {
	h := health.NewServer()
	h.SetServingStatus("", grpc_health.HealthCheckResponse_NOT_SERVING)

	return &GRPCProbe{h: h}
}

func (p *GRPCProbe) HealthServer() *health.Server { return p.h }

func (p *GRPCProbe) Healthy() {
	p.h.Resume()
}

func (p *GRPCProbe) NotHealthy(err error) {
	p.h.Shutdown()
}

func (p *GRPCProbe) Ready() {
	p.h.SetServingStatus("", grpc_health.HealthCheckResponse_SERVING)
}

func (p *GRPCProbe) NotReady(err error) {
	p.h.SetServingStatus("", grpc_health.HealthCheckResponse_NOT_SERVING)
}
