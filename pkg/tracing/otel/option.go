package otel

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type options struct {
	serviceName string
	exporter    trace.SpanExporter
	sampler     trace.Sampler
	attributes  []attribute.KeyValue
}

type Option interface {
	apply(*options)
}

type OptionFunc func(*options)

func (f OptionFunc) apply(o *options) {
	f(o)
}

func WithServiceName(name string) Option {
	return OptionFunc(func(o *options) {
		o.serviceName = name
	})
}

func WithExporter(exporter trace.SpanExporter) Option {
	return OptionFunc(func(o *options) {
		o.exporter = exporter
	})
}

func WithSampler(sampler trace.Sampler) Option {
	return OptionFunc(func(o *options) {
		o.sampler = sampler
	})
}

func WithAttributes(attributes ...attribute.KeyValue) Option {
	return OptionFunc(func(o *options) {
		o.attributes = append(attributes, o.attributes...)
	})
}
