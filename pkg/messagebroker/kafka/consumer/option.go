package consumer

import "time"

type options struct {
	topics                 []string
	autoCommitTimeDuration time.Duration
}

// Option overrides behavior of Consumer.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithTopics sets topics for consumer
func WithTopics(topics ...string) Option {
	return optionFunc(func(o *options) {
		o.topics = append(o.topics, topics...)
	})
}

// WithAutoCommitTimeDuration - sets the frequency in milliseconds that the consumer offsets are committed (written)
// to offset storage. (0 = disable). This setting is used by the high-level consumer.
func WithAutoCommitTimeDuration(autoCommitTimeDuration time.Duration) Option {
	return optionFunc(func(o *options) {
		o.autoCommitTimeDuration = autoCommitTimeDuration
	})
}
