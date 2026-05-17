package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Response is a generic API response wrapper used by all HTTP handlers.
type Response struct {
	Data             map[string]any    `json:"data,omitempty"`
	Error            string            `json:"error,omitempty"`
	Reason           string            `json:"reason,omitempty"`
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

// ErrorBytes returns a JSON-encoded API error response with the given HTTP status code.
func ErrorBytes(httpCode int, err error) []byte {
	response := Response{
		Error: http.StatusText(httpCode),
	}

	if err != nil {
		response.Reason = err.Error()
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return []byte(`{"error":"Internal Server Error","reason":"Failed to marshal error response"}`)
	}

	return bytes
}

// ValidationErrorBytes returns a JSON-encoded API validation error response with detail of invalid fields.
func ValidationErrorBytes(httpCode int, err error) []byte {
	response := Response{
		Error:  http.StatusText(httpCode),
		Reason: "validation failed",
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		errsMap := make(map[string]string)

		for _, e := range validationErrs {
			errsMap[e.Field()] = fmt.Sprintf("invalid value for '%s'", e.Tag())
		}

		response.ValidationErrors = errsMap
	} else {
		response.Reason = err.Error()
	}

	bytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return []byte(`{"error":"Internal Server Error"}`)
	}

	return bytes
}

// DataBytes returns a valid JSON-encoded API response.
func DataBytes(name string, data any) (bytes []byte, err error) {
	response := Response{
		Data: map[string]any{
			name: data,
		},
	}

	bytes, err = json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal data: %w", err)
	}

	return bytes, nil
}
