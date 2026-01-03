package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	sp sarama.SyncProducer
}

func NewProducer(brokers []string, clientID string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Idempotent = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Retry.Backoff = 200 * time.Millisecond
	prod, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{sp: prod}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	return p.sp.SendMessage(msg)
}

func (p *Producer) Close() error {
	return p.sp.Close()
}


