package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-boilerplate-clean/internal/config"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	kafkarunner "go-boilerplate-clean/internal/transport/event/kafka"

	"github.com/IBM/sarama"
	confLoader "github.com/viantonugroho11/go-config-library"
)

func main() {
	// Load configuration (Consul/env/file)
	cfg := config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Root context with cancellation on signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Kafka Consumer
	consumerHandler := func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		return kafkarunner.ExampleHandler(ctx, msg.Key, msg.Value)
	}
	consumer, err := kafkainfra.NewConsumer(
		cfg.KafkaBrokersList(),
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		consumerHandler,
	)
	if err != nil {
		log.Fatalf("kafka consumer init error: %v", err)
	}
	defer func() {
		if cerr := consumer.Close(); cerr != nil {
			log.Printf("kafka consumer close error: %v", cerr)
		}
	}()

	// Start consuming
	kafkarunner.RegisterConsumers(ctx, consumer)
	log.Printf("consumer started, group=%s topic=%s", cfg.Kafka.GroupID, cfg.Kafka.Topic)

	// Wait for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, stopping consumer...")
}
