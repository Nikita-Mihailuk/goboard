package kafka

import (
	"encoding/json"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flashTimeout = 5000 // ms
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address string) (*Producer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers": address,
	}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
	}, nil
}

func (p *Producer) Produce(message dto.UpdateUserMessage, topic string) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msgBytes,
		Key:   nil,
	}

	kafkaChan := make(chan kafka.Event)
	if err = p.producer.Produce(&kafkaMsg, kafkaChan); err != nil {
		return err
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errors.New("unknown message type")
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flashTimeout)
	p.producer.Close()
}
