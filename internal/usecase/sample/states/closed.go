package states

import (
	"context"
	"go-boilerplate-clean/internal/entity/sample"
	"gorm.io/gorm"
)

type closed struct {
	stateMachine *stateMachineSample
	onPending    IOnStateTransition
	onClosed     IOnStateTransition
}

func (s closed) Do(ctx context.Context, tx *gorm.DB, update sample.Sample) (sample.Sample, error) {
	s.stateMachine.data = &update

	return s.onClosed.OnStateTransition(ctx, tx, update)
}