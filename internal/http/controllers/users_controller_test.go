package controllers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"simple-microservice/internal/http/controllers"
	"simple-microservice/internal/persistence/models"
	"simple-microservice/internal/test/mock"
	"simple-microservice/internal/test/utils"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrDBConnectionLost      = errors.New("db connection lost")
	ErrDbConstraintViolation = errors.New("db constraint violation")
)

func TestGetUserById(t *testing.T) {
	t.Parallel()

	t.Run("GetUserById_ShouldReturnOk", func(t *testing.T) {
		t.Parallel()

		var (
			name       = "John Doe"
			email      = "john@doe.com"
			sampleUser = &models.User{
				ID:          1,
				ExternalId:  uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				Name:        &name,
				Email:       &email,
				DateOfBirth: time.Now(),
			}
		)

		repo := &mock.MockUserRepo{
			GetUserFunc: func(ctx context.Context, id uint64) (*models.User, error) {
				assert.EqualValues(t, 1, id)
				return sampleUser, nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		controller.GetUserById(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		resp := utils.DecodeResponse[utils.UserResponse](t, rec.Body.Bytes())

		assert.Equal(t, uint64(1), resp.Data.User.ID)
		assert.Equal(t, "John Doe", resp.Data.User.Name)
		assert.Equal(t, "john@doe.com", resp.Data.User.Email)
		assert.Equal(t,
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440000").String(),
			resp.Data.User.ExternalID,
		)
	})

	t.Run("GetUserById_NonNumericId_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			GetUserFunc: func(ctx context.Context, id uint64) (*models.User, error) {
				t.Fatal("GetUser should not be called")
				//nolint:nilnil
				return nil, nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("id", "abc")

		rec := httptest.NewRecorder()

		controller.GetUserById(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())

		assert.Equal(t, http.StatusText(http.StatusBadRequest), resp.Error)
		assert.Contains(t, resp.Reason, "abc")
	})

	t.Run("GetUserById_NonExistingId_ShouldReturnNotFound", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			GetUserFunc: func(ctx context.Context, id uint64) (*models.User, error) {
				assert.EqualValues(t, 999, id)
				return nil, gorm.ErrRecordNotFound
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("id", "999")

		rec := httptest.NewRecorder()

		controller.GetUserById(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusNotFound), resp.Error)
	})

	t.Run("GetUserById_InternalIssue_ShouldReturnInternalServerError", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			GetUserFunc: func(ctx context.Context, id uint64) (*models.User, error) {
				return nil, ErrDBConnectionLost
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		controller.GetUserById(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
	})
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("CreateUser_ShouldReturnCreated", func(t *testing.T) {
		t.Parallel()

		var createdID uint64 = 42

		expectedDOB, err := time.Parse(time.RFC3339, "1990-01-01T00:00:00Z")
		require.NoError(t, err)

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				assert.Equal(t, "John Doe", *user.Name)
				assert.Equal(t, "john@example.com", *user.Email)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", user.ExternalId.String())
				assert.Equal(t, expectedDOB, user.DateOfBirth)
				user.ID = createdID

				return nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		body := `{
			"name": "John Doe",
			"email": "john@example.com",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		resp := utils.DecodeResponse[utils.UserResponse](t, rec.Body.Bytes())
		assert.Equal(t, "John Doe", resp.Data.User.Name)
		assert.Equal(t, "john@example.com", resp.Data.User.Email)
	})

	t.Run("CreateUser_InvalidRequest_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				t.Fatal("CreateUser should not be called")
				return nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bad json`))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusBadRequest), resp.Error)
		assert.Contains(t, resp.Reason, "invalid character")
	})

	t.Run("CreateUser_MissingRequiredField_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				t.Fatal("CreateUser should not be called")
				return nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		body := `{
			"name": "",
			"email": "john@example.com",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusBadRequest), resp.Error)
	})

	t.Run("CreateUser_InvalidEmail_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				t.Fatal("CreateUser should not be called")
				return nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		body := `{
			"name": "John Doe",
			"email": "not-an-email",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusBadRequest), resp.Error)
	})

	t.Run("CreateUser_InvalidUuid_ShouldReturnBadRequest", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				t.Fatal("CreateUser should not be called")
				return nil
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		body := `{
			"name": "John Doe",
			"email": "john@example.com",
			"external_id": "not-a-uuid",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusBadRequest), resp.Error)
	})

	t.Run("CreateUser_RepositoryError_ShouldReturnInternalServerError", func(t *testing.T) {
		t.Parallel()

		repo := &mock.MockUserRepo{
			CreateUserFunc: func(ctx context.Context, user *models.User) error {
				return ErrDbConstraintViolation
			},
		}

		controller := controllers.NewUsersController(validator.New(), repo, zap.NewNop())

		body := `{
			"name": "John Doe",
			"email": "john@example.com",
			"external_id": "550e8400-e29b-41d4-a716-446655440000",
			"date_of_birth": "1990-01-01T00:00:00Z"
		}`

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		controller.CreateUser(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp := utils.DecodeResponse[utils.ErrorResponse](t, rec.Body.Bytes())
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), resp.Error)
	})
}
