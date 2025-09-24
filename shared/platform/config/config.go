package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server *ServerCfg `mapstructure:"server"`
	Client *ClientCfg `mapstructure:"client"`
	Http   *HTTPCfg   `mapstructure:"http"`
}

type ServerCfg struct {
	AppVersion string `mapstructure:"app_version"`
	Port       string `mapstructure:"port"`
}

type ClientCfg struct {
	GRPCAddr string `mapstructure:"grpc_addr"`
}

type HTTPCfg struct {
	HTTPAddr string `mapstructure:"http_addr"`
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
