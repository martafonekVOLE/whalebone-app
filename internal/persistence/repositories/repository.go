package repositories

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Repository contains all service repositories.
type Repository struct {
	usersRepository UsersRepository
}

// NewRepository is a repository for Repository.
func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		usersRepository: UsersRepository{
			DB:     db,
			logger: logger.Named("users_repository"),
		},
	}
}

// Users returns a repository for user model.
func (repo *Repository) Users() *UsersRepository {
	return &repo.usersRepository
}
