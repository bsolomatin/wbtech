package kafka

import (
	"context"
	"dockertest/internal/config"
	"dockertest/pkg/logger"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Processor func (ctx context.Context, msg kafka.Message) error

type Consumer struct {
	Reader *kafka.Reader
	Processor Processor
	Logger logger.Logger
}

func NewReader(cfg config.ConsumerConfig, processor Processor, logger logger.Logger) *Consumer{
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
			Topic:   cfg.Topic,
			GroupID: "wb-group",
		}),
		Processor: processor,
		Logger: logger,
	}
}

func (c* Consumer) Consume(ctx context.Context) {
	for {
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			c.Logger.Error(ctx, "Fail to fetch message", zap.Error(err))
		}
		if err := c.Processor(ctx, m); err != nil {
			c.Logger.Error(ctx, "Fail to process message", zap.String("Message key", string(m.Key)), zap.Error(err))
		}
		if err := c.Reader.CommitMessages(ctx, m); err != nil {
			c.Logger.Error(ctx, "Fail to commit message", zap.String("Message key", string(m.Key)), zap.Error(err))
		}

		c.Logger.Info(ctx, "Commit message succesfull", zap.String("Mesasge key", string(m.Key)))
	}
}