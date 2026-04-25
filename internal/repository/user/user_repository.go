package user

import (
	"context"
	userEntity "go-boilerplate-clean/internal/entity/users"

	"gorm.io/gorm"
)

// Interface repository untuk entity User.
// Implementasi (Postgres/Mongo/dll) harus memenuhi kontrak ini.
// Menggunakan model dari usecase untuk penyederhanaan.

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error)
	GetByID(ctx context.Context,tx *gorm.DB, id string) (userEntity.User, error)
	List(ctx context.Context,tx *gorm.DB) ([]userEntity.User, error)
	Update(ctx context.Context,tx *gorm.DB, user userEntity.User) (userEntity.User, error)
	Delete(ctx context.Context,tx *gorm.DB, id string) error
}
