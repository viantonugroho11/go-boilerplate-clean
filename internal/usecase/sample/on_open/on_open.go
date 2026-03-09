package on_open

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"
	"gorm.io/gorm"
)

type onOpen struct{}

func NewOnOpen() *onOpen {
	return &onOpen{}
}

func (s *onOpen) OnStateTransition(ctx context.Context, tx *gorm.DB, update entitysample.Sample) (entitysample.Sample, error) {
	return update, nil
}