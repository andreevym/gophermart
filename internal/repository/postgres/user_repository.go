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
	// ErrUserNotFound представляет ошибку, возникающую при отсутствии пользователя в базе данных.
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository представляет репозиторий пользователей с использованием PostgreSQL.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных PostgreSQL.
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

// GetUserByID возвращает пользователя из базы данных PostgreSQL по его идентификатору.
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

// GetUserByUsername возвращает пользователя из базы данных PostgreSQL по его имени пользователя.
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	sql := `SELECT id, username, password FROM users WHERE username = $1`
	var user repository.User
	err := r.db.QueryRow(ctx, sql, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// UpdateUser обновляет информацию о пользователе в базе данных PostgreSQL.
func (r *UserRepository) UpdateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	sql := `UPDATE users SET username = $1, password = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, sql, user.Username, user.Password, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return user, nil
}

// DeleteUser удаляет пользователя из базы данных PostgreSQL по его идентификатору.
func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	sql := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
