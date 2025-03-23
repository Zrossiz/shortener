package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURI         string
	LogLevel      string
	ServerPort    string
	RedisAddress  string
	RedisPassword string
}

func GetConfig() (*Config, error) {
	var cfg Config

	_ = godotenv.Load()

	cfg.DBURI = getStringEnvOrDefault("DB_URI", "empty db uri")
	cfg.LogLevel = getStringEnvOrDefault("LOG_LEVEL", "debug")
	cfg.ServerPort = getStringEnvOrDefault("SERVER_PORT", ":8080")
	cfg.RedisAddress = getStringEnvOrDefault("REDIS_ADDRESS", "redis:6379")

	return &cfg, nil
}

func getStringEnvOrDefault(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return envValue
	}

	return defaultValue
}
