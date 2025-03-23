package service

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/Zrossiz/shortener/internal/apperrors"

	"go.uber.org/zap"
)

type Service struct {
	log           *zap.Logger
	postgresStore PostgresStorage
	redisStore    RedisStorage
}

type PostgresStorage interface {
	Create(url string, hash string) error
	Get(hash string) (string, error)
}

type RedisStorage interface {
	Create(hash, original string) error
	Get(hash string) (string, error)
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewService(store PostgresStorage, redisStore RedisStorage, log *zap.Logger) Service {
	return Service{
		postgresStore: store,
		redisStore:    redisStore,
		log:           log,
	}
}

func (s *Service) Get(hash string) (string, error) {
	originalURLRedis, err := s.redisStore.Get(hash)
	if err == nil && originalURLRedis != "" {
		return originalURLRedis, nil
	}
	if err != nil && errors.Is(err, apperrors.ErrRedisNotFound) {
		s.log.Error("Failed to get redis store", zap.String("hash", hash), zap.Error(err))
	}

	original, err := s.postgresStore.Get(hash)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	return original, nil
}

func (s *Service) Create(original string) (string, error) {
	hash := s.generateHash(original)

	err := s.postgresStore.Create(original, hash)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	err = s.redisStore.Create(hash, original)
	if err != nil {
		s.log.Error(err.Error())
	}

	return hash, nil
}

func (s *Service) generateHash(original string) string {
	hash := sha256.Sum256([]byte(original))

	num := binary.BigEndian.Uint64(hash[:8])

	shortURL := make([]byte, 7)
	for i := 0; i < 7; i++ {
		shortURL[i] = base62Chars[num%62]
		num /= 62
	}

	return string(shortURL)
}
