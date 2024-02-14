package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/andreevym/gofermart/internal/logger"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"go.uber.org/zap"
)

type UserBalanceWithdrawRequestDTO struct {
	OrderNumber string `json:"order"` // номер заказа
	Sum         int64  `json:"sum"`   // сумма баллов к списанию в счёт оплаты
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

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.transactionService.GetTransactionsByUserAndOperationType(userID, repository.WithdrawOperationType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	bytes, err := json.Marshal(withdrawals)
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

//
//type Withdrawal struct {
//	Order       string    `json:"order"`
//	Sum         *big.Int  `json:"sum"`
//	ProcessedAt time.Time `json:"processed_at"`
//}
//
//func (h ServiceHandlers) getWithdrawalsTransaction(userID int64) ([]Withdrawal, error) {
//	withdrawalTransactions, err := h.transactionService.GetTransactionsByUserAndOperationType(userID, mem.WithdrawOperationType)
//	if err != nil {
//		return nil, err
//	}
//
//	withdrawals := []Withdrawal
//	for _, transaction := range withdrawalTransactions {
//		w := Withdrawal{
//			Order:       transaction.Reason,
//			Sum:         transaction.Amount,
//			ProcessedAt: transaction.Date,
//		}
//		withdrawals = append(withdrawals, w)
//	}
//	return withdrawals, nil
//}

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
		logger.Logger().Warn("middleware.GetUserID", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var userBalanceWithdrawDTO UserBalanceWithdrawRequestDTO
	err = json.Unmarshal(bytes, &userBalanceWithdrawDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByNumber(ctx, userBalanceWithdrawDTO.OrderNumber)
	if err != nil {
		if errors.Is(err, postgres.ErrOrderNotFound) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.transactionService.Withdraw(
		context.Background(),
		userID,
		big.NewInt(userBalanceWithdrawDTO.Sum),
		strconv.FormatInt(order.ID, 10),
	)
	if err != nil {
		logger.Logger().Error("transactionService.Withdraw", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type BalanceDTO struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type GetBalanceResponseDTO struct {
	Current   *big.Int `json:"current"`
	Withdrawn *big.Int `json:"withdrawn"`
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

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		logger.Logger().Warn("middleware.GetUserID", zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userAccountBalance, err := h.userAccountService.GetCurrentBalance(r.Context(), userID)
	if err != nil {
		logger.Logger().Warn("GetCurrentBalance", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawAmount, err := h.userAccountService.GetWithdrawAmount(r.Context(), userID)
	if err != nil {
		logger.Logger().Warn("GetWithdrawAmount", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseDTO := GetBalanceResponseDTO{
		Current:   userAccountBalance,
		Withdrawn: withdrawAmount,
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
