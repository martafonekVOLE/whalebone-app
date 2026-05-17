package utils

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// UserResponse represents an HTTP response structure returned by user-related endpoints.
type UserResponse struct {
	Data struct {
		User struct {
			ID          uint64    `json:"id"`
			Name        string    `json:"name"`
			Email       string    `json:"email"`
			ExternalID  string    `json:"external_id"`
			DateOfBirth time.Time `json:"date_of_birth"`
		} `json:"user"`
	} `json:"data"`
}

// ErrorResponse represents a standard API error response returned by HTTP handlers.
type ErrorResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// DecodeResponse unmarshals a JSON response body into a strongly-typed struct T.
func DecodeResponse[T any](t *testing.T, body []byte) T {
	t.Helper()

	var resp T
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
