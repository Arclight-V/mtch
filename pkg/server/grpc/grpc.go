package grpc

import (
	"context"
	"fmt"
	"math"
	"net"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/Arclight-V/mtch/pkg/prober"
)

type Server struct {
	logger log.Logger

	srv      *grpc.Server
	listener net.Listener

	opts options
}

func NewServer(logger log.Logger, probe *prober.GRPCProbe, opts ...Option) *Server {
	logger = log.With(logger, "service", "gRPC/server")
	options := options{
		network: "tcp",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	options.grpcOpts = append(options.grpcOpts, []grpc.ServerOption{
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.ChainUnaryInterceptor(
			NewUnaryServerRequestIDInterceptor(),
		),
	}...)

	if options.tlsConfig != nil {
		options.grpcOpts = append(options.grpcOpts, grpc.Creds(credentials.NewTLS(options.tlsConfig)))
	}

	s := grpc.NewServer(options.grpcOpts...)

	grpc_health.RegisterHealthServer(s, probe.HealthServer())

	return &Server{
		logger: logger,
		srv:    s,
		opts:   options,
	}
}

func (s *Server) ListenAndServe() error {
	l, err := net.Listen(s.opts.network, s.opts.listen)
	if err != nil {
		return fmt.Errorf("%s, listen gRPC on address %s", err, s.opts.listen)
	}
	s.listener = l

	level.Info(s.logger).Log("msg", "listening for serving gRPC", "address", s.opts.listen)
	return fmt.Errorf("%s serve gRPC", s.srv.Serve(s.listener))
}

func (s *Server) Shutdown(err error) {
	level.Info(s.logger).Log("msg", "internal server is shutting down", "err", err)

	if s.opts.gracePeriod == 0 {
		s.srv.Stop()
		level.Info(s.logger).Log("msg", "internal server is shutdown", "err", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.opts.gracePeriod)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		level.Info(s.logger).Log("msg", "gracefully stopping internal server")
		s.srv.GracefulStop() // Also closes s.listener.
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		level.Info(s.logger).Log("msg", "grace period exceeded enforcing shutdown")
		s.srv.Stop()
		return
	case <-stopped:
		cancel()
	}
	level.Info(s.logger).Log("msg", "internal server is shutdown gracefully", "err", err)
}
