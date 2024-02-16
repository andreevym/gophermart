package postgres_test

import (
	"context"
	"testing"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository(t *testing.T) {
	require.NotNil(t, testDB)
	err := testDB.Ping(context.Background())
	require.NoError(t, err)

	repo := postgres.NewTransactionRepository(testDB)

	// Create a test transaction
	transaction := repository.Transaction{
		FromUserID:    1,
		ToUserID:      1,
		Amount:        100,
		OrderNumber:   "821546",
		OperationType: "Test operation",
	}

	// Test CreateTransaction
	createdTransaction, err := repo.CreateTransaction(context.Background(), transaction)
	require.NoError(t, err)

	require.NotNil(t, createdTransaction)
	require.NotZero(t, createdTransaction.TransactionID)

	transaction.TransactionID = createdTransaction.TransactionID
	require.Equal(t, transaction, createdTransaction)

	// Test GetTransactionByID
	retrievedTransaction, err := repo.GetTransactionByID(context.Background(), createdTransaction.TransactionID)
	require.NoError(t, err)
	require.NotNil(t, retrievedTransaction)
	require.Equal(t, createdTransaction.TransactionID, retrievedTransaction.TransactionID)
	require.Equal(t, createdTransaction.Amount, retrievedTransaction.Amount)
	require.Equal(t, createdTransaction.OrderNumber, retrievedTransaction.OrderNumber)
	require.Equal(t, createdTransaction.OperationType, retrievedTransaction.OperationType)

	// Test GetTransactionByUserIDAndOperationType
	retrievedTransactionByUserAndOperation, err := repo.GetTransactionsByUserIDAndOperationType(context.Background(), 1, "Test operation")
	require.NoError(t, err)
	require.NotNil(t, retrievedTransactionByUserAndOperation)
	require.Len(t, retrievedTransactionByUserAndOperation, 1)
	//)Equal(t, []repository.Transaction{createdTransaction},
	require.Len(t, retrievedTransactionByUserAndOperation, 1)
	require.Equal(t, createdTransaction.TransactionID, retrievedTransactionByUserAndOperation[0].TransactionID)
	require.Equal(t, createdTransaction.Amount, retrievedTransactionByUserAndOperation[0].Amount)
	require.Equal(t, createdTransaction.OrderNumber, retrievedTransactionByUserAndOperation[0].OrderNumber)
	require.Equal(t, createdTransaction.OperationType, retrievedTransactionByUserAndOperation[0].OperationType)
	require.False(t, retrievedTransactionByUserAndOperation[0].Created.IsZero())

	// Test UpdateTransaction
	transactionToUpdate := repository.Transaction{
		TransactionID: createdTransaction.TransactionID,
		FromUserID:    1,
		ToUserID:      1,
		Amount:        200,
		OrderNumber:   "821546",
		OperationType: "Updated operation",
	}
	updatedTransaction, err := repo.UpdateTransaction(context.Background(), transactionToUpdate)
	require.NoError(t, err)
	require.NotNil(t, updatedTransaction)
	require.Equal(t, transactionToUpdate, updatedTransaction)

	// Test DeleteTransaction
	err = repo.DeleteTransaction(context.Background(), createdTransaction.TransactionID)
	require.NoError(t, err)

	// Verify transaction is deleted
	_, err = repo.GetTransactionByID(context.Background(), createdTransaction.TransactionID)
	require.Error(t, err)
	require.EqualError(t, err, postgres.ErrTransactionNotFound.Error())
}
