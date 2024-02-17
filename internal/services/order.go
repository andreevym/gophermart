// services/order_service.go

package services

import (
	"context"
	"errors"
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

var ErrAccrualServiceDisabled = errors.New("accrual service is disabled")

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

func (s *OrderService) OrderProcessingWithRetry(order repository.Order, maxOrderAttempts int) error {
	var err error
	for i := 0; i < maxOrderAttempts; i++ {
		err = s.OrderProcessing(order)
		if err == nil {
			return nil
		}
		logger.Logger().Error(
			"order processing: failed to process order by number, with retry",
			zap.String("orderNumber", order.Number),
			zap.Int("attempt", i),
			zap.Error(err),
		)
		time.Sleep(time.Millisecond * 100)
	}
	if err != nil {
		_, err = s.CancelOrder(order)
		if err != nil {
			return fmt.Errorf("failed to process, order was canceled: %w", err)
		}
	}

	return nil
}

func (s *OrderService) OrderProcessing(
	order repository.Order,
) error {
	child, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// return if this order is already handled
	if order.Status == ProcessedOrderStatus ||
		order.Status == InvalidOrderStatus {
		return nil
	}
	updatedOrder, err := s.SyncOrderWithAccrual(order)
	if err != nil {
		logger.Logger().Error("order service: wait accrual", zap.Error(err))
		return fmt.Errorf("failed to sync order with accrual: %w", err)
	}
	err = s.TransactionService.AccrualAmount(child, updatedOrder)
	if err != nil {
		logger.Logger().Error("transfer service: accrual amount", zap.Error(err))
		return fmt.Errorf("failed to make changes in accrual %w", err)
	}

	return nil
}

func (s OrderService) GetOrderByNumber(context context.Context, number string) (repository.Order, error) {
	return s.OrderRepository.GetOrderByNumber(context, number)
}

// SyncOrderWithAccrual обновляем статус заказа и начислния исходя после запроса в сервис
// возвращает заказ с обновленными данными
func (s OrderService) SyncOrderWithAccrual(order repository.Order) (repository.Order, error) {
	ctx := context.Background()

	orderAccrual, err := s.AccrualService.GetOrderByNumber(order.Number)
	if err != nil {
		logger.Logger().Error("AccrualService.GetOrderByNumber", zap.Error(err))
		return repository.Order{}, err
	}

	order.Status = orderAccrual.Status
	order.Accrual = orderAccrual.Accrual

	_, err = s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return repository.Order{}, err
	}

	return order, nil
}

// CancelOrder отмена заказа
func (s OrderService) CancelOrder(order repository.Order) (repository.Order, error) {
	ctx := context.Background()
	order.Status = InvalidOrderStatus

	_, err := s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return repository.Order{}, err
	}

	return order, nil
}

func (s OrderService) NewOrder(ctx context.Context, orderNumber string, userID int64) error {
	newOrder := repository.Order{
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

func (s *OrderService) GetOrdersByStatus(status string) ([]repository.Order, error) {
	ctx := context.Background()
	orders, err := s.OrderRepository.GetOrdersByStatus(ctx, status)
	if err != nil {
		logger.Logger().Error("orderService.OrderRepository.GetOrderByID", zap.Error(err))
		return nil, err
	}
	return orders, nil
}
