package kafka

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/Zrossiz/LogConsumer/consumer/internal/domain"
)

type KafkaHandler struct {
	PartitionConsumer sarama.PartitionConsumer
	service           domain.Kafka
}

func NewKafkaHandler(
	consumer sarama.Consumer,
	topic string,
	serv domain.Kafka,
) (*KafkaHandler, error) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return &KafkaHandler{}, fmt.Errorf("error init partition consumer %v", err)
	}

	return &KafkaHandler{
		PartitionConsumer: partitionConsumer,
		service:           serv,
	}, nil
}

func (h *KafkaHandler) StartListening() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-h.PartitionConsumer.Messages():
			var messageData domain.RegisterRedirectEventDTO
			err := json.Unmarshal(msg.Value, &messageData)
			if err != nil {
				// error handling
			}

			h.service.Create(messageData)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
}
