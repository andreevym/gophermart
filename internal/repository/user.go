package repository

import (
	"context"
	"time"
)

// User represents a user entity in the application.
type User struct {
	ID       int64      `json:"id"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Created  *time.Time `json:"created_at"`
}

func (u User) IsValidPassword(password string) bool {
	return u.Password == password
}

// UserRepository defines the interface for user repository operations.
//
//go:generate mockgen -source=user.go -destination=./mock/user.go -package=mock
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, userID int64) error
}
