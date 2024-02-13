package repository

import (
	"context"
	"math/big"
)

const (
	WithdrawOperationType = "withdraw"
)

type Transaction struct {
	TransactionID int64    `json:"transactionId"`
	FromUserID    int64    `json:"fromUserId"`
	ToUserID      int64    `json:"toUserId"`
	Amount        *big.Int `json:"amount"`
	Reason        string   `json:"reason"`
	OperationType string   `json:"operationType"`
}

// TransactionRepository defines the interface for user repository operations.
//
//go:generate mockgen -source=transaction.go -destination=./mock/transaction_mock.go -package=mock
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID int64) (*Transaction, error)
	GetTransactionsByUserIDAndOperationType(ctx context.Context, userID int64, operationType string) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID int64) error
}
