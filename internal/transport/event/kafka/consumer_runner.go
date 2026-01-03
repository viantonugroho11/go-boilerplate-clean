package kafka

import (
	"context"
	"log"

	infra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
)

// RegisterConsumers menyiapkan dan menjalankan consumer Kafka dengan handler contoh.
// Anda bisa mengganti handler untuk memanggil usecase tertentu.
func RegisterConsumers(ctx context.Context, consumer *infra.Consumer) {
	consumer.Start(ctx)
}

// ExampleHandler adalah contoh handler pesan Kafka.
// Di sini hanya logging; ganti dengan pemanggilan usecase sebenarnya.
func ExampleHandler(ctx context.Context, msgKey, msgVal []byte) error {
	log.Printf("kafka message: key=%s val=%s", string(msgKey), string(msgVal))
	return nil
}


