package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

var ErrAccrualServiceDisabled = errors.New("accrual service is disabled because url is not set")

type AccrualService struct {
	url string
}

func NewAccrualService(url string) *AccrualService {
	return &AccrualService{url: url}
}

type OrderAccrual struct {
	// Number номер заказа
	Order string `json:"order"`
	// Status статус расчёта начисления
	Status string `json:"status"`
	// Accrual рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	Accrual float32 `json:"accrual"`
}

// GetOrderByNumber получение информации о расчёте начислений баллов лояльности.
func (as AccrualService) GetOrderByNumber(orderNumber string) (*OrderAccrual, error) {
	if as.url == "" {
		return nil, ErrAccrualServiceDisabled
	}
	url := fmt.Sprintf("%s/api/orders/%s", as.url, orderNumber)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Post: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Logger().Error("failed to close resp body", zap.Error(err))
		}
	}(response.Body)

	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response io.ReadAll: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request: %s", url)
	}

	orderAccrual := OrderAccrual{}
	err = json.Unmarshal(readAll, &orderAccrual)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &orderAccrual, nil
}
