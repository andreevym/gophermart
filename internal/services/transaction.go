package services

import (
	"context"
	"fmt"

	"github.com/andreevym/gophermart/internal/repository"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
}

const (
	WithdrawUserID = 1
	AccrualUserID  = 2
)

func (s TransactionService) Withdraw(ctx context.Context, fromUserID int64, amount float32, orderNumber string) error {
	transaction := repository.Transaction{
		FromUserID:    fromUserID,
		ToUserID:      WithdrawUserID,
		Amount:        amount,
		OperationType: repository.WithdrawOperationType,
		OrderNumber:   orderNumber,
	}
	_, err := s.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("transaction storage: save transaction: %w", err)
	}
	return nil
}

func (s TransactionService) GetCurrentBalance(ctx context.Context, userID int64) (float32, error) {
	transactions, err := s.transactionRepository.GetTransactionsByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get user account by user id '%d': %w", userID, err)
	}
	var balance float32
	for _, transaction := range transactions {
		if transaction.FromUserID == userID {
			balance -= transaction.Amount
		} else if transaction.ToUserID == userID {
			balance += transaction.Amount
		} else {
			panic("should not be here")
		}
	}

	return balance, nil
}

func (s TransactionService) GetWithdrawBalance(ctx context.Context, userID int64) (float32, error) {
	transactions, err := s.transactionRepository.GetTransactionsByUserIDAndOperationType(ctx, userID, repository.WithdrawOperationType)
	if err != nil {
		return 0, fmt.Errorf("get user account by user id '%d': %w", userID, err)
	}
	var balance float32
	for _, transaction := range transactions {
		if transaction.FromUserID == userID {
			balance += transaction.Amount
		}
	}

	return balance, nil
}

func (s TransactionService) GetWithdrawTransaction(ctx context.Context, userID int64) ([]repository.Transaction, error) {
	transactions, err := s.transactionRepository.GetTransactionsByUserIDAndOperationType(ctx, userID, repository.WithdrawOperationType)
	if err != nil {
		return nil, fmt.Errorf("get user account by user id '%d': %w", userID, err)
	}
	return transactions, nil
}

func (s TransactionService) AccrualAmount(ctx context.Context, orderUserID int64, orderNumber string, orderAccrual float32, orderStatus string) error {
	err := s.transactionRepository.AccrualAmount(ctx, orderUserID, orderNumber, orderAccrual, orderStatus)
	if err != nil {
		return fmt.Errorf("failed to accrual amount for userID '%d' and order number %s: %w", orderUserID, orderNumber, err)
	}

	return nil
}

func NewTransactionService(transactionRepository repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}
