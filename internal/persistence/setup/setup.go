package setup

import (
	"fmt"
	"time"

	"simple-microservice/internal/persistence/config"
	"simple-microservice/internal/persistence/models"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

const (
	DefaultMaxIdleConnections = 10
	DefaultMaxOpenConnections = 100
	DefaultConnMaxLifetime    = time.Hour
)

// ConnectDatabase connects to a database.
func ConnectDatabase(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB %w", err)
	}

	if config.MaxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(*config.MaxOpenConnections)
	} else {
		sqlDB.SetMaxIdleConns(DefaultMaxIdleConnections)
	}

	if config.MaxOpenConnections != nil {
		sqlDB.SetMaxOpenConns(*config.MaxOpenConnections)
	} else {
		sqlDB.SetMaxOpenConns(DefaultMaxOpenConnections)
	}

	if config.MaxLifeTimeSeconds != nil {
		sqlDB.SetConnMaxLifetime(time.Duration(*config.MaxLifeTimeSeconds) * time.Second)
	} else {
		sqlDB.SetConnMaxLifetime(DefaultConnMaxLifetime)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	return db, nil
}
