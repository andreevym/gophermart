// internal/router/router.go

package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter creates a new HTTP router with the specified handlers and tracer.
func NewRouter(s *ServiceHandlers, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Attach other middlewares
	//r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middlewares...)

	//POST /api/user/register — регистрация пользователя;
	r.Post("/api/user/register", s.PostRegisterUser)
	//POST /api/user/login — аутентификация пользователя;
	r.Post("/api/user/login", s.PostLoginUser)
	//POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	r.Post("/api/user/orders", s.PostOrdersHandler)
	//GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	r.Get("/api/user/orders", s.GetOrdersHandler)
	//GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
	r.Get("/api/user/balance", s.GetBalanceHandler)
	//POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	r.Post("/api/user/balance/withdraw", s.PostWithdrawHandler)
	//GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	r.Get("/api/user/withdrawals", s.GetWithdrawalsHandler)
	r.Get("/api/ping", s.GetPingHandler)
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
	})
	return r
}
