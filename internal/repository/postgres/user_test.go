package postgres_test

import (
	"context"
	"testing"

	"github.com/andreevym/gophermart/internal/repository"
	"github.com/andreevym/gophermart/internal/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	require.NotNil(t, testDB)
	err := testDB.Ping(context.Background())
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)

	// Create a test user
	user := repository.User{
		Username: "testuser",
		Password: "password",
	}

	// Test CreateUser
	err = repo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	// Test GetUserByUsername
	retrievedUserByUsername, err := repo.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotNil(t, retrievedUserByUsername)
	require.Equal(t, user.Username, retrievedUserByUsername.Username)
	require.Equal(t, user.Password, retrievedUserByUsername.Password)

	// Test GetUserByID
	retrievedUserByID, err := repo.GetUserByID(context.Background(), retrievedUserByUsername.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUserByID)
	require.Equal(t, retrievedUserByUsername, retrievedUserByID)

	// Test UpdateUser
	userToUpdate := repository.User{
		ID:       retrievedUserByUsername.ID,
		Username: "updateduser",
		Password: "newpassword",
	}
	err = repo.UpdateUser(context.Background(), userToUpdate)
	require.NoError(t, err)

	// Test DeleteUser
	err = repo.DeleteUser(context.Background(), retrievedUserByUsername.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = repo.GetUserByID(context.Background(), retrievedUserByUsername.ID)
	require.Error(t, err)
	require.EqualError(t, err, postgres.ErrUserNotFound.Error())
}
