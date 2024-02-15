package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	WithdrawUserID              = 1
	AccrualUserID               = 2
	ProcessedOrderStatus string = "PROCESSED"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	var transactionID int64
	sql := `INSERT INTO transactions (from_user_id, to_user_id, amount, reason, operation_type) VALUES ($1, $2, $3, $4, $5) RETURNING transaction_id`
	err := r.db.QueryRow(ctx, sql, transaction.FromUserID, transaction.ToUserID, transaction.Amount, transaction.Reason, transaction.OperationType).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	transaction.TransactionID = transactionID
	return transaction, nil
}

func (r TransactionRepository) AccrualAmount(ctx context.Context, userID int64, orderID int64, accrual float32) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}
	defer tx.Commit(ctx)

	insertTxsql := `INSERT INTO transactions (from_user_id, to_user_id, amount, reason, operation_type) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, insertTxsql, AccrualUserID, userID, accrual, strconv.FormatInt(orderID, 10), repository.WithdrawOperationType)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	sql := `UPDATE orders SET status = $1, accrual = $2 WHERE id = $3`
	_, err = tx.Exec(ctx, sql, ProcessedOrderStatus, accrual, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order, sql %s: %v", sql, err)
	}

	return nil
}

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, transactionID int64) (*repository.Transaction, error) {
	sql := `SELECT from_user_id, to_user_id, amount, reason, operation_type FROM transactions WHERE transaction_id = $1`
	var transaction repository.Transaction
	err := r.db.QueryRow(ctx, sql, transactionID).Scan(&transaction.FromUserID, &transaction.ToUserID, &transaction.Amount, &transaction.Reason,
		&transaction.OperationType)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}

	transaction.TransactionID = transactionID

	return &transaction, nil
}

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
		err := rows.Scan(&transaction.TransactionID, &transaction.FromUserID, &transaction.ToUserID,
			&transaction.Amount, &transaction.Reason)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}

		transaction.OperationType = operationType
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID int64) ([]*repository.Transaction, error) {
	sql := `SELECT transaction_id, from_user_id, to_user_id, amount, reason, operation_type
		FROM transactions WHERE from_user_id = $1 OR to_user_id = $1`
	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	defer rows.Close()

	transactions := make([]*repository.Transaction, 0)
	for rows.Next() {
		var transaction repository.Transaction
		err := rows.Scan(&transaction.TransactionID, &transaction.FromUserID, &transaction.ToUserID,
			&transaction.Amount, &transaction.Reason, &transaction.OperationType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction *repository.Transaction) (*repository.Transaction, error) {
	sql := `UPDATE transactions SET from_user_id = $1, to_user_id = $2, amount = $3, reason = $4, operation_type = $5 WHERE transaction_id = $6`
	_, err := r.db.Exec(ctx, sql, transaction.FromUserID, transaction.ToUserID, transaction.Amount, transaction.Reason, transaction.OperationType, transaction.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %v", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) DeleteTransaction(ctx context.Context, transactionID int64) error {
	sql := `DELETE FROM transactions WHERE transaction_id = $1`
	_, err := r.db.Exec(ctx, sql, transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}

	return nil
}
