package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// RequestAccrualByOrderNumber получение информации о расчёте начислений баллов лояльности.
func (as AccrualService) RequestAccrualByOrderNumber(orderNumber string) (*OrderAccrual, error) {
	if as.url == "" {
		return nil, ErrAccrualServiceDisabled
	}
	url := fmt.Sprintf("%s/api/orders/%s", as.url, orderNumber)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Post: %w", err)
	}

	defer response.Body.Close()

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

	if orderAccrual.Order != orderNumber {
		return nil, fmt.Errorf(
			"failed to get order from AccrualService: received wrong order number %s, bug expected %s",
			orderAccrual.Order,
			orderNumber,
		)
	}

	return &orderAccrual, nil
}
