// services/order_service.go

package services

import (
	"context"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

// OrderService struct represents the service for orders
type OrderService struct {
	OrderRepository repository.OrderRepository
	AccrualService  *accrual.AccrualService
}

func (s OrderService) GetOrderByNumber(context context.Context, number string) (*repository.Order, error) {
	return s.OrderRepository.GetOrderByNumber(context, number)
}

// WaitAccrual ждем расчета начислений по заказу и переводим статус заказа
func (s OrderService) WaitAccrual(orderNumber string) (*repository.Order, error) {
	if s.AccrualService == nil {
		logger.Logger().Warn("WaitAccrual", zap.Error(accrual.ErrAccrualServiceDisabled))
		return nil, nil
	}
	ctx := context.Background()
	foundOrder, err := s.OrderRepository.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return nil, err
	}

	orderAccrual, err := s.AccrualService.GetOrderByNumber(orderNumber)
	if err != nil {
		logger.Logger().Error("AccrualService.GetOrderByNumber", zap.Error(err))
		return nil, err
	}

	foundOrder.Status = orderAccrual.Status
	foundOrder.Accrual = orderAccrual.Accrual

	_, err = s.OrderRepository.UpdateOrder(ctx, foundOrder)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return nil, err
	}

	return foundOrder, nil
}

// NewOrderService creates a new instance of OrderService
func NewOrderService(repo repository.OrderRepository, accrualService *accrual.AccrualService) *OrderService {
	return &OrderService{
		OrderRepository: repo,
		AccrualService:  accrualService,
	}
}
