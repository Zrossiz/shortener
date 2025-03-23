package redisdb

import (
	"fmt"
	"time"

	"github.com/Zrossiz/shortener/internal/config"
	"github.com/go-redis/redis"
)

type RedisStore struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func Connect(cfg *config.Config) (*redis.Client, error) {
	fmt.Println("Connecting to Redis...")
	fmt.Println("address: ", cfg.RedisAddress)
	fmt.Println("password: ", cfg.RedisPassword)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func (r *RedisStore) Create(hash, original string) error {
	err := r.client.Set(hash, original, time.Hour*1).Err()
	if err != nil {
		return fmt.Errorf("error create redis url: %v", err)
	}
	return nil
}

func (r *RedisStore) Get(hash string) (string, error) {
	val, err := r.client.Get(hash).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
