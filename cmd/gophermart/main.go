package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andreevym/gofermart/internal/accrual"
	"github.com/andreevym/gofermart/internal/config"
	"github.com/andreevym/gofermart/internal/handlers"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository/postgres"
	"github.com/andreevym/gofermart/internal/server"
	"github.com/andreevym/gofermart/internal/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Create a new configuration instance
	cfg := config.NewConfig()

	// Parse the configuration from flags and environment variables
	if err := cfg.Parse(); err != nil {
		fmt.Printf("Error parsing configuration: %s\n", err)
		return
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

	accrualService := accrual.NewAccrualService(cfg.AccrualSystemAddress)

	// Create services and repositories
	userService := services.NewUserService(postgres.NewUserRepository(db))
	orderService := services.NewOrderService(postgres.NewOrderRepository(db), accrualService)

	transactionRepository := postgres.NewTransactionRepository(db)

	transactionService := services.NewTransactionService(transactionRepository)

	jwtConfig := config.JWTConfig{}
	authService := services.NewAuthService(userService, jwtConfig)

	serviceHandlers := handlers.NewServiceHandlers(
		authService,
		userService,
		orderService,
		transactionService,
	)

	// Create router with tracer
	router := handlers.NewRouter(
		serviceHandlers,
		middleware.NewAuthMiddleware(authService).WithAuthentication,
		middleware.RequestLoggerMiddleware,
	)

	// Create server
	server := server.NewServer(router)
	if server == nil {
		panic("server can't be nil")
	}

	// Run server
	server.Run(cfg.Address)
	server.Shutdown()
}
