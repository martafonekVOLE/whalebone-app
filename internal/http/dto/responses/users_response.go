package responses

import "time"

// UserResponse represents an outgoing response.
type UserResponse struct {
	ID          uint64    `json:"id"`
	ExternalID  string    `json:"external_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}
