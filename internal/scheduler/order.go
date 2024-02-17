package scheduler

import (
	"time"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/services"
	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

type Scheduler struct {
	orderService   *services.OrderService
	accrualService *accrual.AccrualService
}

func NewScheduler(accrualService *accrual.AccrualService, orderService *services.OrderService) *Scheduler {
	return &Scheduler{
		accrualService: accrualService,
		orderService:   orderService,
	}
}

func (s *Scheduler) ProcessingByDelay(done chan struct{}, stop chan struct{}, pollOrdersDelay time.Duration, maxOrderAttempts int) {
	defer close(done)
	go func() {
		ticker := time.NewTicker(pollOrdersDelay)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case t := <-ticker.C:
				logger.Logger().Debug("poll orders", zap.String("ticker", t.String()))
				orders, err := s.orderService.GetOrdersByStatus(services.NewOrderStatus)
				if err != nil {
					logger.Logger().Error("get orders by status", zap.Error(err))
					return
				}
				for _, order := range orders {
					err = s.orderService.OrderProcessingWithRetry(order, maxOrderAttempts)
					if err != nil {
						logger.Logger().Error("RetryOrderProcessing", zap.Error(err))
						panic(err.Error())
					}
				}
			}
		}
	}()
}
