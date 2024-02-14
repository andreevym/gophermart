package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/andreevym/gofermart/internal/repository"
)

type UserAccountService struct {
	userAccountRepository repository.UserAccountRepository
	transactionRepository repository.TransactionRepository
}

func NewUserAccountService(userAccountRepository repository.UserAccountRepository, transactionRepository repository.TransactionRepository) *UserAccountService {
	return &UserAccountService{
		userAccountRepository: userAccountRepository,
		transactionRepository: transactionRepository,
	}
}

func (s UserAccountService) GetCurrentBalance(ctx context.Context, userID int64) (*big.Int, error) {
	userAccount, err := s.userAccountRepository.GetUserAccountByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user account by user id '%d': %w", userID, err)
	}

	return userAccount.Balance, nil
}

func (s UserAccountService) GetWithdrawAmount(ctx context.Context, userID int64) (*big.Int, error) {
	transactions, err := s.transactionRepository.GetTransactionsByUserIDAndOperationType(ctx, userID, repository.WithdrawOperationType)
	if err != nil {
		return nil, fmt.Errorf("get user account by user id '%d': %w", userID, err)
	}

	var withdrawAmount *big.Int
	for _, transaction := range transactions {
		if transaction.FromUserID == userID {
			withdrawAmount = big.NewInt(0).Add(withdrawAmount, transaction.Amount)
		}
	}

	return withdrawAmount, nil
}
