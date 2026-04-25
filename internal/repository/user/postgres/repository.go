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
	m := model.User{ID: user.ID, Name: user.Name, Email: user.Email}
	if tx != nil {
		err := tx.WithContext(ctx).Create(&m).Error
		if err != nil {
			return userEntity.User{}, err
		}
		return userEntity.User{ID: m.ID, Name: m.Name, Email: m.Email}, nil
	}
	err := r.db.WithContext(ctx).Create(&m).Error
	if err != nil {
		return userEntity.User{}, err
	}
	return userEntity.User{ID: m.ID, Name: m.Name, Email: m.Email}, err
}

func (r *userRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (userEntity.User, error) {
	var u model.User
	if tx != nil {
		err := tx.WithContext(ctx).First(&u, "id = ?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userEntity.User{}, errors.New("user not found")
		}
		return userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email}, nil
	}
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return userEntity.User{}, err
	}
	return userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email}, nil
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
			result = append(result, userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email})
		}
		return result, nil
	}
	err := r.db.WithContext(ctx).Order("name").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, u := range rows {
		result = append(result, userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email})
	}
	return result, nil
}

func (r *userRepository) Update(ctx context.Context, tx *gorm.DB, user userEntity.User) (userEntity.User, error) {
	if tx != nil {
		err := tx.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
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
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
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
	if tx != nil {
		err := tx.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
		if err != nil {
			return err
		}
		if tx.RowsAffected == 0 {
			return errors.New("user not found")
		}
		return nil
	}
	err := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
