package handlers

import (
	"github.com/andreevym/gofermart/internal/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceHandlers struct {
	dbClient           *pgxpool.Pool
	authService        *services.AuthService
	userService        *services.UserService
	orderService       *services.OrderService
	transactionService *services.TransactionService
	newOrderCallback   func(number string)
}

func NewServiceHandlers(
	authService *services.AuthService,
	userService *services.UserService,
	orderService *services.OrderService,
	transactionService *services.TransactionService,
	newOrderCallback func(number string),
) *ServiceHandlers {
	return &ServiceHandlers{
		authService:        authService,
		userService:        userService,
		orderService:       orderService,
		transactionService: transactionService,
		newOrderCallback:   newOrderCallback,
	}
}
