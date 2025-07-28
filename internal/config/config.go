package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	cfg := &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		DBHost:     viper.GetString("DATABASE_HOST"),
		DBPort:     viper.GetString("DATABASE_PORT"),
		DBUser:     viper.GetString("DATABASE_USER"),
		DBPassword: viper.GetString("DATABASE_PASSWORD"),
		DBName:     viper.GetString("DATABASE_NAME"),
	}
	return cfg, nil
}
