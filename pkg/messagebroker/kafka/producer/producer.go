package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/Arclight-V/mtch/pkg/messagebroker"
	kafkaCfg "github.com/Arclight-V/mtch/pkg/messagebroker/kafka"
	"github.com/Arclight-V/mtch/pkg/platform/config"
)

// A Producer defines parameters for kafka.Producer, a wrapper around kafka.Producer.
type Producer struct {
	p *kafka.Producer

	logger log.Logger

	opts options
}

// New creates a new Producer.
func New(cfg config.ProducerConfig, logger log.Logger, opts ...Option) (*Producer, error) {
	level.Info(logger).Log("msg", "creating kafka producer")
	configMap := cfgToProducerConfigMap(cfg)

	options := options{
		flushTimeoutMS: kafkaCfg.DefaultFlushTimoutMS,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(&options)
		}

		setKeyFromOpts(configMap, &options)
	}

	p, err := kafka.NewProducer(configMap)

	if err != nil {
		return nil, fmt.Errorf("error creating kafka producer: %w", err)
	}

	return &Producer{p: p, logger: logger}, nil
}

// Publish a single message.
func (p *Producer) Publish(ctx context.Context, event *messagebroker.Event) error {
	headers := make([]kafka.Header, 0, len(event.Headers))
	for k, v := range event.Headers {
		headers = append(headers, kafka.Header{Key: k, Value: v})
	}

	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &event.Topic,
			Partition: kafka.PartitionAny,
		},
		Value:   event.Value,
		Headers: headers,
	}
	err := p.p.Produce(&msg, nil)

	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			//TODO:backpressure
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			time.Sleep(time.Second)
			level.Info(p.logger).Log("msg", "kafka producer is dead, retrying in 1 seconds", "err", err, "topic", "TODO::backpressure")
		}
		level.Error(p.logger).Log("msg", "Failed to produce event:", "err", err)
		return err
	}

	return nil
}

// Close a Producer instance.
func (p *Producer) Close() error {
	for p.p.Flush(p.opts.flushTimeoutMS) > 0 {
		level.Info(p.logger).Log("msg", "Waiting for messages to flush")
	}
	p.p.Close()
	level.Info(p.logger).Log("msg", "Producer closed")
	return nil
}

// cfgToProducerConfigMap converts a ProducerConfig struct into a *kafka.ConfigMap
// suitable for initializing a Kafka producer.
func cfgToProducerConfigMap(cfg config.ProducerConfig) *kafka.ConfigMap {
	configMap := &kafka.ConfigMap{
		kafkaCfg.BootstrapServers: cfg.Brokers,
		kafkaCfg.ClientID:         cfg.ClientID,
	}

	return configMap
}

// setKeyFromOpts set optional parameters into a *kafka.ConfigMap
func setKeyFromOpts(configMap *kafka.ConfigMap, options *options) {
	if options.acks != kafkaCfg.DefaultAcks {
		_ = configMap.SetKey(kafkaCfg.Acks, options.acks)
	}
	if options.enableIdempotence != kafkaCfg.DefaultEnableIdempotence {
		_ = configMap.SetKey(kafkaCfg.EnableIdempotence, options.enableIdempotence)
	}
	if options.lingerMS != kafkaCfg.DefaultLingerMS {
		_ = configMap.SetKey(kafkaCfg.LingerMS, options.lingerMS)
	}
	if options.compressionType != kafkaCfg.DefaultCompressionType {
		_ = configMap.SetKey(kafkaCfg.CompressionType, options.compressionType)
	}
}
