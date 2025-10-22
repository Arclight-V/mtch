package grpc

import (
	"context"
	"fmt"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"net"
	"runtime/debug"

	logging_mw "github.com/Arclight-V/mtch/pkg/logging"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
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

func NewServer(logger log.Logger, reg prometheus.Registerer, logOpts []grpc_logging.Option, probe *prober.GRPCProbe, opts ...Option) *Server {
	logger = log.With(logger, "service", "gRPC/server")
	options := options{
		network: "tcp",
	}
	for _, o := range opts {
		o.apply(&options)
	}

	met := grpc_prometheus.NewServerMetrics(
		grpc_prometheus.WithServerHandlingTimeHistogram(
			grpc_prometheus.WithHistogramOpts(&prometheus.HistogramOpts{
				Buckets:                        []float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120},
				NativeHistogramMaxBucketNumber: 256,
				NativeHistogramBucketFactor:    1.1,
			}),
		),
	)

	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})

	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		level.Error(logger).Log("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
		return status.Errorf(codes.Internal, "%s", p)
	}

	options.grpcOpts = append(options.grpcOpts, []grpc.ServerOption{
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.ChainUnaryInterceptor(
			NewUnaryServerRequestIDInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
			met.UnaryServerInterceptor(),
			grpc_logging.UnaryServerInterceptor(logging_mw.InterceptorLogger(logger), logOpts...),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}...)

	if options.tlsConfig != nil {
		options.grpcOpts = append(options.grpcOpts, grpc.Creds(credentials.NewTLS(options.tlsConfig)))
	}

	s := grpc.NewServer(options.grpcOpts...)

	// Register all configured servers.
	for _, f := range options.registerServerFuncs {
		f(s)
	}

	met.InitializeMetrics(s)
	reg.MustRegister(met)

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
