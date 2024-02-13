package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreevym/gofermart/internal/config"
	"github.com/golang-jwt/jwt"
)

// AuthService represents a concrete implementation of the AuthService interface.
type AuthService struct {
	userService *UserService
	jwtConfig   config.JWTConfig
}

var (
	ErrAuthAlreadyExists         = errors.New("user already exists")
	ErrAuthWrongLoginAndPassword = errors.New("invalid username or password")
)

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userService *UserService, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{
		userService: userService,
		jwtConfig:   jwtConfig,
	}
}

// Login authenticates a user and returns a JWT token.
func (a *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := a.userService.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("UserRepository.GetUserByUsername: %w", err)
	}
	if user == nil || !user.IsValidPassword(password) {
		return "", ErrAuthWrongLoginAndPassword
	}

	token, err := a.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("GenerateToken: %w", err)
	}
	return token, nil
}

// Register registers a new user.
func (a *AuthService) Register(ctx context.Context, username string, password string) error {
	_, err := a.userService.UserRepository.GetUserByUsername(ctx, username)
	if err == nil {
		return ErrAuthAlreadyExists
	}

	err = a.userService.CreateUser(ctx, username, password)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// Logout invalidates a JWT token.
func (a *AuthService) Logout(tokenString string) error {
	// You can implement logout logic here, such as blacklisting the token
	// or marking it as expired. For simplicity, let's assume we don't need
	// to perform any action for logout in this example.
	return nil
}

// GenerateToken generates a JWT token for the given user ID.
func (a *AuthService) GenerateToken(userID int64) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(a.jwtConfig.SecretKey)
	if err != nil {
		return "", fmt.Errorf("sign the token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and extracts the user ID.
func (a *AuthService) ValidateToken(tokenString string) (int64, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return a.jwtConfig.SecretKey, nil
	})
	if err != nil {
		return -1, fmt.Errorf("jwt parse: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return -1, errors.New("token is not valid")
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, errors.New("invalid token claims")
	}
	userID, ok := claims["userID"].(int64)
	if !ok {
		return -1, errors.New("invalid user ID in token")
	}

	return userID, nil
}
