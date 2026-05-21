package users

import (
	"context"
	"log"
	"strings"

	userEntity "go-boilerplate-clean/internal/entity/users"
	begin "go-boilerplate-clean/internal/repository/begin"
	repouser "go-boilerplate-clean/internal/repository/user"
	"go-boilerplate-clean/internal/shared/apperrors"
)

type UserService interface {
	Create(ctx context.Context, user userEntity.User) (userEntity.User, error)
	GetByID(ctx context.Context, id string) (userEntity.User, error)
	List(ctx context.Context) ([]userEntity.User, error)
	Update(ctx context.Context, user userEntity.User) (userEntity.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo      repouser.UserRepository
	txManager begin.BeginRepository
	publisher UserEventPublisher 
}

func NewUserService(repo repouser.UserRepository, txManager begin.BeginRepository, publisher UserEventPublisher) UserService {
	return &userService{repo: repo, txManager: txManager, publisher: publisher}
}

func (s *userService) Create(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	if err := validateUser(user, true); err != nil {
		return userEntity.User{}, err
	}
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return userEntity.User{}, err
	}

	defer func() {
		if err != nil {
			s.txManager.Rollback(ctx, tx)
		}
	}()
	created, err := s.repo.Create(ctx, tx, user)
	if err != nil {
		return userEntity.User{}, err
	}
	err = s.txManager.Commit(ctx, tx)
	if err != nil {
		return userEntity.User{}, err
	}

	err = s.publisher.PublishUser(ctx, created)
	if err != nil {
		log.Printf("user_usecase: PublishUserCreated: %v", err)
	}
	return created, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (userEntity.User, error) {
	if strings.TrimSpace(id) == "" {
		return userEntity.User{}, apperrors.ErrUserIDRequired
	}
	return s.repo.GetByID(ctx, nil, id)
}

func (s *userService) List(ctx context.Context) ([]userEntity.User, error) {
	return s.repo.List(ctx, nil)
}

func (s *userService) Update(ctx context.Context, user userEntity.User) (userEntity.User, error) {
	if strings.TrimSpace(user.ID) == "" {
		return userEntity.User{}, apperrors.ErrUserIDRequired
	}
	if err := validateUser(user, false); err != nil {
		return userEntity.User{}, err
	}
	updated, err := s.repo.Update(ctx, nil, user)
	if err != nil {
		return userEntity.User{}, err
	}

	err = s.publisher.PublishUser(ctx, updated)
	if err != nil {
		log.Printf("user_usecase: PublishUserUpdated: %v", err)
	}
	return updated, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return apperrors.ErrUserIDRequired
	}
	return s.repo.Delete(ctx, nil, id)
}

func validateUser(user userEntity.User, creating bool) error {
	if creating && strings.TrimSpace(user.ID) != "" {
		// ID akan diisi oleh repository saat create jika kosong
	}
	if strings.TrimSpace(user.Name) == "" {
		return apperrors.ErrUserNameRequired
	}
	if strings.TrimSpace(user.Email) == "" {
		return apperrors.ErrUserEmailRequired
	}
	return nil
}
