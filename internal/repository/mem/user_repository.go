package mem

import (
	"context"
	"errors"
	"math/rand"

	"github.com/andreevym/gofermart/internal/repository"
)

// MemUserRepository is an implementation of UserRepository using an in-memory store (for demonstration purposes)
type MemUserRepository struct {
	users map[int64]*repository.User
}

// NewMemUserRepository creates a new instance of MemUserRepository
func NewMemUserRepository() *MemUserRepository {
	return &MemUserRepository{
		users: make(map[int64]*repository.User),
	}
}

func (r *MemUserRepository) CreateUser(_ context.Context, user *repository.User) (*repository.User, error) {
	// Check if the user already exists
	if _, ok := r.users[user.ID]; ok {
		return nil, errors.New("user already exists")
	}
	// Assign an ID to the user (for demonstration purposes; in a real scenario, you'd use a proper ID generation mechanism)
	user.ID = generateUserID()
	r.users[user.ID] = user
	return user, nil
}

func (r *MemUserRepository) GetUserByID(_ context.Context, userID int64) (*repository.User, error) {
	user, ok := r.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *MemUserRepository) GetUserByUsername(_ context.Context, username string) (*repository.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MemUserRepository) UpdateUser(_ context.Context, user *repository.User) (*repository.User, error) {
	// Check if the user exists
	if _, ok := r.users[user.ID]; !ok {
		return nil, errors.New("user not found")
	}
	// Update the user
	r.users[user.ID] = user
	return user, nil
}

func (r *MemUserRepository) DeleteUser(_ context.Context, userID int64) error {
	// Check if the user exists
	if _, ok := r.users[userID]; !ok {
		return errors.New("user not found")
	}
	// Delete the user
	delete(r.users, userID)
	return nil
}

// generateUserID generates a unique user ID (for demonstration purposes)
func generateUserID() int64 {
	return rand.Int63()
}
