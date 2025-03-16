package service

import (
	"crypto/sha256"
	"encoding/binary"

	"go.uber.org/zap"
)

type Service struct {
	log   *zap.Logger
	store Storage
}

type Storage interface {
	Create(url string, hash string) (string, error)
	Get(hash string) (string, error)
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewService(store Storage, log *zap.Logger) Service {
	return Service{store: store, log: log}
}

func (s *Service) Get(hash string) (string, error) {
	original, err := s.store.Get(hash)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	return original, nil
}

func (s *Service) Create(original string) (string, error) {
	hash := s.generateHash(original)

	_, err := s.store.Create(original, hash)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
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
