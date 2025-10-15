package grpc

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"time"
)

type options struct {
	registerServerFuncs []registerServerFunc

	gracePeriod time.Duration

	listen  string
	network string

	tlsConfig *tls.Config

	grpcOpts []grpc.ServerOption
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

type registerServerFunc func(s *grpc.Server)

func WithServer(f registerServerFunc) Option {
	return optionFunc(func(o *options) {
		o.registerServerFuncs = append(o.registerServerFuncs, f)
	})
}

func WithListen(s string) Option {
	return optionFunc(func(o *options) {
		o.listen = s
	})
}

func WithTLSConfig(cfg *tls.Config) Option {
	return optionFunc(func(o *options) {
		o.tlsConfig = cfg
	})
}

func WithGRPCOptions(opt grpc.ServerOption) Option {
	return optionFunc(func(o *options) {
		o.grpcOpts = append(o.grpcOpts, opt)
	})
}

func WithGracePeriod(d time.Duration) Option {
	return optionFunc(func(o *options) {
		o.gracePeriod = d
	})
}
