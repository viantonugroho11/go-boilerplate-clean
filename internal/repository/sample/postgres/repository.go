package postgres

import (
	"context"
	"errors"

	entitysample "go-boilerplate-clean/internal/entity/sample"
	reposample "go-boilerplate-clean/internal/repository/sample"
	"go-boilerplate-clean/internal/repository/sample/model"
	"go-boilerplate-clean/internal/shared/apperrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type sampleRepository struct {
	db *gorm.DB
}

func NewSampleRepository(db *gorm.DB) reposample.SampleRepository {
	return &sampleRepository{db: db}
}

func (r *sampleRepository) Add(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error) {
	if tx == nil {
		return entitysample.Sample{}, errors.New("transaction is required")
	}
	if sample.ID == "" {
		sample.ID = uuid.NewString()
	}
	if sample.Status == "" {
		sample.Status = entitysample.SampleStatusOpen
	}
	m := model.ToModel(sample)
	if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
		return entitysample.Sample{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *sampleRepository) GetByID(ctx context.Context, id string) (*entitysample.Sample, error) {
	var m model.Sample
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrSampleNotFound
	}
	if err != nil {
		return nil, err
	}
	s := model.ToEntity(&m)
	return &s, nil
}

func (r *sampleRepository) Update(ctx context.Context, tx *gorm.DB, sample entitysample.Sample) (entitysample.Sample, error) {
	if tx == nil {
		return entitysample.Sample{}, errors.New("transaction is required")
	}
	err := tx.WithContext(ctx).Model(&model.Sample{}).Where("id = ?", sample.ID).Updates(map[string]interface{}{
		"code":   sample.Code,
		"name":   sample.Name,
		"email":  sample.Email,
		"status": sample.Status,
	}).Error
	if err != nil {
		return entitysample.Sample{}, err
	}
	var m model.Sample
	if err := tx.WithContext(ctx).First(&m, "id = ?", sample.ID).Error; err != nil {
		return entitysample.Sample{}, err
	}
	return model.ToEntity(&m), nil
}
