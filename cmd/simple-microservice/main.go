package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	nethttp "net/http"
	"os"
	"os/signal"
	"simple-microservice/internal/http"
	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/logging"
	"simple-microservice/internal/persistence/config"
	"simple-microservice/internal/persistence/repositories"
	"simple-microservice/internal/persistence/setup"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const DefaultShutdownTimeout = 10 * time.Second

func main() {
	logger := logging.ConfigureLogger()
	defer func(zapLogger *zap.Logger) {
		_ = zapLogger.Sync()
	}(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Info("no .env file found, falling back to system environment variables")
	}

	err = runApp(logger)
	if err != nil {
		logger.Fatal("application failed to start", zap.Error(err))
	}
}

// runApp serves as an application's entrypoint.
func runApp(logger *zap.Logger) error {
	dbConfig := config.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
	}

	db, err := setup.ConnectDatabase(dbConfig)
	if err != nil {
		logger.Error("failed to connect database", zap.Error(err))
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get underlying sql.DB", zap.Error(err))
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			logger.Error("failed to close underlying sql.DB", zap.Error(err))
		}
	}()

	validate := validator.New(validator.WithRequiredStructEnabled())
	repository := repositories.NewRepository(db, logger)
	controller := controllers.NewController(validate, repository, logger)

	hostname := net.JoinHostPort(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	server, err := http.NewServer(hostname, *controller, logger)
	if err != nil {
		logger.Error("failed to start http server", zap.Error(err))
		return fmt.Errorf("failed to start http server: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting http server", zap.String("hostname", hostname))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server failed unexpectedly: %w", err)

	case <-ctx.Done():
		logger.Info("graceful shutdown initiated")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("server shutdown failed", zap.Error(err))
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Info("server exited properly")

	return nil
}
