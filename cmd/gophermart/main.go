package main

import (
	"fmt"
	"log"

	"github.com/andreevym/gofermart/internal/config"
	"github.com/andreevym/gofermart/internal/database"
	"github.com/andreevym/gofermart/internal/handlers"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository/mem"
	"github.com/andreevym/gofermart/internal/server"
	"github.com/andreevym/gofermart/internal/services"
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

	// Initialize database connection
	db, err := database.Connect(cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

	// Apply database migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to apply database migrations: %v", err)
	}

	// Create services and repositories
	userService := services.NewUserService(mem.NewMemUserRepository())
	orderService := services.NewOrderService(mem.NewMemOrderRepository())
	userAccountRepository := mem.NewMemUserAccountRepository()
	transactionService := services.NewTransactionService(mem.NewMemTransactionRepository(), userAccountRepository)
	userAccountService := services.NewUserAccountService(userAccountRepository)

	jwtConfig := config.JWTConfig{}
	authService := services.NewAuthService(userService, jwtConfig)

	serviceHandlers := handlers.NewServiceHandlers(
		authService,
		userService,
		orderService,
		transactionService,
		userAccountService,
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
	server.Run(":8080")
	server.Shutdown()
}
