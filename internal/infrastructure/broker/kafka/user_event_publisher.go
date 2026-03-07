package kafka

import (
	"context"

	"go-boilerplate-clean/internal/transport/event/events"
	"go-boilerplate-clean/internal/usecase/users"

	"github.com/viantonugroho11/go-lib/kafka"
)

// UserEventPublisherKafka implementasi users.UserEventPublisher (go-lib Producer).
type UserEventPublisherKafka struct {
	producer *kafka.Producer[events.UserCreatedEvent]
}

func NewUserEventPublisherKafka(producer *kafka.Producer[events.UserCreatedEvent]) users.UserEventPublisher {
	return &UserEventPublisherKafka{producer: producer}
}

func (p *UserEventPublisherKafka) PublishUserCreated(ctx context.Context, id, name, email string) error {
	return p.producer.Publish(ctx, events.UserCreatedEvent{
		ID:    id,
		Name:  name,
		Email: email,
	})
}
