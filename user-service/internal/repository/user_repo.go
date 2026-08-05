package repository

import (
	"context"
	"database/sql"
	"strconv"
	"user-service/internal/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, username, passwordHash, role string) (string, error) {
	var idInt int
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, username, passwordHash, role).Scan(&idInt)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(idInt), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var u entities.User
	var idInt int
	err := r.db.QueryRowContext(ctx, "select id, username, password_hash, role from users where username = $1", username).Scan(&idInt, &u.Username, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}
	u.ID = strconv.Itoa(idInt)
	return &u, nil
}
