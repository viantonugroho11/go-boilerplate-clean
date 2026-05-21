package sample

import (
	"context"
	"strings"

	entitysample "go-boilerplate-clean/internal/entity/sample"
	reposample "go-boilerplate-clean/internal/repository/sample"
	"go-boilerplate-clean/internal/shared/apperrors"

	"gorm.io/gorm"
)

type sampleGetter struct {
	repo reposample.SampleRepository
}

func NewSampleGetter(repo reposample.SampleRepository) SampleGetter {
	return &sampleGetter{repo: repo}
}

func (g *sampleGetter) Get(ctx context.Context, tx *gorm.DB, id string) (*entitysample.Sample, error) {
	if strings.TrimSpace(id) == "" {
		return nil, apperrors.ErrSampleIDRequired
	}
	return g.repo.GetByID(ctx, tx, id)
}
