package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"

	configkafka "github.com/Arclight-V/mtch/pkg/messagebroker/kafka"
)

type Config struct {
	Server                    *ServerCfg                    `mapstructure:"server"`
	Client                    *UserServiceClientCfg         `mapstructure:"user_service_client"`
	Http                      *HTTPCfg                      `mapstructure:"http"`
	SMTPClient                *SMTPClient                   `mapstructure:"smtp_client"`
	LogCfg                    *LogCfg                       `mapstructure:"logging"`
	UserServiceServer         *UserServiceServerCfg         `mapstructure:"user_service_server"`
	NotificationServiceServer *NotificationServiceServerCfg `mapstructure:"notification_server"`
	FrontEnd                  *FrontEndConfig               `mapstructure:"front_end"`
	Kafka                     *KafkaConfig                  `mapstructure:"kafka"`
	FlagD                     *FlagDConfig                  `mapstructure:"flagd"`
	Meta                      *Meta                         `mapstructure:"meta"`
}

type ServerCfg struct {
	AppVersion string `mapstructure:"app_version"`
	Port       string `mapstructure:"port"`
}

type UserServiceClientCfg struct {
	GRPCAddr string `mapstructure:"grpc_addr"`
}

type UserServiceServerCfg struct {
	Port        string        `mapstructure:"port"`
	GracePeriod time.Duration `mapstructure:"grace_period"`
}

// NotificationServiceServerCfg configuration parameters for the Notification service
type NotificationServiceServerCfg struct {
	Port        string        `mapstructure:"port"`
	GracePeriod time.Duration `mapstructure:"grace_period"`
}

type HTTPCfg struct {
	ListenAddr        string `mapstructure:"listen_addr"`
	MetricsListenAddr string `mapstructure:"metrics_listen_addr"`
}

type SMTPClient struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"userservice"`
	Pass string `mapstructure:"pass"`
	From string `mapstructure:"from"`
}

type LogCfg struct {
	Level     string `mapstructure:"level"`
	Format    string `mapstructure:"format"`
	DebugName string `mapstructure:"debug_name"`
}

type FrontEndConfig struct {
	FrontendPath string `mapstructure:"frontend_path"`
}

type KafkaConfig struct {
	Producer ProducerConfig `mapstructure:"producer"`
	Consumer ConsumerConfig `mapstructure:"consumer"`
}

type CommonKafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	ClientID string   `mapstructure:"client_id"`
}
type ProducerConfig struct {
	CommonKafkaConfig `mapstructure:",squash"`
	CompressionType   string `mapstructure:"compression_type"`
	Acks              int    `mapstructure:"acks"`
	LingerMS          int    `mapstructure:"linger_ms"`
	FlushTimeoutMS    int    `mapstructure:"flush_timeout_ms"`
	EnableIdempotence bool   `mapstructure:"enable_idempotence"`
}

type ConsumerConfig struct {
	CommonKafkaConfig   `mapstructure:",squash"`
	BrokerAddressFamily string `mapstructure:"broker_address_family"`
	GroupID             string `mapstructure:"group_id"`
	SessionTimeoutMS    int64  `mapstructure:"session_timeout_ms"`
	MaxPollIntervalMs   int64  `mapstructure:"max_poll_interval_ms"`
	AutoOffsetReset     string `mapstructure:"auto_offset_reset"`
}

type FlagDConfig struct {
	FlagsPath string `mapstructure:"flags_path"`
}

type Meta struct {
	Template bool `mapstructure:"template"`
}

// LoadConfig Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	//Kafka default values
	// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	// producer
	_ = v.BindEnv("kafka.producer.brokers", "KAFKA_BOOTSTRAP")
	v.SetDefault("kafka.producer.acks", configkafka.DefaultAcks)
	v.SetDefault("kafka.producer.linger_ms", configkafka.DefaultLingerMS)
	v.SetDefault("kafka.producer.flush_timeout_ms", configkafka.DefaultFlushTimoutMS)
	v.SetDefault("kafka.producer.compression_type", configkafka.DefaultCompressionType)

	//consumer
	_ = v.BindEnv("kafka.consumer.brokers", "KAFKA_BOOTSTRAP")

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// ParseConfig Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// GetConfig Get config
func GetConfig(path string) (*Config, error) {
	cfgFile, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	var cfg *Config

	if cfg, err = ParseConfig(cfgFile); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, err
}

// Validate validates the structure that it is not a template.
func (c *Config) Validate() error {
	if c.Meta != nil && c.Meta.Template == true {
		return errors.New("config is a template: please fill required fields (ports, etc.)")
	}

	return nil
}
