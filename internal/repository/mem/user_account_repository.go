package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/andreevym/gofermart/internal/repository"
)

var (
	ErrUserAccountNotFound      = errors.New("user account not found")
	ErrUserAccountAlreadyExists = errors.New("user account already exists")
)

// MemUserAccountRepository represents an in-memory implementation of UserAccountRepository
type MemUserAccountRepository struct {
	accounts map[int64]*repository.UserAccount
	mu       sync.RWMutex
}

// NewMemUserAccountRepository creates a new instance of MemUserAccountRepository
func NewMemUserAccountRepository() *MemUserAccountRepository {
	return &MemUserAccountRepository{
		accounts: make(map[int64]*repository.UserAccount),
	}
}

func (r *MemUserAccountRepository) CreateUserAccount(_ context.Context, account *repository.UserAccount) (*repository.UserAccount, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.UserID]; exists {
		return nil, ErrUserAccountAlreadyExists
	}

	r.accounts[account.UserID] = account

	return account, nil
}

func (r *MemUserAccountRepository) GetUserAccountByUserID(_ context.Context, userID int64) (*repository.UserAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[userID]
	if !exists {
		return nil, ErrUserAccountNotFound
	}

	return account, nil
}

func (r *MemUserAccountRepository) UpdateUserAccount(_ context.Context, account *repository.UserAccount) (*repository.UserAccount, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.UserID]; !exists {
		return nil, ErrUserAccountNotFound
	}

	r.accounts[account.UserID] = account

	return account, nil
}

func (r *MemUserAccountRepository) DeleteUserAccount(_ context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[userID]; !exists {
		return ErrUserAccountNotFound
	}

	delete(r.accounts, userID)

	return nil
}
