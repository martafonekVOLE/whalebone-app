package http

import (
	"net/http"
	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/http/dto/responses"
	"simple-microservice/internal/http/middleware"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	defaultServerTimeout           = 5 * time.Second
	defaultServerReadHeaderTimeout = 2 * time.Second
	defaultServerIdleTimeout       = 30 * time.Second
)

// NewServer creates and configures an HTTP server.
func NewServer(
	hostname string,
	controller controllers.Controller,
	logger *zap.Logger,
) (server *http.Server, err error) {
	router := chi.NewRouter()
	router.NotFound(NotFound)
	router.MethodNotAllowed(MethodNotAllowed)

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(chimiddleware.RequestID)
		r.Use(middleware.Logger(logger))
		r.Use(chimiddleware.Recoverer)
		r.Use(chimiddleware.Timeout(defaultServerTimeout))

		r.Use(middleware.ContentType)

		r.Mount("/users", NewUsersRouter(controller.Users()))
	})

	return &http.Server{
		Addr:              hostname,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: defaultServerReadHeaderTimeout,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       defaultServerIdleTimeout,
		Handler:           router,
	}, nil
}

// NotFound handles requests to undefined routes.
func NotFound(w http.ResponseWriter, _ *http.Request) {
	code := http.StatusNotFound
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responses.ErrorBytes(code, nil))
}

// MethodNotAllowed handles requests with an unsupported HTTP method.
func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	code := http.StatusMethodNotAllowed
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responses.ErrorBytes(code, nil))
}
