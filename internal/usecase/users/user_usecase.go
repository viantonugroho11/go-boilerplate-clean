package users

import (
	"context"
	"errors"
	"strings"

	"go-boilerplate-clean/internal/entity"
	repouser "go-boilerplate-clean/internal/repository/user"
)

type UserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	List(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repouser.UserRepository
}

func NewUserService(repo repouser.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	if err := validateUser(user, true); err != nil {
		return entity.User{}, err
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (entity.User, error) {
	if strings.TrimSpace(id) == "" {
		return entity.User{}, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]entity.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) Update(ctx context.Context, user entity.User) (entity.User, error) {
	if strings.TrimSpace(user.ID) == "" {
		return entity.User{}, errors.New("id is required")
	}
	if err := validateUser(user, false); err != nil {
		return entity.User{}, err
	}
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}

func validateUser(user entity.User, creating bool) error {
	if creating && strings.TrimSpace(user.ID) != "" {
		// ID akan diisi oleh repository saat create jika kosong
	}
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	return nil
}
