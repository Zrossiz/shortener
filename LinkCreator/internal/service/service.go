package service

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/Zrossiz/LinkCreator/creator/internal/domain"
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
