package states

import (
	"context"

	"go-boilerplate-clean/internal/entity/sample"
	"gorm.io/gorm"
)

type open struct {
	stateMachine *stateMachineSample

	onCreation IOnStateTransition
	onHold     IOnStateTransition
	onClose    IOnStateTransition
	onCancelled IOnStateTransition
}

func (s open) Do(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case sample.SampleStatusOnHold:
		return s.onHold.OnStateTransition(ctx, tx, update)
	case sample.SampleStatusClosed:
		return s.onClose.OnStateTransition(ctx, tx, update)
	default:
		return s.onCreation.OnStateTransition(ctx, tx, update)
	}
}
