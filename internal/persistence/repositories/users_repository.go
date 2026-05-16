package repositories

import (
	"simple-microservice/internal/persistence/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UsersRepository struct {
	DB     *gorm.DB
	logger *zap.Logger
}

func (u *UsersRepository) GetUser(id uint64) (*models.User, error) {
	var user models.User

	result := u.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *UsersRepository) CreateUser(user *models.User) error {
	result := u.DB.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
