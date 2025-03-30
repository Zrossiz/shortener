package redisdb

import (
	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"
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

func (r *RedisRepo) Get(hash string) (string, error) {
	val, err := r.client.Get(hash).Result()
	if err == redis.Nil {
		return "", apperrors.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}
