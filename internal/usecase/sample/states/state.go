package states

import (
	"context"
	"fmt"

	"go-boilerplate-clean/internal/entity/sample"

	"gorm.io/gorm"
)

type ISampleState interface {
	Do(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error)
}
type ISampleStateMachine interface {
	ISampleState
	Sample() *sample.Sample
}
type ISampleNewStateMachine interface {
	NewSampleStateMachine(ctx context.Context, tx *gorm.DB, current *sample.Sample) (ISampleStateMachine, error)
}

type IOnStateTransition interface {
	OnStateTransition(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error)
}

type stateMachineSample struct {
	data        *sample.Sample
	current     ISampleState
	open        ISampleState
	onHold      ISampleState
	closed      ISampleState
}

type stateMachineFactorySample struct {
	onCreation IOnStateTransition
	onHold     IOnStateTransition
	onClose    IOnStateTransition

	onBuybackEligible IOnStateTransition
	onBuybackSuccess  IOnStateTransition
	onCancelled       IOnStateTransition
}

func NewSampleStateMachineFactory(
	onCreation IOnStateTransition,
	onHold IOnStateTransition,
	onClose IOnStateTransition,

	onBuybackEligible IOnStateTransition,
	onBuybackSuccess IOnStateTransition,
	onCancelled IOnStateTransition,
) *stateMachineFactorySample {
	return &stateMachineFactorySample{
		onCreation:        onCreation,
		onHold:            onHold,
		onClose:           onClose,
		onBuybackEligible: onBuybackEligible,
		onBuybackSuccess:  onBuybackSuccess,
		onCancelled:       onCancelled,
	}
}

func (smf stateMachineFactorySample) NewStateMachine(ctx context.Context, current *sample.Sample) (ISampleStateMachine, error) {
	sm := &stateMachineSample{}

	sm.open = open{
		stateMachine: sm,
		onCreation:   smf.onCreation,
		onHold:       smf.onHold,
		onClose:      smf.onClose,
		onCancelled:  smf.onCancelled,
	}

	sm.data = current
	if sm.data.ID == "" {
		
		return nil, fmt.Errorf("sample ID is required")
	}


	switch sm.data.Status {
	case sample.SampleStatusOpen:
		sm.current = sm.open
	case sample.SampleStatusOnHold:
		sm.current = sm.onHold
	case sample.SampleStatusClosed:
		sm.current = sm.closed
	default:
		return nil, fmt.Errorf("unknown status: %s", sm.data.Status)
	}

	return sm, nil
}
func (s stateMachineSample) Do(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error) {
	return s.current.Do(ctx, tx, update)
}

func (s stateMachineSample) Sample() *sample.Sample {
	return s.data
}