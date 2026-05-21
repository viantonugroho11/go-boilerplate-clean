package on_pending

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"

	"gorm.io/gorm"
)

type onPending struct{}

func NewOnPending() *onPending {
	return &onPending{}
}

func (h *onPending) OnStateTransition(ctx context.Context, tx *gorm.DB, update entitysample.Sample) (entitysample.Sample, error) {
	_ = ctx
	_ = tx
	update.Status = entitysample.SampleStatusOnHold
	return update, nil
}
