package repositories

import (
	"context"
	"simple-microservice/internal/persistence/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UsersRepository handles persistence of User entities.
type UsersRepository struct {
	DB     *gorm.DB
	logger *zap.Logger
}

// GetUser retrieves a user by its ID.
func (u *UsersRepository) GetUser(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User

	result := u.DB.WithContext(ctx).First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// CreateUser saves a new User record in the database.
func (u *UsersRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := u.DB.WithContext(ctx).Create(user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
