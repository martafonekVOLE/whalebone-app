package setup

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"simple-microservice/internal/persistence/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const StartupTimeout = 15 * time.Second

// NewTestDatabase starts a testing DB in a disposable PostgreSQL container.
func NewTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "postgres:15-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForSQL("5432/tcp", "pgx", func(host, port string) string {
			portNum := port
			if before, _, ok := strings.Cut(port, "/"); ok {
				portNum = before
			}

			return fmt.Sprintf("host=%s port=%s user=test password=test dbname=test sslmode=disable", host, portNum)
		}).
			WithStartupTimeout(StartupTimeout).
			WithQuery("SELECT 10"),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		if err != nil {
			t.Logf("WARNING: failed to terminate test container: %v", err)
		}
	})

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@localhost:%s/test?sslmode=disable", mappedPort.Port())

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}

	t.Cleanup(func() {
		db, err := gormDB.DB()
		if err != nil {
			t.Fatalf("Failed to get underlying sql.DB from GORM: %v", err)
		}

		err = db.Close()
		if err != nil {
			t.Fatalf("Failed to close GORM database connection: %v", err)
		}
	})

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return gormDB
}
