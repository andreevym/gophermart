// config/config.go

package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/andreevym/gophermart/pkg/logger"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

// Config represents the application configuration.
type Config struct {
	Address              string        `json:"address" env:"RUN_ADDRESS"`
	DatabaseURI          string        `json:"databaseURI" env:"DATABASE_URI"`
	AccrualSystemAddress string        `json:"accrualSystemAddress" env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string        `json:"logLevel" env:"LOG_LEVEL"`
	PollOrdersDelay      time.Duration `json:"pollDuration" env:"POLL_ORDERS_DURATION"`
	MaxOrderAttempts     int           `json:"maxOrderAttempts" env:"MAX_ORDER_ATTEMPTS"`
	JWTSecretKey         string        `json:"secretKey" env:"JWT_SECRET_KEY"`
}

// NewConfig creates a new Config instance with default values.
func NewConfig() *Config {
	return &Config{}
}

// Parse parses the configuration from flags and environment variables.
func (c *Config) Parse() error {
	// Define usage message for the flag
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Define flags
	flag.StringVar(&c.Address, "a", "", "Run Address (overrides environment variable)")
	flag.StringVar(&c.DatabaseURI, "d", "", "Database URI (overrides environment variable)")
	flag.StringVar(&c.AccrualSystemAddress, "r", "", "Accrual System Address (overrides environment variable)")
	flag.StringVar(&c.LogLevel, "l", "info", "Logging level [INFO, DEBUG, ERROR]")
	flag.IntVar(&c.MaxOrderAttempts, "maxOrderAttempts", 3, "Logging level [INFO, DEBUG, ERROR]")
	flag.DurationVar(&c.PollOrdersDelay, "pollOrdersDuration", 10*time.Millisecond, "duration for handle orders")
	flag.StringVar(&c.JWTSecretKey, "j", "", "JWTConfig SecretKey")

	// Parse flags
	flag.Parse()

	// Parse environment variables
	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}

// Print prints the configuration to stdout.
func (c *Config) Print() {
	logger.Logger().Info(
		"Service Configuration",
		zap.String("Run Address", c.Address),
		zap.String("Database URI", c.DatabaseURI),
		zap.String("Accrual System Address", c.AccrualSystemAddress),
		zap.String("PollOrdersDelay", c.PollOrdersDelay.String()),
		zap.Int("MaxOrderAttempts", c.MaxOrderAttempts),
		zap.String("LogLevel", c.LogLevel),
		zap.String("JWT Secret Key", c.JWTSecretKey),
	)
}
