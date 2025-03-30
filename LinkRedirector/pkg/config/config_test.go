package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	originalEnv := map[string]string{
		"SERVER_ADDRESS":   os.Getenv("SERVER_ADDRESS"),
		"SERVER_LOG_LEVEL": os.Getenv("SERVER_LOG_LEVEL"),
		"REDIS_ADDRESS":    os.Getenv("REDIS_ADDRESS"),
		"REDIS_PASSWORD":   os.Getenv("REDIS_PASSWORD"),
		"POSTGRES_DB_URI":  os.Getenv("POSTGRES_DB_URI"),
	}

	_ = os.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	_ = os.Setenv("SERVER_LOG_LEVEL", "debug")
	_ = os.Setenv("REDIS_ADDRESS", "localhost:6379")
	_ = os.Setenv("REDIS_PASSWORD", "secret")
	_ = os.Setenv("POSTGRES_DB_URI", "postgres://user:pass@localhost:5432/db")

	cfg := LoadConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "127.0.0.1:8080", cfg.Server.Address)
	assert.Equal(t, "debug", cfg.Server.LogLevel)
	assert.Equal(t, "localhost:6379", cfg.Redis.Address)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Postgres.DBURI)

	for key, value := range originalEnv {
		if value == "" {
			_ = os.Unsetenv(key)
		} else {
			_ = os.Setenv(key, value)
		}
	}
}
