package broker

import "context"

// Producer mendefinisikan kontrak untuk message broker producer.
type Producer interface {
	Publish(ctx context.Context, topic string, key, value []byte) (partition int32, offset int64, err error)
	Close() error
}

// Consumer mendefinisikan kontrak untuk message broker consumer.
type Consumer interface {
	Start(ctx context.Context)
	Close() error
}
