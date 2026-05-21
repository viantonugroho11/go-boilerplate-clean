package kafka

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"

	"github.com/viantonugroho11/go-lib/kafka"
)

type SampleEventPublisherKafka struct {
	producer *kafka.Producer[entitysample.Sample]
}

func NewSampleEventPublisherKafka(producer *kafka.Producer[entitysample.Sample]) *SampleEventPublisherKafka {
	return &SampleEventPublisherKafka{producer: producer}
}

func (p *SampleEventPublisherKafka) Publish(ctx context.Context, sample entitysample.Sample) error {
	return p.producer.Publish(ctx, entitysample.Sample{
		ID:    sample.ID,
		Code:  sample.Code,
		Name:  sample.Name,
		Email: sample.Email,
		Status: sample.Status,
	})
}