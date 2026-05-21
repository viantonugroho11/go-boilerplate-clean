package bootstrap

import (
	entitysample "go-boilerplate-clean/internal/entity/sample"
	kafkainfra "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	beginpg "go-boilerplate-clean/internal/repository/begin/postgres"
	samplepg "go-boilerplate-clean/internal/repository/sample/postgres"
	userpg "go-boilerplate-clean/internal/repository/user/postgres"
	"go-boilerplate-clean/internal/transport/event/events"
	usecasesample "go-boilerplate-clean/internal/usecase/sample"
	"go-boilerplate-clean/internal/usecase/sample/on_closed"
	"go-boilerplate-clean/internal/usecase/sample/on_open"
	"go-boilerplate-clean/internal/usecase/sample/on_pending"
	"go-boilerplate-clean/internal/usecase/sample/states"
	usecaseusers "go-boilerplate-clean/internal/usecase/users"
	"log"

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
	txManager := beginpg.NewBeginRepository(db)
	userService := usecaseusers.NewUserService(userRepo, txManager, publisher)
	return userService, func() { _ = producer.Close() }, nil
}

func wireSampleService(db *gorm.DB) usecasesample.SampleService {
	sampleRepo := samplepg.NewSampleRepository(db)
	txManager := beginpg.NewBeginRepository(db)
	stateFactory := states.NewSampleStateMachineFactory(
		on_open.NewOnOpen(),
		on_pending.NewOnPending(),
		on_closed.NewOnClosed(),
	)
	producer, err := kafka.NewProducer[entitysample.Sample](
		Config().KafkaBrokersList(),
		Config().Kafka.Topic,
		kafka.WithKeyFunc[entitysample.Sample](func(e entitysample.Sample) []byte { return []byte(e.ID) }),
		kafka.WithIdempotent(),
		kafka.WithRetryMax(5),
	)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	publisher := kafkainfra.NewSampleEventPublisherKafka(producer)
	return usecasesample.NewSampleSaver(
		stateFactory,
		txManager,
		sampleRepo,
		usecasesample.NewSampleGetter(sampleRepo),
		sampleRepo,
		publisher,
	)
}
