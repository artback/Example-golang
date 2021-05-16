package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource       string `mapstructure:"DB_SOURCE"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	ServiceAddress string `mapstructure:"SERVICE_ADDRESS"`
}

func LoadConfig() (config Config, err error) {
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("SERVICE_ADDRESS")
	err = viper.Unmarshal(&config)
	return
}
