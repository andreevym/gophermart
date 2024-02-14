// handlers/order_handlers.go

package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/andreevym/gofermart/internal/logger"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/mem"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type GetOrderByIDHandlerResponseDTO struct {
	// Number номер заказа
	Number string `json:"number"`
	// Status статус расчёта начисления
	Status string `json:"status"`
	// Accrual рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	Accrual *big.Int `json:"accrual,omitempty"`
}

// GetOrderByIDHandler GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
func (h *ServiceHandlers) GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, err := middleware.GetUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderNumber := chi.URLParam(r, "orderNumber")

	foundOrder, err := h.orderService.GetOrderByNumber(r.Context(), orderNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if foundOrder == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if foundOrder.Status == "" {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Error(
			"order status can't be empty",
			zap.String("orderNumber", orderNumber),
		)
		return
	}
	// order can't be with accrual when order is not processed
	if foundOrder.Status != ProcessedOrderStatus && foundOrder.Accrual != nil && foundOrder.Accrual.Sign() > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Logger().Error(
			"order can't be with accrual when order is not processed",
			zap.String("orderNumber", orderNumber),
		)
		return
	}
	respDTO := &GetOrderByIDHandlerResponseDTO{
		Number:  foundOrder.Number,
		Status:  foundOrder.Status,
		Accrual: foundOrder.Accrual,
	}
	bytes, err := json.Marshal(respDTO)
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

const (
	// RegisteredOrderStatus заказ зарегистрирован, но начисление не рассчитано;
	RegisteredOrderStatus string = "REGISTERED"
	// InvalidOrderStatus заказ не принят к расчёту, и вознаграждение не будет начислено;
	InvalidOrderStatus string = "INVALID"
	// ProcessingOrderStatus расчёт начисления в процессе;
	ProcessingOrderStatus string = "PROCESSING"
	// ProcessedOrderStatus расчёт начисления окончен;
	ProcessedOrderStatus string = "PROCESSED"
)

type GetOrdersResponseDTO struct {
	// Number номер заказа
	Number string `json:"number"`
	// Status статус расчёта начисления
	Status string `json:"status"`
	// Accrual рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
	Accrual    *big.Int `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

// GetOrdersHandler ### Взаимодействие с системой расчёта начислений баллов лояльности
//
// Для взаимодействия с системой доступен один хендлер:
//
// *   `GET /api/orders/{number}` — получение информации о расчёте начислений баллов лояльности.
//
// Формат запроса:
//
// GET /api/orders/{number} HTTP/1.1
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
// "order": "<number>",
// "status": "PROCESSED",
// "accrual": 500
// }
//
// Поля объекта ответа:
//
// *   `order` — номер заказа;
// *   `status` — статус расчёта начисления:
//
// *   `REGISTERED` — заказ зарегистрирован, но вознаграждение не рассчитано;
// *   `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
// *   `PROCESSING` — расчёт начисления в процессе;
// *   `PROCESSED` — расчёт начисления окончен;
// *   `accrual` — рассчитанные баллы к начислению, при отсутствии начисления — поле отсутствует в ответе.
//
// *   `204` — заказ не зарегистрирован в системе расчёта.
//
// *   `429` — превышено количество запросов к сервису.
//
// Формат ответа:
//
// 429 Too Many Requests HTTP/1.1
// Content-Type: text/plain
// Retry-After: 60
//
// # No more than N requests per minute allowed
//
// *   `500` — внутренняя ошибка сервера.
//
// Заказ может быть взят в расчёт в любой момент после его совершения. Время выполнения расчёта системой не регламентировано. Статусы `INVALID` и `PROCESSED` являются окончательными.
//
// Общее количество запросов информации о начислении не ограничено.
func (h *ServiceHandlers) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	foundOrders, err := h.orderService.OrderRepository.GetOrdersByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respDTOs := make([]GetOrdersResponseDTO, 0)
	for _, foundOrder := range foundOrders {
		resp := GetOrdersResponseDTO{
			Number:     foundOrder.Number,
			Status:     foundOrder.Status,
			Accrual:    foundOrder.Accrual,
			UploadedAt: foundOrder.UploadedAt.Format(time.RFC3339),
		}

		respDTOs = append(respDTOs, resp)
	}
	bytes, err := json.Marshal(respDTOs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(respDTOs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PostOrdersHandler загрузка номера заказа
// #### **Загрузка номера заказа**
//
// Хендлер: `POST /api/user/orders`.
//
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
//
// Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна).
//
// Формат запроса:
//
// POST /api/user/orders HTTP/1.1
// Content-Type: text/plain
// ...
//
// 12345678903
//
// Возможные коды ответа:
//
// *   `200` — номер заказа уже был загружен этим пользователем;
// *   `202` — новый номер заказа принят в обработку;
// *   `400` — неверный формат запроса;
// *   `401` — пользователь не аутентифицирован;
// *   `409` — номер заказа уже был загружен другим пользователем;
// *   `422` — неверный формат номера заказа;
// *   `500` — внутренняя ошибка сервера.
func (h *ServiceHandlers) PostOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNumber := string(bytes)

	err = goluhn.Validate(orderNumber)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	existsOrder, err := h.orderService.GetOrderByNumber(r.Context(), orderNumber)
	if err != nil && !errors.Is(err, mem.ErrOrderNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if existsOrder != nil {
		if existsOrder.UserID == userID {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusConflict)
		return
	}

	newOrder := &repository.Order{
		Number: orderNumber,
		UserID: userID,
		Status: RegisteredOrderStatus,
	}
	_, err = h.orderService.OrderRepository.CreateOrder(r.Context(), newOrder)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
