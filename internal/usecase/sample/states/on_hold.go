package states

import (
	"context"
	"go-boilerplate-clean/internal/entity/sample"
	"gorm.io/gorm"
)

type onHold struct {
	stateMachine *stateMachineSample
	onPending    IOnStateTransition
	onClosed     IOnStateTransition
}

func (s onHold) Do(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error) {
	s.stateMachine.data = &update

	switch update.Status {
	case sample.SampleStatusClosed:
		return s.onClosed.OnStateTransition(ctx, tx, update)
	default:
		return s.onPending.OnStateTransition(ctx, tx, update)
	}
	return s.onPending.OnStateTransition(ctx, tx, update)
}