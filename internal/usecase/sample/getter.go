package sample

import (
	"context"
	"strings"

	entitysample "go-boilerplate-clean/internal/entity/sample"
	reposample "go-boilerplate-clean/internal/repository/sample"
)

type sampleGetter struct {
	repo reposample.SampleRepository
}

func NewSampleGetter(repo reposample.SampleRepository) SampleGetter {
	return &sampleGetter{repo: repo}
}

func (g *sampleGetter) Get(ctx context.Context, id string) (*entitysample.Sample, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrSampleIDRequired
	}
	return g.repo.GetByID(ctx, id)
}
