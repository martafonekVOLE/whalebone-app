package main

import (
	"fmt"
	"net"
	"os"

	"simple-microservice/internal/http"
	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/logging"
	"simple-microservice/internal/persistence/config"
	"simple-microservice/internal/persistence/repositories"
	"simple-microservice/internal/persistence/setup"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

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

	validate := validator.New(validator.WithRequiredStructEnabled())
	repository := repositories.NewRepository(db, logger)
	controller := controllers.NewController(validate, repository, logger)

	hostname := net.JoinHostPort(os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	server, err := http.NewServer(hostname, *controller)
	if err != nil {
		logger.Error("failed to start http server", zap.Error(err))
		return fmt.Errorf("failed to start http server: %w", err)
	}

	logger.Info("starting http server", zap.String("hostname", hostname))

	err = server.ListenAndServe()
	if err != nil {
		logger.Error("failed running http server", zap.Error(err))
		return fmt.Errorf("failed running http server: %w", err)
	}

	return nil
}
