package service

import (
	"fmt"

	"github.com/Zrossiz/LogConsumer/consumer/internal/domain"
)

type KafkaService struct {
	db domain.ClickhouseDB
}

func New(db domain.ClickhouseDB) *KafkaService {
	return &KafkaService{db: db}
}

func (k *KafkaService) Create(data domain.RegisterRedirectEventDTO) error {
	err := k.db.Create(data)
	if err != nil {
		return fmt.Errorf("error create event: %v", err)
	}

	return nil
}
