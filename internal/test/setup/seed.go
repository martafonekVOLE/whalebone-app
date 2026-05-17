package setup

import (
	"simple-microservice/internal/persistence/models"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type userTestData struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
}

// GetUser returns a models.User populated with random (but valid) data via faker.
func GetUser(t *testing.T) models.User {
	t.Helper()

	var data userTestData

	err := faker.FakeData(&data)
	if err != nil {
		t.Fatalf("Failed to generate test data: %v", err)
	}

	return models.User{
		ExternalId:  uuid.New(),
		Name:        &data.Name,
		Email:       &data.Email,
		DateOfBirth: time.Now(),
	}
}
