package user

import (
	"context"
	userEntity "go-boilerplate-clean/internal/entity/users"

	"gorm.io/gorm"
)

// UserRepository is the persistence contract for the User entity.
// Implementations (Postgres, etc.) must satisfy this interface.

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error)
	GetByID(ctx context.Context, tx *gorm.DB, id string) (userEntity.User, error)
	List(ctx context.Context, tx *gorm.DB, offset, limit int) ([]userEntity.User, int64, error)
	Update(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
}
