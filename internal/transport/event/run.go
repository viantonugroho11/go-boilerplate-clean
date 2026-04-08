package event

import (
	"context"

	"go-boilerplate-clean/internal/config"
	infrakafka "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	transportkafka "go-boilerplate-clean/internal/transport/event/kafka"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"
)

const (
	ConsumerNameUser  = "user"
	ConsumerNameOrder = "order"
)

// ConsumerNames daftar nama consumer yang didukung (flag -consumer).
func ConsumerNames() []string {
	return []string{ConsumerNameUser, ConsumerNameOrder}
}

// RunUser menjalankan consumer Kafka untuk event user (topic & group dari cfg.Kafka).
func RunUser(ctx context.Context, cfg *config.Configuration, userService usecaseusers.UserService) (interface{ Close() error }, error) {
	h := transportkafka.NewUserCreatedHandler(userService)
	return infrakafka.RunWithConfig(ctx, cfg, cfg.Kafka.GroupID, cfg.Kafka.Topic, h)
}

// RunOrder menjalankan consumer Kafka untuk event order (topic_orders & group_id_orders).
func RunOrder(ctx context.Context, cfg *config.Configuration) (interface{ Close() error }, error) {
	h := transportkafka.NewOrderCreatedHandler()
	return infrakafka.RunWithConfig(ctx, cfg, cfg.Kafka.GroupIDOrders, cfg.Kafka.TopicOrders, h)
}
