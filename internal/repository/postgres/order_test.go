package postgres_test

import (
	"context"
	"testing"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/stretchr/testify/require"
)

func TestOrderRepository(t *testing.T) {
	require.NotNil(t, testDB)
	err := testDB.Ping(context.TODO())
	require.NoError(t, err)

	repo := postgres.NewOrderRepository(testDB)

	// Create a test order
	order := &repository.Order{
		Number: "123456",
		UserID: 1,
		Status: "pending",
	}

	// Test CreateOrder
	createdOrder, err := repo.CreateOrder(context.Background(), order)
	require.NoError(t, err)

	require.NotNil(t, createdOrder)
	require.NotZero(t, createdOrder.ID)

	order.ID = createdOrder.ID
	require.Equal(t, order, createdOrder)

	// Test GetOrderByID
	retrievedOrder, err := repo.GetOrderByID(context.Background(), createdOrder.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedOrder)
	require.Equal(t, createdOrder.ID, retrievedOrder.ID)
	require.Equal(t, createdOrder.UploadedAt, retrievedOrder.UploadedAt)
	require.Equal(t, createdOrder.Accrual, retrievedOrder.Accrual)
	require.Equal(t, createdOrder.Number, retrievedOrder.Number)
	require.Equal(t, createdOrder.UserID, retrievedOrder.UserID)

	// Test GetOrderByNumber
	retrievedOrderByNumber, err := repo.GetOrderByNumber(context.Background(), order.Number)
	require.NoError(t, err)
	require.NotNil(t, retrievedOrderByNumber)
	require.Equal(t, createdOrder, retrievedOrderByNumber)

	// Test UpdateOrder
	orderToUpdate := &repository.Order{
		ID:      createdOrder.ID,
		Number:  "654321",
		UserID:  1,
		Status:  "completed",
		Accrual: 200,
	}
	updatedOrder, err := repo.UpdateOrder(context.Background(), orderToUpdate)
	require.NoError(t, err)
	require.NotNil(t, updatedOrder)
	require.Equal(t, orderToUpdate, updatedOrder)

	// Test GetOrdersByUserID
	ordersByUserID, err := repo.GetOrdersByUserID(context.Background(), orderToUpdate.UserID)
	require.NoError(t, err)
	require.NotNil(t, ordersByUserID)
	require.NotEmpty(t, ordersByUserID)
	require.Contains(t, ordersByUserID, updatedOrder)

	// Test DeleteOrder
	err = repo.DeleteOrder(context.Background(), createdOrder.ID)
	require.NoError(t, err)

	// Verify order is deleted
	_, err = repo.GetOrderByID(context.Background(), createdOrder.ID)
	require.Error(t, err)
	require.EqualError(t, err, postgres.ErrOrderNotFound.Error())
}
