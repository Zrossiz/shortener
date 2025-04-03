package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {

	os.Clearenv()

	cfg := LoadConfig()

	assert.Equal(t, "localhost:8080", cfg.Server.Address)
	assert.Equal(t, "warn", cfg.Server.LogLevel)
	assert.Equal(t, "", cfg.Redis.Address)
	assert.Equal(t, "root", cfg.Redis.Password)
	assert.Equal(t, "invalid", cfg.Postgres.DBURI)
}

func TestLoadConfig_FromEnv(t *testing.T) {

	os.Setenv("SERVER_ADDRESS", "127.0.0.1:9090")
	os.Setenv("SERVER_LOG_LEVEL", "debug")
	os.Setenv("REDIS_ADDRESS", "redis://localhost:6379")
	os.Setenv("REDIS_PASSWORD", "my-secret-password")
	os.Setenv("POSTGRES_DB_URI", "postgresql://user:password@localhost/dbname")

	cfg := LoadConfig()

	assert.Equal(t, "127.0.0.1:9090", cfg.Server.Address)
	assert.Equal(t, "debug", cfg.Server.LogLevel)
	assert.Equal(t, "redis://localhost:6379", cfg.Redis.Address)
	assert.Equal(t, "my-secret-password", cfg.Redis.Password)
	assert.Equal(t, "postgresql://user:password@localhost/dbname", cfg.Postgres.DBURI)
}

func TestLoadConfig_ClearEnv(t *testing.T) {

	os.Setenv("SERVER_ADDRESS", "127.0.0.1:8082")

	cfg := LoadConfig()

	assert.Equal(t, "127.0.0.1:8082", cfg.Server.Address)

	os.Clearenv()
}
