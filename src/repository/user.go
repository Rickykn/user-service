package repository

import (
	"context"
	"database/sql"
	"github.com/Rickykn/user-service/src/model"
)

type IUserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	query := `
        SELECT id, username, email
        FROM users
        WHERE username = $1
        LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4,$5)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}
