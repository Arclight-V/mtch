package postgres

import (
	"context"
	"time"
)

type options struct {
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	maxOpenConns    int
	maxIdleConns    int
	context         context.Context
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) { f(o) }

func WithConnMaxLifetime(d time.Duration) Option {
	return optionFunc(func(o *options) {
		o.connMaxLifetime = d
	})
}

func WithConnMaxIdleTime(d time.Duration) Option {
	return optionFunc(func(o *options) {
		o.connMaxIdleTime = d
	})
}

func WithMaxOpenConns(c int) Option {
	return optionFunc(func(o *options) {
		o.maxOpenConns = c
	})
}

func WithMaxIdleConns(c int) Option {
	return optionFunc(func(o *options) {
		o.maxIdleConns = c
	})
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(o *options) {
		o.context = ctx
	})
}
