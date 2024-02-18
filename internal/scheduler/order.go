package scheduler

import (
	"fmt"
	"time"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/services"
	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

type AccrualScheduler struct {
	orderService     *services.OrderService
	accrualService   *accrual.AccrualService
	stop             chan struct{}
	done             chan struct{}
	pollOrdersDelay  time.Duration
	maxOrderAttempts int
}

func NewAccrualScheduler(
	accrualService *accrual.AccrualService,
	orderService *services.OrderService,
	pollOrdersDelay time.Duration,
	maxOrderAttempts int,
) *AccrualScheduler {
	var (
		stop = make(chan struct{}) // tells the goroutine to stop
		done = make(chan struct{}) // tells us that the goroutine exited
	)
	s := &AccrualScheduler{
		accrualService:   accrualService,
		orderService:     orderService,
		stop:             stop,
		done:             done,
		pollOrdersDelay:  pollOrdersDelay,
		maxOrderAttempts: maxOrderAttempts,
	}

	return s
}

func (s *AccrualScheduler) Run() {
	go s.processingByDelay(s.done, s.stop, s.pollOrdersDelay, s.maxOrderAttempts)
}

func (s *AccrualScheduler) processingByDelay(done chan struct{}, stop chan struct{}, pollOrdersDelay time.Duration, maxOrderAttempts int) {
	defer close(done)
	go func() {
		ticker := time.NewTicker(pollOrdersDelay)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case t := <-ticker.C:
				if err := s.syncOrders(t, maxOrderAttempts); err != nil {
					logger.Logger().Error("sync orders", zap.Error(err))
				}
			}
		}
	}()
}

// syncOrders sync new orders statuses and accrual with accrual service
func (s *AccrualScheduler) syncOrders(t time.Time, maxOrderAttempts int) error {
	logger.Logger().Debug("poll orders", zap.String("ticker", t.String()))
	orders, err := s.orderService.GetOrdersByStatus(services.NewOrderStatus)
	if err != nil {
		logger.Logger().Error("get orders by status", zap.Error(err))
		return fmt.Errorf("get orders by status %w", err)
	}
	for _, order := range orders {
		err = s.orderService.OrderProcessingWithRetry(order, maxOrderAttempts)
		if err != nil {
			logger.Logger().Error("RetryOrderProcessing", zap.Error(err))
			panic(err.Error())
		}
	}
	return nil
}

// Shutdown tells the worker to stop
// and waits until it has finished.
func (s *AccrualScheduler) Shutdown() {
	close(s.stop)
	<-s.done
}
