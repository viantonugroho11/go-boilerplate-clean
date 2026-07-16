package bootstrap

import (
	"context"

	"go-boilerplate-clean/internal/config"
	pginfra "go-boilerplate-clean/internal/infrastructure/database/postgres"
	"gorm.io/gorm"
)

func initDB(cfg *config.Configuration) (*gorm.DB, error) {
	ctx := context.Background()
	db, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		return nil, err
	}
	if err := pginfra.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}
