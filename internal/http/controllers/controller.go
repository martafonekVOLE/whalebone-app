package controllers

import (
	"simple-microservice/internal/persistence/repositories"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// Controller aggregates all application controllers.
type Controller struct {
	usersController UsersController
}

// NewController creates a new controller container instance.
func NewController(validator *validator.Validate, repository *repositories.Repository, logger *zap.Logger) *Controller {
	return &Controller{
		usersController: UsersController{
			validator:  validator,
			repository: repository.Users(),
			logger:     logger.Named("users_controller"),
		},
	}
}

// Users returns an instance of the UsersController.
func (c *Controller) Users() *UsersController {
	return &c.usersController
}
