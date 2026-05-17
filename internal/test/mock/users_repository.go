package mock

import (
	"context"
	"simple-microservice/internal/persistence/models"
)

// MockUserRepo is a mock implementation of UsersRepository for testing.
type MockUserRepo struct {
	GetUserFunc    func(ctx context.Context, id uint64) (*models.User, error)
	CreateUserFunc func(ctx context.Context, user *models.User) error
}

// GetUser executes the mocked GetUserFunc.
func (m *MockUserRepo) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	return m.GetUserFunc(ctx, id)
}

// CreateUser executes the mocked CreateUserFunc.
func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return m.CreateUserFunc(ctx, user)
}
