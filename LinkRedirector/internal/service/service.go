package service

import (
	"errors"
	"github.com/Zrossiz/Redirector/redirector/internal/domain"
	"github.com/Zrossiz/Redirector/redirector/pkg/apperrors"

	"go.uber.org/zap"
)

type Service struct {
	log           *zap.Logger
	postgresStore domain.UrlPostresRepo
	redisStore    domain.UrlRedisRepo
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewService(
	log *zap.Logger,
	postgresStore domain.UrlPostresRepo,
	redisStore domain.UrlRedisRepo,
) *Service {
	return &Service{
		log:           log,
		postgresStore: postgresStore,
		redisStore:    redisStore,
	}
}

func (s *Service) Get(hash string) (string, error) {
	originalURLRedis, err := s.redisStore.Get(hash)
	if err == nil && originalURLRedis != "" {
		return originalURLRedis, nil
	}
	if err != nil && errors.Is(err, apperrors.ErrNotFound) {
		s.log.Error("Failed to get redis store", zap.String("hash", hash), zap.Error(err))
	}

	original, err := s.postgresStore.Get(hash)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	return original, nil
}
