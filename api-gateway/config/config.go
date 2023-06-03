package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool   `mapstructure:"DEBUG"`
	Port  string `mapstructure:"PORT"`

	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`

	Database `mapstructure:",squash"`
	Redis    `mapstructure:",squash"`
}

type Database struct {
	Driver string `mapstructure:"DATABASE_DRIVER"`
	URL    string `mapstructure:"DATABASE_URL"`
	DBName string `mapstructure:"PGDATABASE"`
}

type Redis struct {
	Addr       string `mapstructure:"REDIS_ADDR"`
	ClientName string `mapstructure:"REDIS_CLIENTNAME"`
	Username   string `mapstructure:"REDIS_USERNAME"`
	Password   string `mapstructure:"REDIS_PASSWORD"`
	DB         int    `mapstructure:"REDIS_DB"`
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
