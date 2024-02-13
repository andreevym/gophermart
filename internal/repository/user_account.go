package repository

import (
	"context"
	"math/big"
)

type UserAccount struct {
	UserID  int64    `json:"userId"`
	Balance *big.Int `json:"balance"`
}

// UserAccountRepository defines the interface for user account repository operations.
//
//go:generate mockgen -source=user_account.go -destination=./mock/user_account_mock.go -package=mock
type UserAccountRepository interface {
	CreateUserAccount(ctx context.Context, user *UserAccount) (*UserAccount, error)
	GetUserAccountByUserID(ctx context.Context, userID int64) (*UserAccount, error)
	UpdateUserAccount(ctx context.Context, user *UserAccount) (*UserAccount, error)
	DeleteUserAccount(ctx context.Context, userID int64) error
}
