package handlers

import "github.com/andreevym/gofermart/internal/services"

type ServiceHandlers struct {
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
