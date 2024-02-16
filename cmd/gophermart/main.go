package main

import (
	"context"
	"log"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/config"
	"github.com/andreevym/gofermart/internal/handlers"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/andreevym/gofermart/internal/server"
	"github.com/andreevym/gofermart/internal/services"
	"github.com/andreevym/gofermart/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Create a new configuration instance
	cfg := config.NewConfig()

	// Parse the configuration from flags and environment variables
	if err := cfg.Parse(); err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}

	// Print the configuration
	cfg.Print()

	ctx := context.Background()

	db, err := pgxpool.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Apply database migrations
	err = postgres.Migration(ctx, db)
	if err != nil {
		log.Fatalf("Failed to apply database migrations: %v", err)
	}

	// Create repositories
	transactionRepository := postgres.NewTransactionRepository(db)
	userRepository := postgres.NewUserRepository(db)
	orderRepository := postgres.NewOrderRepository(db)

	// Create services
	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)
	userService := services.NewUserService(userRepository)
	transactionService := services.NewTransactionService(transactionRepository)
	orderService := services.NewOrderService(transactionService, orderRepository, accrualService)

	jwtConfig := config.JWTConfig{}
	authService := services.NewAuthService(userService, jwtConfig)

	newOrderNumbersCh := make(chan string)
	defer close(newOrderNumbersCh)

	// запуск отдельного процесса для процессинга заявок, только если при запуске сервиса был передан адрес accrualService
	if accrualService != nil {
		go func() {
			for orderNumber := range newOrderNumbersCh {
				err := orderService.RetryOrderProcessing(orderNumber)
				if err != nil {
					logger.Logger().Error("RetryOrderProcessing", zap.Error(err))
					panic(err.Error())
				}
				return
			}
		}()
	}

	// объявляем функцию, которая будет вызвана при создании заявки
	newOrderCallback := func(number string) {
		if accrualService != nil {
			newOrderNumbersCh <- number
		}
	}
	// объявляем все сервисы в одной структуре т.к так удобнее изменять кол-во сервисов
	// которые мы будем использовать в обработчике
	serviceHandlers := handlers.NewServiceHandlers(
		authService,
		userService,
		orderService,
		transactionService,
		newOrderCallback,
	)
	router := handlers.NewRouter(
		serviceHandlers,
		middleware.NewAuthMiddleware(authService).WithAuthentication,
		middleware.RequestLoggerMiddleware,
	)

	server := server.NewServer(router)
	if server == nil {
		log.Fatalf("Server can't be nil: %v", err)
	}
	server.Run(cfg.Address)
	server.Shutdown()
}
