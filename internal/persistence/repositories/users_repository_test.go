package repositories_test

import (
	"context"
	"errors"
	"simple-microservice/internal/logging"
	"simple-microservice/internal/persistence/models"
	"simple-microservice/internal/persistence/repositories"
	"simple-microservice/internal/test/setup"
	"testing"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

func TestUsersRepository(t *testing.T) {
	t.Parallel()

	logger := logging.ConfigureTestLogger()

	t.Run("CreateUser_ShouldCreate", func(t *testing.T) {
		t.Parallel()

		db := setup.NewTestDatabase(t)
		repo := repositories.NewRepository(db, logger)

		user := setup.GetUser(t)

		err := repo.Users().CreateUser(context.Background(), &user)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		if user.ID == 0 {
			t.Error("expected non-zero ID after create")
		}
	})

	t.Run("CreateUser_EmptyFields", func(t *testing.T) {
		t.Parallel()

		db := setup.NewTestDatabase(t)
		repo := repositories.NewRepository(db, logger)

		user := models.User{
			ExternalId: uuid.New(),
		}

		err := repo.Users().CreateUser(context.Background(), &user)
		if err == nil {
			t.Fatal("expected CreateUser to fail due to empty fields, but got no error")
		}
	})

	t.Run("GetUser_ShouldSucceed", func(t *testing.T) {
		t.Parallel()

		db := setup.NewTestDatabase(t)
		repo := repositories.NewRepository(db, logger)

		user := setup.GetUser(t)

		err := repo.Users().CreateUser(context.Background(), &user)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		if user.ID == 0 {
			t.Fatal("expected non-zero ID after create")
		}

		got, err := repo.Users().GetUser(context.Background(), user.ID)
		if err != nil {
			t.Fatalf("GetUser failed: %v", err)
		}

		if got.ID != user.ID {
			t.Errorf("ID mismatch: received = %d, wanted = %d", got.ID, user.ID)
		}
	})

	t.Run("GetUser_NonExistentId_ShouldReturnErrRecordNotFound", func(t *testing.T) {
		t.Parallel()

		db := setup.NewTestDatabase(t)
		repo := repositories.NewRepository(db, logger)

		_, err := repo.Users().GetUser(context.Background(), 99999)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("GetUser_ZeroId_ShouldReturnErrRecordNotFound", func(t *testing.T) {
		t.Parallel()

		db := setup.NewTestDatabase(t)
		repo := repositories.NewRepository(db, logger)

		_, err := repo.Users().GetUser(context.Background(), 0)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}
