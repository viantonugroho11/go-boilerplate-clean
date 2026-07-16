package postgres

import (
	"context"
	"errors"

	userEntity "go-boilerplate-clean/internal/entity/users"
	"go-boilerplate-clean/internal/repository/user"
	"go-boilerplate-clean/internal/repository/user/model"
	"go-boilerplate-clean/internal/shared/apperrors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) conn(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, u userEntity.User) (userEntity.User, error) {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	m := model.ToModel(u)
	if err := r.conn(tx).WithContext(ctx).Create(&m).Error; err != nil {
		return userEntity.User{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *userRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (userEntity.User, error) {
	var m model.User
	err := r.conn(tx).WithContext(ctx).First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userEntity.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return userEntity.User{}, err
	}
	return model.ToEntity(&m), nil
}

func (r *userRepository) List(ctx context.Context, tx *gorm.DB, offset, limit int) ([]userEntity.User, int64, error) {
	db := r.conn(tx).WithContext(ctx)
	var total int64
	if err := db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.User
	if err := db.Order("name").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	result := make([]userEntity.User, 0, len(rows))
	for i := range rows {
		result = append(result, model.ToEntity(&rows[i]))
	}
	return result, total, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, u userEntity.User) (userEntity.User, error) {
	m := model.ToModel(u)
	res := r.conn(tx).WithContext(ctx).Model(&m).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"name":  u.Name,
		"email": u.Email,
	})
	if res.Error != nil {
		return userEntity.User{}, res.Error
	}
	if res.RowsAffected == 0 {
		return userEntity.User{}, apperrors.ErrUserNotFound
	}
	return u, nil
}

func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	res := r.conn(tx).WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}
