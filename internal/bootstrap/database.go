package bootstrap

import (
	"context"

	pginfra "go-boilerplate-clean/internal/infrastructure/database/postgres"
	"gorm.io/gorm"
)

// InitDB connect ke Postgres dan jalankan migrate. Pakai Config() global.
func initDB() (*gorm.DB, error) {
	ctx := context.Background()
	db, err := pginfra.Connect(ctx, Config().PGDSN())
	if err != nil {
		return nil, err
	}
	if err := pginfra.Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}
