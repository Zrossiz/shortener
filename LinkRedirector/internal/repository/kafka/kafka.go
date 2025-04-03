package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/Zrossiz/Redirector/redirector/internal/domain"
	"github.com/Zrossiz/Redirector/redirector/pkg/config"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func New(cfg config.Config) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error sync producer %v", err)
	}

	return &KafkaProducer{
		producer: producer,
		topic:    cfg.Kafka.Topic,
	}, nil
}

func (k *KafkaProducer) Send(message domain.UrlKafkaDTO) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
