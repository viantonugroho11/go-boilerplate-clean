package bootstrap

import (
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/event/events"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"

	"github.com/viantonugroho11/go-lib/kafka"
	"gorm.io/gorm"
)

// WireUserService membuat user repo, Kafka producer/publisher, dan UserService. Pakai Config() global. cleanup menutup producer.
func wireUserService(db *gorm.DB) (usecaseusers.UserService, func(), error) {
	userRepo := userpg.NewUserRepository(db)
	c := Config()

	producer, err := kafka.NewProducer[events.UserCreatedEvent](
		c.KafkaBrokersList(),
		c.Kafka.Topic,
		kafka.WithKeyFunc[events.UserCreatedEvent](func(e events.UserCreatedEvent) []byte { return []byte(e.ID) }),
		kafka.WithIdempotent(),
		kafka.WithRetryMax(5),
	)
	if err != nil {
		return nil, nil, err
	}
	publisher := kafkainfra.NewUserEventPublisherKafka(producer)
	userService := usecaseusers.NewUserServiceWithPublisher(userRepo, publisher)
	return userService, func() { _ = producer.Close() }, nil
}
