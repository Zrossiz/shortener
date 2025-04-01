package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Address  string
		LogLevel string
	}
	Kafka struct {
		Brokers        []string
		Topic          string
		GroupID        string
		MinBytes       int
		MaxBytes       int
		MaxWait        time.Duration
		CommitInterval time.Duration
	}
	Clickhouse struct {
		DBURI string
	}
}

func LoadConfig() *Config {
	var cfg Config

	_ = godotenv.Load()

	cfg.Server.Address = getStringEnvOrDefault("SERVER_ADDRESS", "localhost:8080")
	cfg.Server.LogLevel = getStringEnvOrDefault("SERVER_LOG_LEVEL", "warn")

	cfg.Clickhouse.DBURI = getStringEnvOrDefault("DB_URI", "invalid")

	cfg.Kafka.Brokers = strings.Split(getStringEnvOrDefault("KAFKA_BROKERS", "localhost:9090,localhost:9000"), ",")
	cfg.Kafka.Topic = getStringEnvOrDefault("KAFKA_TOPIC", "url_events")
	cfg.Kafka.GroupID = getStringEnvOrDefault("KAFKA_GROUP_ID", "url_group")
	cfg.Kafka.MinBytes = getIntEnvOrDefault("KAFKA_MIN_BYTES", 10e3)
	cfg.Kafka.MaxBytes = getIntEnvOrDefault("KAFKA_MAX_BYTES", 10e6)
	cfg.Kafka.MaxWait = getDurationEnvOrDefault("KAFKA_MAX_WAIT", 500*time.Millisecond)
	cfg.Kafka.CommitInterval = getDurationEnvOrDefault("KAFKA_COMMIT_INTERVAL", 0)

	return &cfg
}

func getStringEnvOrDefault(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return defaultValue
	}

	return envValue
}

func getIntEnvOrDefault(envName string, defaultValue int) int {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(envValue)
	if err != nil {
		fmt.Printf("error parsing %v, %v\n", envValue, err)
		return defaultValue
	}

	return intValue
}

func getDurationEnvOrDefault(envName string, defaultValue time.Duration) time.Duration {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(envValue)
	if err != nil {
		fmt.Printf("error parsing %v, %v\n", envValue, err)
		return defaultValue
	}

	return durationValue
}
