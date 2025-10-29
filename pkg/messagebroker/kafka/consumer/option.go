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

func WithTopics(topics ...string) Option {
	return optionFunc(func(o *options) {
		o.topics = append(o.topics, topics...)
	})
}

func WithAutoCommitTimeDuration(autoCommitTimeDuration time.Duration) Option {
	return optionFunc(func(o *options) {
		o.autoCommitTimeDuration = autoCommitTimeDuration
	})
}
