package redisdb

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisRepo struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func Connect(addr, password string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func (r *RedisRepo) Create(hash, original string) error {
	err := r.client.Set(hash, original, time.Hour*1).Err()
	if err != nil {
		return fmt.Errorf("error create redis url: %v", err)
	}
	return nil
}
