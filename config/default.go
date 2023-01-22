package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort          string `mapstructure:"SERVER_PORT"`
	DBUri               string `mapstructure:"POSTGRESQL_URL"`
	Port                string `mapstructure:"POSTGRESQL_PORT"`
	JWTSecretKey        string `mapstructure:"JWT_SECRET_KEY"`
	JWTRefreshSecretKey string `mapstructure:"JWT_REFRESH_SECRET_KEY"`
}

var (
	AppConfig *Config
)

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
