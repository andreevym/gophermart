package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

type userWithdrawal struct {
	OrderWithdrawNumber string    `json:"order"` // номер заказа к которому привязан вывод средств
	Sum                 float32   `json:"sum"`   // сумма баллов к списанию в счёт оплаты
	ProcessedAt         time.Time `json:"processed_at"`
}

type WithdrawRequestDTO struct {
	OrderWithdrawNumber string  `json:"order"` // номер заказа к которому привязан вывод средств
	Sum                 float32 `json:"sum"`   // сумма баллов к списанию в счёт оплаты
}

// GetWithdrawalsHandler получение информации о выводе средств
// #### **Получение информации о выводе средств**
//
// Хендлер: `GET /api/user/withdrawals`.
//
// Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.
//
// Формат запроса:
//
// GET /api/user/withdrawals HTTP/1.1
// Content-Length: 0
//
// Возможные коды ответа:
//
// *   `200` — успешная обработка запроса.
//
// Формат ответа:
//
// 200 OK HTTP/1.1
// Content-Type: application/json
// ...
//
// [
// {
// "order": "2377225624",
// "sum": 500,
// "processed_at": "2020-12-09T16:09:57+03:00"
// }
// ]
//
// *   `204` — нет ни одного списания.
// *   `401` — пользователь не авторизован.
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionService.GetWithdrawTransaction(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userWithdrawals := make([]*userWithdrawal, 0)
	for _, transaction := range transactions {
		userWithdrawals = append(userWithdrawals, &userWithdrawal{
			OrderWithdrawNumber: transaction.OrderNumber,
			Sum:                 transaction.Amount,
			ProcessedAt:         transaction.Created,
		})
	}
	bytes, err := json.Marshal(userWithdrawals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PostWithdrawHandler запрос на списание средств
// #### **Запрос на списание средств**
//
// Хендлер: `POST /api/user/balance/withdraw`
//
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя, в счёт оплаты которого списываются баллы.
//
// Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
//
// Формат запроса:
//
// POST /api/user/balance/withdraw HTTP/1.1
// Content-Type: application/json
//
// {
// "order": "2377225624",
// "sum": 751
// }
//
// Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.
//
// Возможные коды ответа:
//
// *   `200` — успешная обработка запроса;
// *   `401` — пользователь не авторизован;
// *   `402` — на счету недостаточно средств;
// *   `422` — неверный номер заказа;
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) PostWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		logger.Logger().Debug("middleware.GetUserID", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Logger().Debug("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var withdrawRequestDTO WithdrawRequestDTO
	err = json.Unmarshal(bytes, &withdrawRequestDTO)
	if err != nil {
		logger.Logger().Debug("json.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.transactionService.Withdraw(
		context.Background(),
		userID,
		withdrawRequestDTO.Sum,
		withdrawRequestDTO.OrderWithdrawNumber,
	)
	if err != nil {
		logger.Logger().Debug("transactionService.Withdraw", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type BalanceDTO struct {
	Current   float32 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type GetBalanceResponseDTO struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type GetWithdrawResponseDTO struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

// GetBalanceHandler получение текущего баланса пользователя
// #### **Получение текущего баланса пользователя**
//
// Хендлер: `GET /api/user/balance`.
//
// Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
//
// Формат запроса:
//
// GET /api/user/balance HTTP/1.1
// Content-Length: 0
//
// Возможные коды ответа:
//
// *   `200` — успешная обработка запроса.
//
// Формат ответа:
//
// 200 OK HTTP/1.1
// Content-Type: application/json
// ...
//
// {
// "current": 500.5,
// "withdrawn": 42
// }
//
// *   `401` — пользователь не авторизован.
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		logger.Logger().Warn("middleware.GetUserID", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentBalance, err := h.transactionService.GetCurrentBalance(ctx, userID)
	if err != nil {
		logger.Logger().Warn("GetCurrentBalance", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawBalance, err := h.transactionService.GetWithdrawBalance(ctx, userID)
	if err != nil {
		logger.Logger().Warn("GetWithdrawAmount", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseDTO := GetBalanceResponseDTO{
		Current:   currentBalance,
		Withdrawn: withdrawBalance,
	}

	bytes, err := json.Marshal(responseDTO)
	if err != nil {
		logger.Logger().Warn("marshal GetBalanceResponseDTO", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		logger.Logger().Warn("write", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
