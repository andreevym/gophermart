package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/andreevym/gofermart/internal/repository"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

// MemTransactionRepository represents an in-memory implementation of TransactionRepository
type MemTransactionRepository struct {
	transactions map[int64]*repository.Transaction
	mu           sync.RWMutex
}

// NewMemTransactionRepository creates a new instance of MemTransactionRepository
func NewMemTransactionRepository() *MemTransactionRepository {
	return &MemTransactionRepository{
		transactions: make(map[int64]*repository.Transaction),
	}
}

func (r *MemTransactionRepository) CreateTransaction(_ context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.TransactionID]; exists {
		return nil, errors.New("transaction already exists")
	}

	r.transactions[transaction.TransactionID] = transaction

	return transaction, nil
}

func (r *MemTransactionRepository) GetTransactionByID(_ context.Context, transactionID int64) (*repository.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[transactionID]
	if !exists {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}

func (r *MemTransactionRepository) GetTransactionsByUserIDAndOperationType(_ context.Context, userID int64, operationType string) ([]*repository.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions := make([]*repository.Transaction, 0)
	for _, transaction := range r.transactions {
		if (transaction.ToUserID == userID ||
			transaction.FromUserID == userID) &&
			transaction.OperationType == operationType {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

func (r *MemTransactionRepository) UpdateTransaction(_ context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.TransactionID]; !exists {
		return nil, ErrTransactionNotFound
	}

	r.transactions[transaction.TransactionID] = transaction

	return transaction, nil
}

func (r *MemTransactionRepository) DeleteTransaction(_ context.Context, transactionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transactionID]; !exists {
		return ErrTransactionNotFound
	}

	delete(r.transactions, transactionID)

	return nil
}
