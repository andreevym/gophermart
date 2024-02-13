package mem

import (
	"context"
	"errors"
	"sync"

	"github.com/andreevym/gofermart/internal/repository"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

// MemOrderRepository represents an in-memory implementation of OrderRepository
type MemOrderRepository struct {
	orders map[int64]*repository.Order
	mu     sync.RWMutex
}

// NewMemOrderRepository creates a new instance of MemOrderRepository
func NewMemOrderRepository() *MemOrderRepository {
	return &MemOrderRepository{
		orders: make(map[int64]*repository.Order),
	}
}

func (r *MemOrderRepository) CreateOrder(_ context.Context, order *repository.Order) (*repository.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return nil, ErrOrderAlreadyExists
	}

	r.orders[order.ID] = order

	return order, nil
}

func (r *MemOrderRepository) GetOrderByID(_ context.Context, orderID int64) (*repository.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[orderID]
	if !exists {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (r *MemOrderRepository) UpdateOrder(_ context.Context, order *repository.Order) (*repository.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return nil, ErrOrderNotFound
	}

	r.orders[order.ID] = order

	return order, nil
}

func (r *MemOrderRepository) DeleteOrder(_ context.Context, orderID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[orderID]; !exists {
		return ErrOrderNotFound
	}

	delete(r.orders, orderID)

	return nil
}

func (r *MemOrderRepository) GetOrdersByUserID(_ context.Context, userID int64) ([]*repository.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userOrders []*repository.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	return userOrders, nil
}

func (r *MemOrderRepository) GetOrderByNumber(_ context.Context, number string) (*repository.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, order := range r.orders {
		if order.Number == number {
			return order, nil
		}
	}

	return nil, nil
}
