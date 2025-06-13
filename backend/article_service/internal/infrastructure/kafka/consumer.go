package kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

const (
	sessionTimeout     = 7000 // ms
	consumeTimeout     = -1
	autoCommitInterval = 5000 // ms
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg []byte) error
}

type Consumer struct {
	consumer       *kafka.Consumer
	messageHandler MessageHandler
	log            *zap.Logger
	stop           bool
}

func NewConsumer(messageHandler MessageHandler, address, topic, consumerGroup string, log *zap.Logger) (*Consumer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  autoCommitInterval,
		"auto.offset.reset":        "earliest",
	}
	c, err := kafka.NewConsumer(&cfg)
	if err != nil {
		return nil, err
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:       c,
		messageHandler: messageHandler,
		log:            log,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			return
		}
		kafkaMsg, err := c.consumer.ReadMessage(consumeTimeout)
		if err != nil {
			c.log.Error("error read message:", zap.Error(err))
			continue
		}

		if kafkaMsg == nil {
			continue
		}

		if err = c.messageHandler.HandleMessage(context.TODO(), kafkaMsg.Value); err != nil { // TODO: change context
			c.log.Error("error handle message:", zap.Error(err))
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
