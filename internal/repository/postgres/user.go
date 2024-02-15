package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	sql := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	var userID int64
	err := r.db.QueryRow(ctx, sql, user.Username, user.Password).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	user.ID = userID
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (*repository.User, error) {
	sql := `SELECT id, username, password FROM users WHERE id = $1`
	var user repository.User
	err := r.db.QueryRow(ctx, sql, userID).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	sql := `SELECT id, username, password FROM users WHERE username = $1`
	var user repository.User
	err := r.db.QueryRow(ctx, sql, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username %s: %v", username, err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	sql := `UPDATE users SET username = $1, password = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, sql, user.Username, user.Password, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	sql := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user by userID %d: %v", userID, err)
	}

	return nil
}
