package postgres_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestUserAccountRepository(t *testing.T) {
	require.NotNil(t, testDB)
	err := testDB.Ping(context.Background())
	require.NoError(t, err)

	repo := postgres.NewUserAccountRepository(testDB)

	// Create a test user account
	userAccount := &repository.UserAccount{
		UserID:  1,
		Balance: big.NewInt(100),
	}

	// Test CreateUserAccount
	createdUserAccount, err := repo.CreateUserAccount(userAccount)
	require.NoError(t, err)

	require.NotNil(t, createdUserAccount)
	require.Equal(t, userAccount.UserID, createdUserAccount.UserID)
	require.Equal(t, userAccount.Balance, createdUserAccount.Balance)

	// Test GetUserAccountByUserID
	retrievedUserAccount, err := repo.GetUserAccountByUserID(userAccount.UserID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUserAccount)
	require.Equal(t, createdUserAccount, retrievedUserAccount)

	// Update user account balance
	userAccount.Balance.SetInt64(200)

	// Test UpdateUserAccount
	updatedUserAccount, err := repo.UpdateUserAccount(userAccount)
	require.NoError(t, err)
	require.NotNil(t, updatedUserAccount)
	require.Equal(t, userAccount, updatedUserAccount)

	// Test DeleteUserAccount
	err = repo.DeleteUserAccount(userAccount.UserID)
	require.NoError(t, err)

	// Verify user account is deleted
	_, err = repo.GetUserAccountByUserID(userAccount.UserID)
	require.Error(t, err)
}
