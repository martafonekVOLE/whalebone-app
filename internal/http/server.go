package http

import (
	"net/http"
	"time"

	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/http/dto/responses"

	"github.com/go-chi/chi/v5"
)

const (
	defaultServerReadHeaderTimeout  = 2 * time.Second
	defaultServerIdleTimeoutTimeout = 30 * time.Second
)

func NewServer(
	hostname string,
	controller controllers.Controller,
) (server *http.Server, err error) {
	router := chi.NewRouter()
	router.NotFound(notFound)
	router.MethodNotAllowed(methodNotAllowed)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", NewUsersRouter(controller.Users()))
	})

	return &http.Server{
		Addr:              hostname,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: defaultServerReadHeaderTimeout,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       defaultServerIdleTimeoutTimeout,
		Handler:           router,
	}, nil
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	code := http.StatusNotFound
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responses.ErrorBytes(code, nil))
}

func methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	code := http.StatusMethodNotAllowed
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responses.ErrorBytes(code, nil))
}
