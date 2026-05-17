package mappers

import (
	"fmt"
	"simple-microservice/internal/http/dto/requests"
	"simple-microservice/internal/http/dto/responses"
	"simple-microservice/internal/persistence/models"

	"github.com/google/uuid"
)

// MapUserToResponse converts User model to required response format.
func MapUserToResponse(user *models.User) responses.UserResponse {
	return responses.UserResponse{
		ID:          user.ID,
		ExternalID:  user.ExternalId.String(),
		Name:        *user.Name,
		Email:       *user.Email,
		DateOfBirth: user.DateOfBirth,
	}
}

// MapRequestToUser converts incoming request into User model.
func MapRequestToUser(request requests.CreateUserRequest) (*models.User, error) {
	parsedExternalId, err := uuid.Parse(request.ExternalId)
	if err != nil {
		return nil, fmt.Errorf("parsing uuid failed %w", err)
	}

	return &models.User{
		ExternalId:  parsedExternalId,
		Name:        &request.Name,
		Email:       &request.Email,
		DateOfBirth: request.DateOfBirth,
	}, nil
}
