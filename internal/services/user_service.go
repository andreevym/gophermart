// services/user_service.go

package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/andreevym/gofermart/internal/repository"
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
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (us *UserService) CreateUser(ctx context.Context, username string, hashedPassword []byte) error {
	if len(hashedPassword) == 0 {
		return errors.New("password can't be empty")
	}

	password := string(hashedPassword)

	// Check if the user already exists
	_, err := us.UserRepository.GetUserByUsername(ctx, username)
	if err == nil {
		return errors.New("user already exists")
	}

	// Create a new user
	user := &repository.User{
		Username: username,
		Password: password, // For simplicity, you should hash the password before storing it.
	}

	// Save the user to the repository
	_, err = us.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("user repository create user: %w", err)
	}

	return nil
}
