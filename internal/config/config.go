package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	SMTPHost   string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	PGPKeyPath string
}

func LoadConfig() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "bank_api"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		SMTPHost:   getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:   smtpPort,
		SMTPUser:   getEnv("SMTP_USER", "noreply@example.com"),
		SMTPPass:   getEnv("SMTP_PASS", "strong_password"),
		PGPKeyPath: getEnv("PGP_KEY_PATH", "./keys"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
