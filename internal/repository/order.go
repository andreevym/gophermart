package repository

import (
	"context"
	"time"
)

// Order represents an order entity in the application.
type Order struct {
	Number     string    `json:"number"`
	UserID     int64     `json:"userId"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// OrderRepository represents the interface for order repository operations.
//
//go:generate mockgen -source=order.go -destination=./mock/order.go -package=mock
type OrderRepository interface {
	CreateOrder(ctx context.Context, order Order) error
	UpdateOrder(ctx context.Context, order Order) error
	DeleteOrder(ctx context.Context, orderNumber string) error

	GetOrderByNumber(ctx context.Context, number string) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]Order, error)
	GetOrdersByStatus(ctx context.Context, status string) ([]Order, error)
}
