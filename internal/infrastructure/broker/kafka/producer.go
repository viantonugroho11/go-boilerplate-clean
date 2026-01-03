package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	sp sarama.SyncProducer
}

// ProducerOption memungkinkan kustomisasi konfigurasi producer sebelum dibuat.
type ProducerOption func(cfg *sarama.Config)

func NewProducer(brokers []string, clientID string, opts ...ProducerOption) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Idempotent = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Retry.Backoff = 200 * time.Millisecond

	// apply options
	for _, opt := range opts {
		if opt != nil {
			opt(cfg)
		}
	}

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

// options config
// with max retry
func WithMaxRetry(maxRetry int) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Retry.Max = maxRetry
	}
}

// with retry backoff
func WithRetryBackoff(retryBackoff time.Duration) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Retry.Backoff = retryBackoff
	}
}


// with idempotent
func WithIdempotent(idempotent bool) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Idempotent = idempotent
	}
}

// with required acks
func WithRequiredAcks(requiredAcks sarama.RequiredAcks) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.RequiredAcks = requiredAcks
	}
}

// with return successes
func WithReturnSuccesses(returnSuccesses bool) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Return.Successes = returnSuccesses
	}
}

// with return errors
func WithReturnErrors(returnErrors bool) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Return.Errors = returnErrors
	}
}

// with compression codec
func WithCompressionCodec(compressionCodec sarama.CompressionCodec) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.Compression = compressionCodec
	}
}

// with compression level
func WithCompressionLevel(compressionLevel int) ProducerOption {
	return func(cfg *sarama.Config) {
		cfg.Producer.CompressionLevel = compressionLevel
	}
}