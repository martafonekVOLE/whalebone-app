package integration_test

import (
	"context"
	"io"
	"net/http"
	"simple-microservice/internal/persistence/models"
	"simple-microservice/internal/persistence/repositories"
	"simple-microservice/internal/test/setup"
	"simple-microservice/internal/test/utils"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedUser(t *testing.T, repo *repositories.Repository) *models.User {
	t.Helper()

	u := setup.GetUser(t)
	err := repo.Users().CreateUser(context.Background(), &u)
	require.NoError(t, err)
	require.NotZero(t, u.ID, "expected non-zero ID after create")

	return &u
}

func TestUsersAPI_GetUserById(t *testing.T) {
	t.Parallel()

	t.Run("GetUserById_ShouldReturnUserResponse", func(t *testing.T) {
		t.Parallel()

		srv, repo := setup.NewTestServer(t)
		user := seedUser(t, repo)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/users/" + strconv.FormatUint(user.ID, 10)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body := utils.DecodeResponse[utils.UserResponse](t, readAll(t, resp.Body))
		u := body.Data.User

		assert.Equal(t, *user.Name, u.Name)
		assert.Equal(t, *user.Email, u.Email)
		assert.NotEmpty(t, u.ExternalID)
		assert.NotEmpty(t, u.DateOfBirth)
		assert.Equal(t, user.ID, u.ID)
	})

	t.Run("GetUserById_InvalidId_ShouldReturnNotFound", func(t *testing.T) {
		t.Parallel()

		srv, repo := setup.NewTestServer(t)
		user := seedUser(t, repo)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/users/" + strconv.FormatUint(user.ID+1, 10)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusNotFound), body.Error)
	})

	t.Run("GetUserById_NonNumericId_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/users/abc"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusBadRequest), body.Error)
		assert.Contains(t, body.Reason, "abc")
	})

	t.Run("GetUserById_NonExistingId_ShouldReturnNotFound", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/users/99999"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusNotFound), body.Error)
	})
}

func TestUsersAPI_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("CreateUser_ShouldReturnCreated", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		body := `{
			"name": "John Doe",
			"email": "john@example.com",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/users/save"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		responseBody := utils.DecodeResponse[utils.UserResponse](t, readAll(t, resp.Body))
		u := responseBody.Data.User

		assert.Equal(t, "John Doe", u.Name)
		assert.Equal(t, "john@example.com", u.Email)
		assert.NotEmpty(t, u.ID)
		assert.NotEmpty(t, u.ExternalID)
	})

	t.Run("CreateUser_InvalidJson_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			srv.URL+"/api/v1/users/save",
			strings.NewReader(`{bad json`),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusBadRequest), body.Error)
		assert.Contains(t, body.Reason, "invalid character")
	})

	t.Run("CreateUser_MissingFieldName_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		body := `{
			"name": "",
			"email": "john@example.com",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			srv.URL+"/api/v1/users/save",
			strings.NewReader(body),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		responseBody := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusBadRequest), responseBody.Error)
	})

	t.Run("CreateUser_InvalidFieldEmail_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		body := `{
			"name": "John Doe",
			"email": "not-an-email",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			srv.URL+"/api/v1/users/save",
			strings.NewReader(body),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		responseBody := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusBadRequest), responseBody.Error)
	})

	t.Run("CreateUser_InvalidUuid_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		body := `{
			"name": "John Doe",
			"email": "john@example.com",
			"external_id": "not-a-uuid",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			srv.URL+"/api/v1/users/save",
			strings.NewReader(body),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		responseBody := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusBadRequest), responseBody.Error)
	})
}

func TestUsersAPI_HTTPEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("API_InvalidRoute_ShouldReturnNotFound", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url := srv.URL + "/api/v1/nonexistent"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusNotFound), body.Error)
	})

	t.Run("API_InvalidRoute_ShouldReturnMethodNotAllowed", func(t *testing.T) {
		t.Parallel()

		srv, _ := setup.NewTestServer(t)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, srv.URL+"/api/v1/users/save", nil)
		require.NoError(t, err)

		resp, err := srv.Client().Do(req)
		require.NoError(t, err)

		//nolint:errcheck
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		body := utils.DecodeResponse[utils.ErrorResponse](t, readAll(t, resp.Body))
		assert.Equal(t, http.StatusText(http.StatusMethodNotAllowed), body.Error)
	})
}

func readAll(t *testing.T, r io.Reader) []byte {
	t.Helper()

	b, err := io.ReadAll(r)
	require.NoError(t, err)

	return b
}
