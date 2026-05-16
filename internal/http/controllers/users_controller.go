package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"simple-microservice/internal/http/dto/mappers"
	"simple-microservice/internal/http/dto/requests"
	"simple-microservice/internal/http/dto/responses"
	"simple-microservice/internal/persistence/models"
)

type UserRepository interface {
	GetUser(id uint64) (*models.User, error)
	CreateUser(user *models.User) error
}

type UsersController struct {
	validator  *validator.Validate
	repository UserRepository
	logger     *zap.Logger
}

func (u *UsersController) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		u.logger.Error("failed to parse id", zap.Error(err))

		code := http.StatusBadRequest
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	user, err := u.repository.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		u.logger.Error("failed to find user", zap.Uint64("id", id), zap.Error(err))

		code := http.StatusNotFound
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	if err != nil {
		u.logger.Error("error while trying to find user", zap.Uint64("id", id), zap.Error(err))

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	bytes, err := responses.DataBytes("user", mappers.MapUserToResponse(user))
	if err != nil {
		u.logger.Error("failed to marshal user", zap.Uint64("id", id), zap.Error(err))

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	_, _ = w.Write(bytes)
}

func (u *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		u.logger.Error("failed to read request body", zap.Error(err))

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, nil))

		return
	}

	defer func() {
		_ = body.Close()
	}()

	var request requests.CreateUserRequest

	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		u.logger.Error("failed to unmarshal request body", zap.Error(err))

		code := http.StatusBadRequest
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	err = u.validator.Struct(request)
	if err != nil {
		u.logger.Error("failed to validate request", zap.Error(err))

		code := http.StatusBadRequest
		w.WriteHeader(code)
		_, _ = w.Write(responses.ValidationErrorBytes(code, err))

		return
	}

	user, err := mappers.MapRequestToUser(request)
	if err != nil {
		u.logger.Error("failed to parse external id", zap.Error(err))

		code := http.StatusBadRequest
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	err = u.repository.CreateUser(user)
	if err != nil {
		u.logger.Error("failed to create user", zap.Error(err))

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	bytes, err := responses.DataBytes("user", mappers.MapUserToResponse(user))
	if err != nil {
		u.logger.Error("failed to marshal response", zap.Error(err))

		code := http.StatusInternalServerError
		w.WriteHeader(code)
		_, _ = w.Write(responses.ErrorBytes(code, err))

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(bytes)
}
