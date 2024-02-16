// services/order_service.go

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

const (
	// NewOrderStatus новый заказ;
	NewOrderStatus string = "NEW"
	// RegisteredOrderStatus заказ зарегистрирован, но начисление не рассчитано;
	RegisteredOrderStatus string = "REGISTERED"
	// InvalidOrderStatus заказ не принят к расчёту, и вознаграждение не будет начислено;
	InvalidOrderStatus string = "INVALID"
	// ProcessingOrderStatus расчёт начисления в процессе;
	ProcessingOrderStatus string = "PROCESSING"
	// ProcessedOrderStatus расчёт начисления окончен;
	ProcessedOrderStatus string = "PROCESSED"
)

// OrderService struct represents the service for orders
type OrderService struct {
	TransactionService *TransactionService
	OrderRepository    repository.OrderRepository
	AccrualService     *accrual.AccrualService
}

// NewOrderService creates a new instance of OrderService
func NewOrderService(transactionService *TransactionService, orderRepository repository.OrderRepository, accrualService *accrual.AccrualService) *OrderService {
	return &OrderService{
		TransactionService: transactionService,
		OrderRepository:    orderRepository,
		AccrualService:     accrualService,
	}
}

func (s *OrderService) RetryOrderProcessing(orderNumber string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = s.OrderProcessing(orderNumber)
		if err == nil {
			return nil
		}
		logger.Logger().Error(
			"order processing: processing order by number",
			zap.String("orderNumber", orderNumber),
			zap.Error(err),
		)
		time.Sleep(time.Millisecond * 100)
	}
	if err != nil {
		_, err = s.CancelOrder(orderNumber)
		if err != nil {
			return fmt.Errorf("cancel order: %w", err)
		}
	}

	return nil
}

func (s *OrderService) OrderProcessing(
	orderNumber string,
) error {
	child, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	order, err := s.GetOrderByNumber(child, orderNumber)
	if err != nil {
		logger.Logger().Error("order service: get order by number", zap.Error(err))
		return err
	}
	// return if this order is already handled
	if order.Status == ProcessedOrderStatus ||
		order.Status == InvalidOrderStatus {
		return nil
	}
	updatedOrder, err := s.SyncOrderWithAccrual(order.Number)
	if err != nil {
		logger.Logger().Error("order service: wait accrual", zap.Error(err))
		return err
	}
	if updatedOrder != nil {
		err = s.TransactionService.AccrualAmount(child, order.UserID, order.Number, updatedOrder.Accrual)
		if err != nil {
			logger.Logger().Error("transfer service: accrual amount", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s OrderService) GetOrderByNumber(context context.Context, number string) (*repository.Order, error) {
	return s.OrderRepository.GetOrderByNumber(context, number)
}

// SyncOrderWithAccrual обновляем статус заказа и начислния исходя после запроса в сервис
// возвращает заказ с обновленными данными
func (s OrderService) SyncOrderWithAccrual(orderNumber string) (*repository.Order, error) {
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

// CancelOrder отмена заказа
func (s OrderService) CancelOrder(orderNumber string) (*repository.Order, error) {
	ctx := context.Background()
	foundOrder, err := s.OrderRepository.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return nil, err
	}

	foundOrder.Status = InvalidOrderStatus

	_, err = s.OrderRepository.UpdateOrder(ctx, foundOrder)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return nil, err
	}

	return foundOrder, nil
}

func (s OrderService) NewOrder(ctx context.Context, orderNumber string, userID int64) error {
	newOrder := &repository.Order{
		Number:     orderNumber,
		UserID:     userID,
		Status:     NewOrderStatus,
		UploadedAt: time.Now(),
	}
	_, err := s.OrderRepository.CreateOrder(ctx, newOrder)
	if err != nil {
		return fmt.Errorf("creating order: %w", err)
	}
	return nil
}
