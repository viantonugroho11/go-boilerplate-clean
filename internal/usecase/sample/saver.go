package sample

import (
	"context"
	"go-boilerplate-clean/internal/config"
	"go-boilerplate-clean/internal/entity/sample"
	"go-boilerplate-clean/internal/usecase/sample/states"
)

type (

	StateMachineSample interface {
		Do(ctx context.Context, update sample.Sample) (sample.Sample, error)
	}

	NewSampleStateMachine interface {
		NewStateMachine(ctx context.Context, current *sample.Sample) (states.ISampleStateMachine, error)
	}
)

type lenderRepaymentSaver struct {
	stateMachine NewSampleStateMachine
	conf         *config.Configuration
}

func NewLenderRepaymentSaver(
	stateMachine NewSampleStateMachine,
	conf *config.Configuration,
) *lenderRepaymentSaver {
	return &lenderRepaymentSaver{
		stateMachine: stateMachine,
		conf:         conf,
	}
}

func (s *lenderRepaymentSaver) Save(ctx context.Context, sample sample.Sample) (sample.Sample, error) {
	stateMachine, err := s.stateMachine.NewStateMachine(ctx, &sample)
	if err != nil {
		return sample, err
	}
	return stateMachine.Do(ctx, nil, sample)
}