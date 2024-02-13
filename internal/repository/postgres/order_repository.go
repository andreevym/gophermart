package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// ErrOrderNotFound represents an error when an order is not found in the database.
	ErrOrderNotFound = errors.New("order not found")
)

// OrderRepository represents the repository for orders using PostgreSQL.
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository creates a new instance of OrderRepository.
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order in the PostgreSQL database.
func (r *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	var orderID int64
	if order.UploadedAt.IsZero() {
		sql := `INSERT INTO orders (number, user_id, status, accrual) VALUES ($1, $2, $3, $4) RETURNING id`
		err := r.db.QueryRow(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual.String()).Scan(&orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %v", err)
		}
	} else {
		sql := `INSERT INTO orders (number, user_id, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err := r.db.QueryRow(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual.String(), order.UploadedAt).Scan(&orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %v", err)
		}
	}

	order.ID = orderID
	return order, nil
}

// GetOrderByID retrieves an order from the PostgreSQL database by its ID.
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*repository.Order, error) {
	sql := `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE id = $1`
	var order repository.Order
	var accrual int64                         // Change type to string
	var uploadedAtNullable pgtype.Timestamptz // Use pgtype for nullable time.Time
	err := r.db.QueryRow(ctx, sql, orderID).Scan(&order.Number, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	order.ID = orderID
	order.Accrual = big.NewInt(accrual)

	// Check if uploaded_at is NULL
	if uploadedAtNullable.Status == pgtype.Present {
		order.UploadedAt = uploadedAtNullable.Time // Assign uploaded_at if not NULL
	} else {
		order.UploadedAt = time.Time{} // Set to zero time if NULL
	}

	return &order, nil
}

// GetOrderByNumber retrieves an order from the PostgreSQL database by its number.
func (r *OrderRepository) GetOrderByNumber(ctx context.Context, number string) (*repository.Order, error) {
	sql := `SELECT id, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1`
	var order repository.Order
	var accrual int64                         // Change type to string
	var uploadedAtNullable pgtype.Timestamptz // Use pgtype for nullable time.Time
	err := r.db.QueryRow(ctx, sql, number).Scan(&order.ID, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	order.Number = number
	order.Accrual = big.NewInt(accrual)

	// Check if uploaded_at is NULL
	if uploadedAtNullable.Status == pgtype.Present {
		order.UploadedAt = uploadedAtNullable.Time // Assign uploaded_at if not NULL
	} else {
		order.UploadedAt = time.Time{} // Set to zero time if NULL
	}

	return &order, nil
}

// UpdateOrder updates order information in the PostgreSQL database.
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	if order.UploadedAt.IsZero() {
		sql := `UPDATE orders SET number = $1, user_id = $2, status = $3, accrual = $4 WHERE id = $5`
		_, err := r.db.Exec(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual.String(), order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update order: %v", err)
		}
	} else {
		sql := `UPDATE orders SET number = $1, user_id = $2, status = $3, accrual = $4, uploaded_at = $5 WHERE id = $6`
		_, err := r.db.Exec(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual.String(), order.UploadedAt, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update order: %v", err)
		}
	}

	return order, nil
}

// DeleteOrder deletes an order from the PostgreSQL database by its ID.
func (r *OrderRepository) DeleteOrder(ctx context.Context, orderID int64) error {
	sql := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}

	return nil
}

// GetOrdersByUserID retrieves a list of orders from the PostgreSQL database for the specified user.
func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]*repository.Order, error) {
	sql := `SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := r.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer rows.Close()

	orders := make([]*repository.Order, 0)
	for rows.Next() {
		var uploadedAtNullable pgtype.Timestamptz
		var order repository.Order
		var accrual int64
		err := rows.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order row: %v", err)
		}

		order.Accrual = big.NewInt(accrual)

		// Check if uploaded_at is NULL
		if uploadedAtNullable.Status == pgtype.Present {
			order.UploadedAt = uploadedAtNullable.Time // Assign uploaded_at if not NULL
		} else {
			order.UploadedAt = time.Time{} // Set to zero time if NULL
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order rows: %v", err)
	}

	return orders, nil
}
