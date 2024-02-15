package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreevym/gofermart/internal/repository"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	var orderID int64
	if order.UploadedAt.IsZero() {
		sql := `INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3) RETURNING id`
		err := r.db.QueryRow(ctx, sql, order.Number, order.UserID, order.Status).Scan(&orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %v", err)
		}
	} else {
		sql := `INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err := r.db.QueryRow(ctx, sql, order.Number, order.UserID, order.Status, order.UploadedAt).Scan(&orderID)
		if err != nil {
			return nil, fmt.Errorf("failed to create order: %v", err)
		}
	}

	order.ID = orderID
	return order, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*repository.Order, error) {
	sql := `SELECT number, user_id, status, accrual, uploaded_at FROM orders WHERE id = $1`
	var order repository.Order
	var accrual pgtype.Float4
	var uploadedAtNullable pgtype.Timestamptz
	err := r.db.QueryRow(ctx, sql, orderID).Scan(&order.Number, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	order.ID = orderID
	if accrual.Status == pgtype.Present {
		order.Accrual = accrual.Float
	}

	if uploadedAtNullable.Status == pgtype.Present {
		order.UploadedAt = uploadedAtNullable.Time
	} else {
		order.UploadedAt = time.Time{}
	}

	return &order, nil
}

func (r *OrderRepository) GetOrderByNumber(ctx context.Context, number string) (*repository.Order, error) {
	sql := `SELECT id, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1`
	var order repository.Order
	var accrual pgtype.Float4
	var uploadedAtNullable pgtype.Timestamptz
	err := r.db.QueryRow(ctx, sql, number).Scan(&order.ID, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}

	order.Number = number
	if accrual.Status == pgtype.Present {
		order.Accrual = accrual.Float
	}

	if uploadedAtNullable.Status == pgtype.Present {
		order.UploadedAt = uploadedAtNullable.Time
	} else {
		order.UploadedAt = time.Time{}
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	if order.UploadedAt.IsZero() {
		sql := `UPDATE orders SET number = $1, user_id = $2, status = $3, accrual = $4 WHERE id = $5`
		_, err := r.db.Exec(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update order, sql %s: %v", sql, err)
		}
	} else {
		sql := `UPDATE orders SET number = $1, user_id = $2, status = $3, accrual = $4, uploaded_at = $5 WHERE id = $6`
		_, err := r.db.Exec(ctx, sql, order.Number, order.UserID, order.Status, order.Accrual, order.UploadedAt, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update order, sql %s: %v", sql, err)
		}
	}

	return order, nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, orderID int64) error {
	sql := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %v", err)
	}

	return nil
}

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
		var accrual pgtype.Float4
		err := rows.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &accrual, &uploadedAtNullable)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order row: %v", err)
		}

		if accrual.Status == pgtype.Present {
			order.Accrual = accrual.Float
		}

		if uploadedAtNullable.Status == pgtype.Present {
			order.UploadedAt = uploadedAtNullable.Time
		} else {
			order.UploadedAt = time.Time{}
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over order rows: %v", err)
	}

	return orders, nil
}
