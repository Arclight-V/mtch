package kafka

import "time"

// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md

// default values in Kafka
const (
	DefaultAcks                = -1
	DefaultBrokerAddressFamily = "v4"
	DefaultCompressionType     = "none"
	DefaultLingerMS            = 5
	DefaultEnableIdempotence   = false
)

// default values for Kafka
const (
	DefaultFlushTimoutMS          = 10_000
	DefaultSessionTimeoutMS       = 6000
	DefaultMaxPollIntervalMS      = 600000
	DefaultAutoOffsetReset        = "earliest"
	DefaultAutoCommitTimeDuration = 5 * time.Second
)

// Global Configuration Properties
const (
	Acks                  = "acks"
	BootstrapServers      = "bootstrap.servers"
	ClientID              = "client.id"
	BrokerAddressFamily   = "broker.address.family"
	LingerMS              = "linger.ms"
	EnableIdempotence     = "enable.idempotence"
	CompressionType       = "compression.type"
	SessionTimeoutMS      = "session.timeout.ms"
	MaxPollIntervalMs     = "max.poll.interval.ms"
	AutoOffsetReset       = "auto.offset.reset"
	GroupID               = "group.id"
	EnableAutoOffsetStore = "enable.auto.offset.store"
	EnableAutocommit      = "enable.auto.commit"
)
