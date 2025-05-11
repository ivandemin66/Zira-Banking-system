package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}

	Database struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}

	JWT struct {
		Host         string
		Port         int
		User         string
		Pass         string
		ExpiresHours time.Duration // Время жизни JWT токена в часах
	}

	CBR struct {
		URL string
	}
}

func LoadConfig(path string) *Config {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// интерпретируем jwt.ExpiresHours как часы
	cfg.JWT.ExpiresHours = cfg.JWT.ExpiresHours * time.Hour

	return &cfg
}
