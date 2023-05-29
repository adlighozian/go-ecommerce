package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database       string `mapstructure:"DATABASE"`
	DatabaseDriver string `mapstructure:"DATABASE_DRIVER"`
	DatabaseURL    string `mapstructure:"DATABASE_URL"`
	PgDatabase     string `mapstructure:"PGDATABASE"`
	PgHost         string `mapstructure:"PGHOST"`
	PgPassword     string `mapstructure:"PGPASSWORD"`
	PgPort         string `mapstructure:"PGPORT"`
	PgUser         string `mapstructure:"PGUSER"`
	Port           string `mapstructure:"PORT"`
	Debug          bool   `mapstructure:"DEBUG"`
	RabbitMQ       string `mapstructure:"RABBITMQ"`
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
