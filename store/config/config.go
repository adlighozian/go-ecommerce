package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug 		bool   `mapstructure:"DEBUG"`
	Port  		string `mapstructure:"PORT"`
	RabbitMQURL string `mapstructure:"RABBITMQURL"`

	Database `mapstructure:",squash"`
}

type Database struct {
	Driver string `mapstructure:"DATABASE_DRIVER"`
	URL    string `mapstructure:"DATABASE_URL"`
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