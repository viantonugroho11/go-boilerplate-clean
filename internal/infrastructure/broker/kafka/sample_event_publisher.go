package kafka

import (
	"context"

	"go-boilerplate-clean/internal/transport/event/events"

	"github.com/viantonugroho11/go-lib/kafka"
)

type SampleEventPublisherKafka struct {
	producer *kafka.Producer[events.SampleEvent]
}

func NewSampleEventPublisherKafka(producer *kafka.Producer[events.SampleEvent]) *SampleEventPublisherKafka {
	return &SampleEventPublisherKafka{producer: producer}
}

func (p *SampleEventPublisherKafka) Publish(ctx context.Context, event events.SampleEvent) error {
	return p.producer.Publish(ctx, event)
}
