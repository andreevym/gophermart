package main

import (
	"context"
	"log"

	"github.com/andreevym/gophermart/internal/accrual"
	"github.com/andreevym/gophermart/internal/config"
	"github.com/andreevym/gophermart/internal/handlers"
	"github.com/andreevym/gophermart/internal/middleware"
	"github.com/andreevym/gophermart/internal/repository/postgres"
	"github.com/andreevym/gophermart/internal/scheduler"
	"github.com/andreevym/gophermart/internal/server"
	"github.com/andreevym/gophermart/internal/services"
	"github.com/jackc/pgx/v4/pgxpool"
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
	authService := services.NewAuthService(userService, cfg.JWTSecretKey)

	// запуск отдельного процесса для процессинга заявок, только если при запуске сервиса был передан адрес accrualService
	if accrualService != nil {
		accrualScheduler := scheduler.NewAccrualScheduler(accrualService, orderService, cfg.PollOrdersDelay, cfg.MaxOrderAttempts)
		defer accrualScheduler.Shutdown()
	}

	// объявляем все сервисы в одной структуре т.к так удобнее изменять кол-во сервисов
	// которые мы будем использовать в обработчике
	serviceHandlers := handlers.NewServiceHandlers(
		authService,
		userService,
		orderService,
		transactionService,
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
	defer server.Shutdown()
	server.Run(cfg.Address)
}
