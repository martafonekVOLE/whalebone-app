package setup

import (
	"net/http/httptest"
	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/logging"
	"simple-microservice/internal/persistence/repositories"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	customHttp "simple-microservice/internal/http"
)

// NewTestServer prepares a new server for testing.
func NewTestServer(t *testing.T) (*httptest.Server, *repositories.Repository) {
	t.Helper()

	// 1. Load .env
	err := godotenv.Load()
	if err != nil {
		t.Log("no .env file found, falling back to system environment variables")
	}

	// 2. Prepare DB container
	db := NewTestDatabase(t)

	logger := logging.ConfigureTestLogger()
	validate := validator.New(validator.WithRequiredStructEnabled())
	repository := repositories.NewRepository(db, logger)
	controller := controllers.NewController(validate, repository, logger)

	router := chi.NewRouter()
	router.NotFound(customHttp.NotFound)
	router.MethodNotAllowed(customHttp.MethodNotAllowed)
	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", customHttp.NewUsersRouter(controller.Users()))
	})

	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	return srv, repository
}
