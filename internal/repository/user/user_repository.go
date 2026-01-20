package user

import (
	"context"
	"go-boilerplate-clean/internal/entity"
)

// Interface repository untuk entity User.
// Implementasi (Postgres/Mongo/dll) harus memenuhi kontrak ini.
// Menggunakan model dari usecase untuk penyederhanaan.

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	List(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
}
