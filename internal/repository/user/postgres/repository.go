package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-boilerplate-clean/internal/usecase"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user usecase.User) (usecase.User, error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	const q = `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, q, user.ID, user.Name, user.Email)
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (usecase.User, error) {
	const q = `SELECT id, name, email FROM users WHERE id = $1`
	var u usecase.User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return usecase.User{}, errors.New("user not found")
	}
	return u, err
}

func (r *UserRepository) List(ctx context.Context) ([]usecase.User, error) {
	const q = `SELECT id, name, email FROM users ORDER BY name`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []usecase.User
	for rows.Next() {
		var u usecase.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, user usecase.User) (usecase.User, error) {
	const q = `UPDATE users SET name=$2, email=$3 WHERE id=$1`
	ct, err := r.pool.Exec(ctx, q, user.ID, user.Name, user.Email)
	if err != nil {
		return usecase.User{}, err
	}
	if ct.RowsAffected() == 0 {
		return usecase.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM users WHERE id=$1`
	ct, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}


