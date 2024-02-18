package services

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
)

// AuthService represents a concrete implementation of the AuthService interface.
type AuthService struct {
	userService  *UserService
	jwtSecretKey string
}

var (
	ErrAuthBadCredentials = errors.New("username or password is incorrect")
)

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userService *UserService, jwtSecretKey string) *AuthService {
	return &AuthService{
		userService:  userService,
		jwtSecretKey: jwtSecretKey,
	}
}

// Login authenticates a user and returns a JWT token.
func (a *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := a.userService.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("UserRepository.GetUserByUsername: %w", err)
	}
	if user == nil || !user.IsValidPassword(password) {
		return "", ErrAuthBadCredentials
	}

	token, err := a.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("GenerateToken: %w", err)
	}
	return token, nil
}

// Register registers a new user.
func (a *AuthService) Register(ctx context.Context, username string, password string) (string, error) {
	err := a.userService.CreateUser(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	return a.Login(ctx, username, password)
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
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.MapClaims{
		"userID": strconv.FormatInt(userID, 10),
	})

	var jwtSecretKey *ecdsa.PrivateKey
	if a.jwtSecretKey == "" {
		jwtSecretKey = GenPrivateKeyMust()
		privateKey, err := x509.MarshalECPrivateKey(jwtSecretKey)
		if err != nil {
			return "", fmt.Errorf("x509.MarshalECPrivateKey: %w", err)
		}
		a.jwtSecretKey = string(privateKey)
	} else {
		var err error
		jwtSecretKey, err = x509.ParseECPrivateKey([]byte(a.jwtSecretKey))
		if err != nil {
			return "", fmt.Errorf("x509.MarshalECPrivateKey: %w", err)
		}
	}
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("sign the token: %w", err)
	}
	return t, nil
}

// ValidateToken validates a JWT token and extracts the user ID.
func (a *AuthService) ValidateToken(tokenString string) (int64, error) {
	privateKey, err := x509.ParseECPrivateKey([]byte(a.jwtSecretKey))
	if err != nil {
		return -1, fmt.Errorf("x509.ParseECPrivateKey: %w", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		return -1, fmt.Errorf("jwt parse: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return -1, errors.New("token is not valid")
	}

	// Extract the user ID from the token claims
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return -1, errors.New("invalid token claims")
	}
	id := mapClaims["userID"].(string)
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("strconv.ParseInt, %s: invalid user ID in token: %w", id, err)
	}

	return userID, nil
}

func GenPrivateKeyMust() *ecdsa.PrivateKey {
	key, err := TestGenKey()
	if err != nil {
		panic(err)
	}
	return key
}

func TestGenKey() (*ecdsa.PrivateKey, error) {
	// Генерируем новый ECDSA ключ
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Ошибка при генерации ключа:", err)
		return nil, err
	}

	return privateKey, nil
}

func MarshalPrivateKey(privateKey *ecdsa.PrivateKey) error {
	// Кодируем приватный ключ в формат PEM
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Println("Ошибка при кодировании приватного ключа:", err)
		return err
	}

	privateKeyPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	file := os.Stdout
	err = pem.Encode(file, privateKeyPEM)
	if err != nil {
		fmt.Println("Ошибка при кодировании PEM-блока:", err)
		return err
	}
	fmt.Println("Приватный ключ сохранен в файл private.key")

	return nil
}
