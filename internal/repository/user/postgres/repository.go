package postgres

import (
	"context"
	"errors"

	userEntity "go-boilerplate-clean/internal/entity/users"
	"go-boilerplate-clean/internal/repository/user"
	"go-boilerplate-clean/internal/repository/user/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	var m model.User
	m = model.ToModel(user)
	err := tx.WithContext(ctx).Create(&m).Error
	if err != nil {
		return userEntity.User{}, err
	}
	return model.ToEntity(&m), nil

}

func (r *userRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (userEntity.User, error) {
	var u model.User
	if tx != nil {
		err := tx.WithContext(ctx).First(&u, "id = ?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEntity.User{}, errors.New("user not found")
		}
		return model.ToEntity(&u), nil
	}
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return userEntity.User{}, err
	}
	return model.ToEntity(&u), nil
}

func (r *userRepository) List(ctx context.Context, tx *gorm.DB) ([]userEntity.User, error) {
	var result []userEntity.User
	var rows []model.User
	if tx != nil {
		err := tx.WithContext(ctx).Order("name").Find(&rows).Error
		if err != nil {
			return nil, err
		}
		for _, u := range rows {
			result = append(result, model.ToEntity(&u))
		}
		return result, nil
	}
	err := r.db.WithContext(ctx).Order("name").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, u := range rows {
		result = append(result, model.ToEntity(&u))
	}
	return result, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error) {
	if tx == nil {
		return userEntity.User{}, errors.New("not implemented")
	}
	var u model.User
	u = model.ToModel(user)
	err := tx.WithContext(ctx).Model(&u).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	}).Error
	if err != nil {
		return userEntity.User{}, err
	}
	if tx.RowsAffected == 0 {
		return userEntity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	err := tx.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}
	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
