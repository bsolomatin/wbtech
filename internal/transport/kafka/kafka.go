package kafka

import (
	"dockertest/internal/config"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewReader(cfg config.ConsumerConfig) *Consumer{
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Host + ":" + fmt.Sprint(cfg.Port)},
			Topic:   cfg.Topic,
			GroupID: "wb-group",
		}),
		//Processor: processor,
	}
}