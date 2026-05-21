package on_closed

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"

	"gorm.io/gorm"
)

type onClosed struct{}

func NewOnClosed() *onClosed {
	return &onClosed{}
}

func (h *onClosed) OnStateTransition(ctx context.Context, tx *gorm.DB, update entitysample.Sample) (entitysample.Sample, error) {
	_ = ctx
	_ = tx
	update.Status = entitysample.SampleStatusClosed
	return update, nil
}
