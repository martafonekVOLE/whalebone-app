package http

import (
	"simple-microservice/internal/http/controllers"

	"github.com/go-chi/chi/v5"
)

// NewUsersRouter creates a chi router with registered user-related HTTP endpoints.
func NewUsersRouter(c *controllers.UsersController) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/save", c.CreateUser)
	router.Get("/{id}", c.GetUserById)

	return router
}
