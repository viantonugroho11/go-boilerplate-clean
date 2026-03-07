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

func (r *userRepository) Create(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	m := model.User{ID: user.ID, Name: user.Name, Email: user.Email}
	err := r.db.WithContext(ctx).Create(&m).Error
	return userEntity.User{ID: m.ID, Name: m.Name, Email: m.Email}, err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (userEntity.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userEntity.User{}, errors.New("user not found")
	}
	return userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email}, err
}

func (r *userRepository) List(ctx context.Context) ([]userEntity.User, error) {
	var result []userEntity.User
	var rows []model.User
	if err := r.db.WithContext(ctx).Order("name").Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, u := range rows {
		result = append(result, userEntity.User{ID: u.ID, Name: u.Name, Email: u.Email})
	}
	return result, nil
}

func (r *userRepository) Update(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	tx := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	})
	if tx.Error != nil {
		return userEntity.User{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return userEntity.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	tx := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
