package postgres_test

import (
	"context"
	"testing"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	require.NotNil(t, testDB)
	err := testDB.Ping(context.Background())
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)

	// Create a test user
	user := &repository.User{
		Username: "testuser",
		Password: "password",
	}

	// Test CreateUser
	createdUser, err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotZero(t, createdUser.ID)

	user.ID = createdUser.ID
	require.Equal(t, user, createdUser)

	// Test GetUserByID
	retrievedUserByID, err := repo.GetUserByID(context.Background(), createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedUserByID)
	require.Equal(t, createdUser, retrievedUserByID)

	// Test GetUserByUsername
	retrievedUserByUsername, err := repo.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotNil(t, retrievedUserByUsername)
	require.Equal(t, createdUser, retrievedUserByUsername)

	// Test UpdateUser
	userToUpdate := &repository.User{
		ID:       createdUser.ID,
		Username: "updateduser",
		Password: "newpassword",
	}
	updatedUser, err := repo.UpdateUser(context.Background(), userToUpdate)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Equal(t, userToUpdate, updatedUser)

	// Test DeleteUser
	err = repo.DeleteUser(context.Background(), createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = repo.GetUserByID(context.Background(), createdUser.ID)
	require.Error(t, err)
	require.EqualError(t, err, postgres.ErrUserNotFound.Error())
}
