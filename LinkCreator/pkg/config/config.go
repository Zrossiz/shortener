package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Address  string
		LogLevel string
	}
	Redis struct {
		Address  string
		Password string
	}
	Postgres struct {
		DBURI string
	}
}

func LoadConfig() *Config {
	var cfg Config

	_ = godotenv.Load()

	cfg.Server.Address = getStringEnvOrDefault("SERVER_ADDRESS", "localhost:4000")
	cfg.Server.LogLevel = getStringEnvOrDefault("SERVER_LOG_LEVEL", "warn")
	cfg.Redis.Address = getStringEnvOrDefault("REDIS_ADDRESS", "localhost:6379")
	cfg.Redis.Password = getStringEnvOrDefault("REDIS_PASSWORD", "root")
	cfg.Postgres.DBURI = getStringEnvOrDefault("POSTGRES_DB_URI", "invalid")

	return &cfg
}

func getStringEnvOrDefault(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}
