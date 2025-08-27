package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xanuthatusu/tepia/internal/models"
)

func CreateUserWithPassword(ctx context.Context, pool *pgxpool.Pool, name, email, hash string) error {
	if _, err := pool.Exec(ctx,
		"INSERT INTO users (id, name, email, pass_hash) VALUES ($1, $2, $3, $4)",
		uuid.New(), name, email, hash); err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (models.User, error) {
	row := pool.QueryRow(ctx, "SELECT id, name, email, pass_hash FROM users WHERE email = $1", email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return models.User{}, err
	}

	return user, nil
}
