package postgres

import (
	"context"
	"fmt"
	"math/big"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

// UserAccountRepository represents the repository for user accounts using PostgreSQL.
type UserAccountRepository struct {
	db *pgxpool.Pool
}

// NewUserAccountRepository creates a new instance of UserAccountRepository.
func NewUserAccountRepository(db *pgxpool.Pool) *UserAccountRepository {
	return &UserAccountRepository{db: db}
}

// CreateUserAccount creates a new user account in the PostgreSQL database.
func (r *UserAccountRepository) CreateUserAccount(user *repository.UserAccount) (*repository.UserAccount, error) {
	sql := `INSERT INTO user_accounts (user_id, balance) VALUES ($1, $2) RETURNING user_id`
	err := r.db.QueryRow(context.Background(), sql, user.UserID, user.Balance.String()).Scan(&user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user account: %v", err)
	}

	return user, nil
}

// GetUserAccountByUserID retrieves a user account from the PostgreSQL database by user ID.
func (r *UserAccountRepository) GetUserAccountByUserID(userID int64) (*repository.UserAccount, error) {
	sql := `SELECT balance FROM user_accounts WHERE user_id = $1`
	var balance int64
	err := r.db.QueryRow(context.Background(), sql, userID).Scan(&balance)
	if err != nil {
		return nil, fmt.Errorf("failed to get user account: %v", err)
	}

	userAccount := &repository.UserAccount{
		UserID:  userID,
		Balance: big.NewInt(balance),
	}

	return userAccount, nil
}

// UpdateUserAccount updates a user account information in the PostgreSQL database.
func (r *UserAccountRepository) UpdateUserAccount(user *repository.UserAccount) (*repository.UserAccount, error) {
	sql := `UPDATE user_accounts SET balance = $1 WHERE user_id = $2`
	_, err := r.db.Exec(context.Background(), sql, user.Balance.String(), user.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user account: %v", err)
	}

	return user, nil
}

// DeleteUserAccount deletes a user account from the PostgreSQL database by user ID.
func (r *UserAccountRepository) DeleteUserAccount(userID int64) error {
	sql := `DELETE FROM user_accounts WHERE user_id = $1`
	_, err := r.db.Exec(context.Background(), sql, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user account: %v", err)
	}

	return nil
}
