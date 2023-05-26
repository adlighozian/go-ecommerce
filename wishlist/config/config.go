package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string `mapStructure:"DATABASE_URL"`
    Database    string `mapStructure:"PGDATABASE"`
    Host        string `mapStructure:"PGHOST"`
    User        string `mapStructure:"PGUSER"`
    Password    string `mapStructure:"PGPASSWORD"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile(".env")

    err := viper.ReadInConfig()
    if err != nil {
        _, ok := err.(viper.ConfigFileNotFoundError)
        if ok {
            return nil, errors.New(".env not found")
        }
        return nil, fmt.Errorf("fatal error config file %s", err)
    }

    config := Config{}
    err = viper.Unmarshal(&config)
    if err != nil {
        return nil, fmt.Errorf("fatal error decode : %s", err)
    }
    return &config, nil
}