package repository

import (
	"context"
	"time"
)

const (
	WithdrawOperationType = "withdraw"
	AccrualOperationType  = "withdraw"
)

type Transaction struct {
	TransactionID int64   `json:"transactionId"`
	FromUserID    int64   `json:"fromUserId"`
	ToUserID      int64   `json:"toUserId"`
	Amount        float32 `json:"amount"`
	OrderNumber   string  `json:"order_number"`
	OperationType string  `json:"operationType"`
	// Created is the combined date and time, filled by database while insert
	Created time.Time `json:"created,omitempty"`
}

// TransactionRepository defines the interface for user repository operations.
//
//go:generate mockgen -source=transaction.go -destination=./mock/transaction.go -package=mock
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction Transaction) (Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID int64) (Transaction, error)
	GetTransactionsByUserIDAndOperationType(ctx context.Context, userID int64, operationType string) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, transaction Transaction) (Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID int64) error
	AccrualAmount(ctx context.Context, userID int64, orderNumber string, accrual float32) error
	GetTransactionsByUserID(ctx context.Context, userID int64) ([]Transaction, error)
}
