package postgres

import (
	"context"

	"go-boilerplate-clean/internal/repository/begin"

	"gorm.io/gorm"
)

type beginRepository struct {
	db *gorm.DB
}

func NewBeginRepository(db *gorm.DB) begin.BeginRepository {
	return &beginRepository{db: db}
}

func (r *beginRepository) Begin(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *beginRepository) Commit(ctx context.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Commit().Error
}

func (r *beginRepository) Rollback(ctx context.Context, tx *gorm.DB) error {
	return tx.WithContext(ctx).Rollback().Error
}