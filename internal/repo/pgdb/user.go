package pgdb

import (
	"context"
	"fmt"

	"github.com/azoma13/marketplace-service/internal/entity"
	"github.com/azoma13/marketplace-service/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (int, error) {
	query := `
		INSERT INTO users
			(username, password)
			VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := r.Pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %w", err)
	}

	return id, nil
}

func (r *UserRepo) GetUsernameByUserID(ctx context.Context, userId int) (string, error) {
	query := `
			SELECT username
				FROM users
			WHERE id = $1
		`
	var username string
	err := r.Pool.QueryRow(ctx, query, userId).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("UserRepo.GetUsernameByUserID - r.Pool.QueryRow: %v", err)
	}

	return username, nil
}

func (r *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (entity.User, error) {
	query := `
			SELECT id, username, password, created_at
				FROM users
			WHERE username = $1 AND password = $2
		`
	var user entity.User
	err := r.Pool.QueryRow(ctx, query, username, password).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}
