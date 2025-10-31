package producer

type options struct {
	compressionType   string
	acks              int
	lingerMS          int
	flushTimeoutMS    int
	enableIdempotence bool
}

// Option overrides behavior of Consumer.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithAcks(acks int) Option {
	return optionFunc(func(o *options) {
		o.acks = acks
	})
}

func WithEnableIdempotence(enableIdempotence bool) Option {
	return optionFunc(func(o *options) {
		o.enableIdempotence = enableIdempotence
	})
}

func WithLingerMS(lingerMS int) Option {
	return optionFunc(func(o *options) {
		o.lingerMS = lingerMS
	})
}

func WithCompressionType(compressionType string) Option {
	return optionFunc(func(o *options) {
		o.compressionType = compressionType
	})
}

func WithFlushTimeoutMS(flushTimeoutMS int) Option {
	return optionFunc(func(o *options) {
		o.flushTimeoutMS = flushTimeoutMS
	})
}
