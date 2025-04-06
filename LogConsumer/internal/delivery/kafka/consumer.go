package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

func NewKafkaConsumer(brokers []string) (sarama.Consumer, error) {
	var consumer sarama.Consumer
	var err error
	maxRetries := 5
	retryInterval := time.Second * 5

	for i := 0; i < maxRetries; i++ {
		consumer, err = sarama.NewConsumer(brokers, sarama.NewConfig())
		if err == nil {
			return consumer, nil
		}

		fmt.Printf("Attempt %d/%d to connect to Kafka failed: %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("error init consumer after %d retries: %v", maxRetries, err)
}
