package kafka

import (
	"encoding/json"
	"fmt"
	"time"

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

	var producer sarama.SyncProducer
	var err error

	maxRetries := 5
	retryInterval := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
		if err == nil {
			fmt.Print("Kafka producer initiated\n")
			return &KafkaProducer{
				producer: producer,
				topic:    cfg.Kafka.Topic,
			}, nil
		}

		fmt.Printf("Attempt %d/%d to connect to Kafka failed: %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("error creating sync producer after %d retries: %v", maxRetries, err)
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
