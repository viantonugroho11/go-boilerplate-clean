package bootstrap

import (
	"go-boilerplate-clean/internal/config"
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

	"github.com/viantonugroho11/go-lib/kafka"
	"gorm.io/gorm"
)

func wireUserService(cfg *config.Configuration, db *gorm.DB) (usecaseusers.UserService, func(), error) {
	userRepo := userpg.NewUserRepository(db)
	producer, err := kafka.NewProducer[events.UserCreatedEvent](
		cfg.KafkaBrokersList(),
		cfg.Kafka.Topic,
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

func wireSampleService(cfg *config.Configuration, db *gorm.DB) (usecasesample.SampleService, func(), error) {
	sampleRepo := samplepg.NewSampleRepository(db)
	txManager := beginpg.NewBeginRepository(db)
	stateFactory := states.NewSampleStateMachineFactory(
		on_open.NewOnOpen(),
		on_pending.NewOnPending(),
		on_closed.NewOnClosed(),
	)
	producer, err := kafka.NewProducer[events.SampleEvent](
		cfg.KafkaBrokersList(),
		cfg.Kafka.Topic,
		kafka.WithKeyFunc[events.SampleEvent](func(e events.SampleEvent) []byte { return []byte(e.ResourceID) }),
		kafka.WithIdempotent(),
		kafka.WithRetryMax(5),
	)
	if err != nil {
		return nil, nil, err
	}
	publisher := kafkainfra.NewSampleEventPublisherKafka(producer)
	svc := usecasesample.NewSampleSaver(
		stateFactory,
		txManager,
		sampleRepo,
		usecasesample.NewSampleGetter(sampleRepo),
		sampleRepo,
		publisher,
	)
	return svc, func() { _ = producer.Close() }, nil
}
