package mappers

import (
	"fmt"
	"github.com/google/uuid"

	"simple-microservice/internal/http/dto/requests"
	"simple-microservice/internal/http/dto/responses"
	"simple-microservice/internal/persistence/models"
)

func MapUserToResponse(user *models.User) responses.UserResponse {
	return responses.UserResponse{
		ID:          user.ID,
		ExternalID:  user.ExternalId.String(),
		Name:        user.FirstName + " " + user.LastName,
		Email:       user.Email,
		DateOfBirth: user.DateOfBirth.String(),
	}
}

func MapRequestToUser(request requests.CreateUserRequest) (*models.User, error) {
	parsedExternalId, err := uuid.Parse(request.ExternalId)
	if err != nil {
		return nil, fmt.Errorf("parsing uuid failed %w", err)
	}

	return &models.User{
		ExternalId:  parsedExternalId,
		Email:       request.Email,
		DateOfBirth: request.DateOfBirth,
	}, nil
}
