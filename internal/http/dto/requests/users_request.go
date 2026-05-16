package requests

import "time"

type CreateUserRequest struct {
	Name        string    `json:"name"          validate:"required"`
	Email       string    `json:"email"         validate:"required,email"`
	ExternalId  string    `json:"external_id"   validate:"required,uuid"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}
