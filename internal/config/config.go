// config/config.go

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

// JWTConfig represents JWT configuration.
type JWTConfig struct {
	SecretKey string `json:"secretKey" env:"JWT_SECRET_KEY"`
}

// Config represents the application configuration.
type Config struct {
	Address              string    `json:"address" env:"RUN_ADDRESS"`
	DatabaseURI          string    `json:"databaseURI" env:"DATABASE_URI"`
	AccrualSystemAddress string    `json:"accrualSystemAddress" env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTConfig            JWTConfig `json:"jwt"`
	LogLevel             string    `json:"logLevel" env:"LOG_LEVEL"`
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
	flag.StringVar(&c.JWTConfig.SecretKey, "j", "secretkey", "JWTConfig SecretKey")

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
	fmt.Println("Service Configuration:")
	fmt.Printf("Run Address: %s\n", c.Address)
	fmt.Printf("Database URI: %s\n", c.DatabaseURI)
	fmt.Printf("Accrual System Address: %s\n", c.AccrualSystemAddress)
	fmt.Printf("JWT Secret Key: %s\n", c.JWTConfig.SecretKey)
}
