package config

import (
	"github.com/joho/godotenv"
	"os"
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

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	cfg.Server.LogLevel = os.Getenv("SERVER_LOG_LEVEL")
	cfg.Redis.Address = os.Getenv("REDIS_ADDRESS")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Postgres.DBURI = os.Getenv("POSTGRES_DB_URI")

	return &cfg
}
