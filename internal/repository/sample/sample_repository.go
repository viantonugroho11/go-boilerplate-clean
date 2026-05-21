package sample

import (
	"context"

	entitysample "go-boilerplate-clean/internal/entity/sample"

	"gorm.io/gorm"
)

type SampleRepository interface {
	Add(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
	GetByID(ctx context.Context, id string) (*entitysample.Sample, error)
	Update(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error)
}
