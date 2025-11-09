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

// WithAcks sets acks. This field indicates the number of acknowledgements the leader broker must receive from ISR
// brokers before responding to the request: 0=Broker does not send any response/ack to client, -1 or all=Broker will
// block until message is committed by all in sync replicas (ISRs). If there are less than min.insync.replicas
// (broker configuration) in the ISR set the produce request will fail.
func WithAcks(acks int) Option {
	return optionFunc(func(o *options) {
		o.acks = acks
	})
}

// WithEnableIdempotence set enable.idempotence.
// When set to true, the producer will ensure that messages are successfully produced exactly once and in the original
// produce order. The following configuration properties are adjusted automatically (if not modified by the user) when
// idempotence is enabled: max.in.flight.requests.per.connection=5 (must be less than or equal to 5), retries=INT32_MAX
// (must be greater than 0), acks=all, queuing.strategy=fifo.
// Producer instantation will fail if user-supplied configuration is incompatible.
func WithEnableIdempotence(enableIdempotence bool) Option {
	return optionFunc(func(o *options) {
		o.enableIdempotence = enableIdempotence
	})
}

// WithLingerMS sets linger.ms. Delay in milliseconds to wait for messages in the producer queue to accumulate before
// constructing message batches (MessageSets) to transmit to brokers. A higher value allows larger and more effective
// (less overhead, improved compression) batches of messages to accumulate at the expense of increased message delivery
// latency.
func WithLingerMS(lingerMS int) Option {
	return optionFunc(func(o *options) {
		o.lingerMS = lingerMS
	})
}

// WithCompressionType sets compression.type. Compression codec to use for compressing message sets. This is the default
// value for all topics, may be overridden by the topic configuration property compression.codec.
func WithCompressionType(compressionType string) Option {
	return optionFunc(func(o *options) {
		o.compressionType = compressionType
	})
}

// WithFlushTimeoutMS sets timeout for Flush
func WithFlushTimeoutMS(flushTimeoutMS int) Option {
	return optionFunc(func(o *options) {
		o.flushTimeoutMS = flushTimeoutMS
	})
}
