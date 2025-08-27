package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xanuthatusu/tepia/internal/models"
)

func CreateUser(ctx context.Context, pool *pgxpool.Pool, name, email string) (string, error) {
	uuid := uuid.New()
	if _, err := pool.Exec(ctx,
		"INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
		uuid, name, email); err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func ListUsers(ctx context.Context, pool *pgxpool.Pool) ([]models.User, error) {
	rows, err := pool.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
