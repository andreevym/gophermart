// services/order_service.go

package services

import (
	"context"

	"github.com/andreevym/gofermart/internal/repository"
)

// OrderService struct represents the service for orders
type OrderService struct {
	OrderRepository repository.OrderRepository
}

func (s OrderService) GetOrderByNumber(context context.Context, number string) (*repository.Order, error) {
	return s.OrderRepository.GetOrderByNumber(context, number)
}

// NewOrderService creates a new instance of OrderService
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{OrderRepository: repo}
}
