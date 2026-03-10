package event


import (
	"context"
	"go-boilerplate-clean/internal/config"
	"go-boilerplate-clean/internal/transport/event/events"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	users "go-boilerplate-clean/internal/usecase/users"

	kafka "github.com/viantonugroho11/go-lib/kafka"

	kafkahandler "go-boilerplate-clean/internal/transport/event/kafka"

)

// RunUser menjalankan consumer user-created.
func RunUser(ctx context.Context, cfg *config.Configuration, userService users.UserService) (kafka.Consumer, error) {
	return kafkainfra.RunWithConfig[events.UserCreatedEvent](ctx, cfg, cfg.Kafka.GroupIDOrders, cfg.Kafka.TopicOrders, kafkahandler.NewUserCreatedHandler(userService))
}

// RunOrder menjalankan consumer order-created.
func RunOrder(ctx context.Context, cfg *config.Configuration) (kafka.Consumer, error) {
	groupID, topic := cfg.Kafka.GroupIDOrders, cfg.Kafka.TopicOrders
	if groupID == "" {
		groupID = "order-consumer-group"
	}
	if topic == "" {
		topic = "orders"
	}
	return kafkainfra.RunWithConfig[events.OrderCreatedEvent](ctx, cfg, groupID, topic, kafkahandler.NewOrderCreatedHandler())
}