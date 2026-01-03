package user

import (
)

// Interface repository untuk entity User.
// Implementasi (Postgres/Mongo/dll) harus memenuhi kontrak ini.
// Menggunakan model dari usecase untuk penyederhanaan.
import (
	"context"
	"go-boilerplate-clean/internal/usecase"
)

type Repository interface {
	Create(ctx context.Context, user usecase.User) (usecase.User, error)
	GetByID(ctx context.Context, id string) (usecase.User, error)
	List(ctx context.Context) ([]usecase.User, error)
	Update(ctx context.Context, user usecase.User) (usecase.User, error)
	Delete(ctx context.Context, id string) error
}
