package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ErrTransactionNotFound represents an error when a transaction is not found in the database.
var ErrTransactionNotFound = errors.New("transaction not found")

// TransactionRepository represents the repository for transactions using PostgreSQL.
type TransactionRepository struct {
	db *pgxpool.Pool
}

// NewTransactionRepository creates a new instance of TransactionRepository.
func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction creates a new transaction in the PostgreSQL database.
func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	var transactionID int64
	sql := `INSERT INTO transactions (from_user_id, to_user_id, amount, reason, operation_type) 
		VALUES ($1, $2, $3, $4, $5) RETURNING transaction_id`
	err := r.db.QueryRow(ctx, sql, transaction.FromUserID, transaction.ToUserID, transaction.Amount.String(), transaction.Reason, transaction.OperationType).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	transaction.TransactionID = transactionID
	return transaction, nil
}

// GetTransactionByID retrieves a transaction from the PostgreSQL database by its ID.
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, transactionID int64) (*repository.Transaction, error) {
	sql := `SELECT from_user_id, to_user_id, amount, reason, operation_type 
		FROM transactions WHERE transaction_id = $1`
	var transaction repository.Transaction
	var amount int64
	err := r.db.QueryRow(ctx, sql, transactionID).Scan(&transaction.FromUserID, &transaction.ToUserID, &amount, &transaction.Reason,
		&transaction.OperationType)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	transaction.Amount = big.NewInt(amount)
	transaction.TransactionID = transactionID

	return &transaction, nil
}

// GetTransactionsByUserIDAndOperationType retrieves a transaction from the PostgreSQL database by user ID and operation type.
func (r *TransactionRepository) GetTransactionsByUserIDAndOperationType(ctx context.Context, userID int64, operationType string) ([]*repository.Transaction, error) {
	sql := `SELECT transaction_id, from_user_id, to_user_id, amount, reason 
		FROM transactions WHERE (from_user_id = $1 OR to_user_id = $1) AND operation_type = $2`
	rows, err := r.db.Query(ctx, sql, userID, operationType)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	defer rows.Close()

	transactions := make([]*repository.Transaction, 0)
	for rows.Next() {
		var transaction repository.Transaction
		var amount int64
		err := rows.Scan(&transaction.TransactionID, &transaction.FromUserID, &transaction.ToUserID,
			&amount, &transaction.Reason)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}

		transaction.Amount = big.NewInt(amount)
		transaction.OperationType = operationType

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	sql := `UPDATE transactions SET from_user_id = $1, to_user_id = $2, amount = $3, reason = $4, operation_type = $5 WHERE transaction_id = $6`
	_, err := r.db.Exec(ctx, sql, transaction.FromUserID, transaction.ToUserID, transaction.Amount.String(), transaction.Reason, transaction.OperationType, transaction.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %v", err)
	}

	return transaction, nil
}

// DeleteTransaction deletes a transaction from the PostgreSQL database by its ID.
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, transactionID int64) error {
	sql := `DELETE FROM transactions WHERE transaction_id = $1`
	_, err := r.db.Exec(ctx, sql, transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	return nil
}
