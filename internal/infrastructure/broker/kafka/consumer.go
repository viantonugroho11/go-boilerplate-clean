package kafka

import (
	"context"
	"log"

	"go-boilerplate-clean/internal/config"

	"github.com/IBM/sarama"
	"github.com/viantonugroho11/go-lib/kafka"
)

// Run menjalankan consumer untuk event tipe E. Pakai RunWithConfig bila punya *config.Configuration.
func Run[E any](
	ctx context.Context,
	brokers []string,
	groupID, topic string,
	handler kafka.EventHandler[E],
	opts ...kafka.ConsumerOption,
) (kafka.Consumer, error) {
	c, err := kafka.NewConsumer[E](brokers, groupID, topic, handler, opts...)
	if err != nil {
		return nil, err
	}
	c.Start(ctx)
	log.Printf("kafka consumer started: group=%s topic=%s", groupID, topic)
	return c, nil
}

// RunWithConfig seperti Run tapi ambil brokers & opsi dari cfg. Supaya transport tidak ulang brokers/opts.
func RunWithConfig[E any](ctx context.Context, cfg *config.Configuration, groupID, topic string, handler kafka.EventHandler[E]) (kafka.Consumer, error) {
	return Run[E](ctx, cfg.KafkaBrokersList(), groupID, topic, handler, DefaultConsumerOptions()...)
}

// DefaultConsumerOptions opsi default (offset oldest).
func DefaultConsumerOptions() []kafka.ConsumerOption {
	return []kafka.ConsumerOption{kafka.WithInitialOffset(sarama.OffsetOldest)}
}
