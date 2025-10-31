/**
* Copyright mtch - 2025
 */

// Package consumer based on https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/consumer_example/consumer_example.go
// Function-based high-level Apache Kafka consumer
package consumer

import (
	"context"
	"errors"
	"fmt"
	kafkaCfg "github.com/Arclight-V/mtch/pkg/messagebroker/kafka"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/Arclight-V/mtch/pkg/messagebroker"
	"github.com/Arclight-V/mtch/pkg/platform/config"
)

type Consumer struct {
	c *kafka.Consumer

	logger log.Logger
	opts   options
}

func NewConsumer(cfg config.ConsumerConfig, logger log.Logger, opts ...Option) (*Consumer, error) {
	level.Info(logger).Log("msg", "creating kafka consumer")

	configMap := kafka.ConfigMap{
		kafkaCfg.EnableAutocommit:      false,
		kafkaCfg.EnableAutoOffsetStore: false,

		kafkaCfg.BrokerAddressFamily: kafkaCfg.DefaultBrokerAddressFamily,
		kafkaCfg.SessionTimeoutMS:    kafkaCfg.DefaultSessionTimeoutMS,
		kafkaCfg.MaxPollIntervalMs:   kafkaCfg.DefaultMaxPollIntervalMS,
		kafkaCfg.AutoOffsetReset:     kafkaCfg.DefaultAutoOffsetReset,
	}

	err := cfgToConsumerConfigMap(&cfg, &configMap)
	if err != nil {
		return nil, err
	}

	options := options{
		autoCommitTimeDuration: kafkaCfg.DefaultAutoCommitTimeDuration,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	if len(options.topics) == 0 {
		return nil, errors.New("no topics defined")
	}

	c, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka consumer: %w", err)
	}

	if err := c.SubscribeTopics(options.topics, nil); err != nil {
		return nil, fmt.Errorf("error subscribing to topics: %w", err)
	}

	return &Consumer{c: c, logger: logger, opts: options}, nil

}
func (c *Consumer) Start(ctx context.Context, h messagebroker.Handler) error {
	ticker := time.NewTicker(c.opts.autoCommitTimeDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			_, err := c.c.Commit()
			if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {

				return err
			}
			return nil
		case <-ticker.C:
			_, err := c.c.Commit()
			if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {

				return err
			}
		default:
			ev := c.c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				level.Info(c.logger).Log("msg", "Message on:", "topic", e.TopicPartition.Topic, "partition", e.TopicPartition.Partition, "value", string(e.Value))
				if e.Headers != nil {
					level.Info(c.logger).Log("msg", "Message headers on", "headers", e.Headers)
				}

				evt := messagebroker.Event{
					Topic:   *e.TopicPartition.Topic,
					Key:     e.Key,
					Value:   e.Value,
					Headers: toMap(e.Headers),
					Type:    getHeader(e.Headers, "event-type"),
				}
				if err := h(ctx, evt); err != nil {
					level.Error(c.logger).Log("msg", "handler error", "err", err)
					continue
				}

				// We can store the offsets of the messages manually or let
				// the library do it automatically based on the setting
				// enable.auto.offset.store. Once an offset is stored, the
				// library takes care of periodically committing it to the broker
				// if enable.auto.commit isn't set to false (the default is true).
				// By storing the offsets manually after completely processing
				// each message, we can ensure atleast once processing.
				if _, err := c.c.StoreMessage(e); err != nil {
					level.Error(c.logger).Log("msg", "store offset failed", "err", err)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				level.Error(c.logger).Log("err", e, "code", e.Code())
				if e.Code() == kafka.ErrAllBrokersDown {
					return fmt.Errorf("all kafka brokers down: %w", e)
				}
			default:
				level.Info(c.logger).Log("msg", "Unknown message type", "type", fmt.Sprintf("%T", e))
			}
		}
	}
}

func (c *Consumer) Close() error {
	level.Info(c.logger).Log("msg", "Closing consumer")

	return c.c.Close()
}

func cfgToConsumerConfigMap(cfg *config.ConsumerConfig, configMap *kafka.ConfigMap) error {
	if len(cfg.Brokers) == 0 {
		return errors.New("no brokers defined")
	}
	if cfg.GroupID == "" {
		return errors.New("no group id defined")
	}

	_ = configMap.SetKey(kafkaCfg.BootstrapServers, cfg.Brokers)
	_ = configMap.SetKey(kafkaCfg.GroupID, cfg.GroupID)
	if cfg.BrokerAddressFamily != kafkaCfg.DefaultBrokerAddressFamily {
		_ = configMap.SetKey(kafkaCfg.BootstrapServers, cfg.BrokerAddressFamily)
	}
	if cfg.SessionTimeoutMS != kafkaCfg.DefaultSessionTimeoutMS {
		_ = configMap.SetKey(kafkaCfg.SessionTimeoutMS, cfg.SessionTimeoutMS)
	}
	if cfg.MaxPollIntervalMs != kafkaCfg.DefaultMaxPollIntervalMS {
		_ = configMap.SetKey(kafkaCfg.MaxPollIntervalMs, cfg.MaxPollIntervalMs)
	}
	if cfg.AutoOffsetReset != kafkaCfg.DefaultAutoOffsetReset {
		_ = configMap.SetKey(kafkaCfg.AutoOffsetReset, cfg.AutoOffsetReset)
	}

	return nil
}

func toMap(headers []kafka.Header) map[string][]byte {
	m := make(map[string][]byte, len(headers))

	for _, h := range headers {
		m[h.Key] = h.Value
	}

	return m
}

func getHeader(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}
