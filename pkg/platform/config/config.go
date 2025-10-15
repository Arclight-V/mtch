package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server            *ServerCfg            `mapstructure:"server"`
	Client            *UserServiceClientCfg `mapstructure:"user_service_client"`
	Http              *HTTPCfg              `mapstructure:"http"`
	SMTPClient        *SMTPClient           `mapstructure:"smtp_client"`
	LogCfg            *LogCfg               `mapstructure:"logging"`
	UserServiceServer *UserServiceServerCfg `mapstructure:"user_service_server"`
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

// LoadConfig Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
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
	cfg, err := ParseConfig(cfgFile)
	return cfg, err
}
