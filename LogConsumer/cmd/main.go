package main

import (
	"github.com/Zrossiz/LogConsumer/consumer/internal/delivery/kafka"
	"github.com/Zrossiz/LogConsumer/consumer/internal/repository/clickhouse"
	"github.com/Zrossiz/LogConsumer/consumer/internal/service"
	"github.com/Zrossiz/LogConsumer/consumer/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	clickhouseConn, err := clickhouse.ClickhouseConnect(cfg.Clickhouse.DBURI)
	if err != nil {
		panic(err)
	}

	clickhouseDB := clickhouse.NewClickhouse(clickhouseConn)

	service := service.New(clickhouseDB)

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Brokers)
	if err != nil {
		panic(err)
	}
	defer kafkaConsumer.Close()

	partitionConsumer, err := kafka.NewKafkaHandler(kafkaConsumer, cfg.Kafka.Topic, service)
	if err != nil {
		panic(err)
	}
	defer partitionConsumer.PartitionConsumer.Close()

	partitionConsumer.StartListening()
}
