package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

func NewKafkaConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error init consumer %v", err)
	}

	return consumer, nil
}
