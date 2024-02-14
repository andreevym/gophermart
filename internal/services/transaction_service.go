package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/andreevym/gofermart/internal/repository"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
	userAccountRepository repository.UserAccountRepository
}

const (
	SystemUserID = 1
)

func (ts TransactionService) Withdraw(ctx context.Context, fromUserID int64, amount *big.Int, reason string) error {
	fromUserAccount, err := ts.userAccountRepository.GetUserAccountByUserID(ctx, fromUserID)
	if err != nil {
		return fmt.Errorf("userAccountRepository.GetUserAccountByUserID, user %d: %w", fromUserID, err)
	}
	fromUserAccount.Balance = big.NewInt(0).Sub(fromUserAccount.Balance, amount)
	_, err = ts.userAccountRepository.UpdateUserAccount(ctx, fromUserAccount)
	if err != nil {
		return fmt.Errorf("balance storage: save balance, user %d: %w", fromUserAccount.UserID, err)
	}

	toUserAccount, err := ts.userAccountRepository.GetUserAccountByUserID(ctx, SystemUserID)
	if err != nil {
		return fmt.Errorf("userAccountRepository.GetUserAccountByUserID, user %d: %w", SystemUserID, err)
	}
	toUserAccount.Balance = big.NewInt(0).Sub(toUserAccount.Balance, amount)
	_, err = ts.userAccountRepository.UpdateUserAccount(ctx, toUserAccount)
	if err != nil {
		return fmt.Errorf("balance storage: save balance, user %d: %w", toUserAccount.UserID, err)
	}

	transaction := &repository.Transaction{
		FromUserID:    fromUserID,
		ToUserID:      SystemUserID,
		Amount:        amount,
		OperationType: repository.WithdrawOperationType,
		Reason:        reason,
	}
	_, err = ts.transactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return fmt.Errorf("transaction storage: save transaction: %w", err)
	}
	return nil
}

func (ts TransactionService) GetTransactionsByUserAndOperationType(userID int64, operationType string) ([]*repository.Transaction, error) {
	return nil, errors.New("")
}

func NewTransactionService(transactionRepository repository.TransactionRepository,
	userAccountRepository repository.UserAccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		userAccountRepository: userAccountRepository,
	}
}
