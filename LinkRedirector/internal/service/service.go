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
	kafkaProducer domain.UrlKafka
}

func NewService(
	log *zap.Logger,
	postgresStore domain.UrlPostresRepo,
	redisStore domain.UrlRedisRepo,
	kafkaProducer domain.UrlKafka,
) *Service {
	return &Service{
		log:           log,
		postgresStore: postgresStore,
		redisStore:    redisStore,
		kafkaProducer: kafkaProducer,
	}
}

func (s *Service) Get(data domain.GetUrlDTO) (string, error) {
	originalURLRedis, err := s.redisStore.Get(data.Short)
	if err == nil && originalURLRedis != "" {
		return originalURLRedis, nil
	}
	if err != nil && errors.Is(err, apperrors.ErrNotFound) {
		s.log.Error("Failed to get redis store", zap.String("hash", data.Short), zap.Error(err))
	}

	original, err := s.postgresStore.Get(data.Short)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}

	err = s.kafkaProducer.Send(domain.UrlKafkaDTO{
		Original: original,
		Short:    data.Short,
		OS:       data.OS,
		UserIP:   data.UserIP,
	})
	if err != nil {
		s.log.Error(err.Error())
	}

	return original, nil
}
