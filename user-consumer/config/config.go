package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool `mapstructure:"DEBUG"`

	Database `mapstructure:",squash"`
	RabbitMQ `mapstructure:",squash"`
}

type Database struct {
	Driver string `mapstructure:"DATABASE_DRIVER"`
	URL    string `mapstructure:"DATABASE_URL"`
	DBName string `mapstructure:"PGDATABASE"`
}

type RabbitMQ struct {
	URL string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	errRead := viper.ReadInConfig()
	if errRead != nil {
		return nil, errRead
	}

	config := new(Config)
	errUn := viper.Unmarshal(&config)
	if errUn != nil {
		return nil, errUn
	}
	return config, nil
}
