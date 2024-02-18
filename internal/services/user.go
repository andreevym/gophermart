// services/user_service.go

package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreevym/gophermart/internal/repository"
)

var (
	ErrUserPasswordEmpty   = errors.New("password can't be empty")
	ErrUserPasswordInvalid = errors.New("invalid password")
	ErrUserAlreadyExists   = errors.New("user already exists")
)

// UserService struct represents the service for users
type UserService struct {
	UserRepository repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

// AuthenticateUser authenticates a user.
func (us *UserService) AuthenticateUser(ctx context.Context, username string, password string) (*repository.User, error) {
	// Retrieve the user by username
	user, err := us.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user repository get user by username %w", err)
	}

	// Verify the password
	if user.Password != password { // For simplicity, you should compare hashed passwords here.
		return nil, ErrUserPasswordInvalid
	}

	return user, nil
}

func (us *UserService) CreateUser(ctx context.Context, username string, password string) error {
	if len(password) == 0 {
		return ErrUserPasswordEmpty
	}

	// Check if the user already exists
	_, err := us.UserRepository.GetUserByUsername(ctx, username)
	if err == nil {
		return ErrUserAlreadyExists
	}

	// Create a new user
	newUser := repository.User{
		Username: username,
		Password: password, // For simplicity, you should hash the password before storing it.
	}

	// Save the user to the repository
	err = us.UserRepository.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("user repository create user: %w", err)
	}

	return nil
}
