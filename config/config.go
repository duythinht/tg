package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		DSN string `yaml:"dsn" mapstructure:"dsn"`
	} `yaml:"db" mapstructure:"db"`
	Queue struct {
		Brokers string `yaml:"brokers" mapstructure:"brokers"`
		Topic   string `yaml:"topic" mapstructure:"topic"`
		GroupID string `yaml:"group_id" mapstructure:"group_id"`
	} `yaml:"queue" mapstructure:"queue"`
}

func Load() (*Config, error) {
	c := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/opt/tg/config")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config %w", err)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config %w", err)
	}

	return c, nil
}
