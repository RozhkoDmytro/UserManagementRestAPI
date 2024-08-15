package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	myValidate "gitlab.com/jkozhemiaka/web-layout/internal/validate"
	"go.uber.org/zap"
)

func TestCreateUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	reqBody := &CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password@123",
		RoleID:    1,
	}

	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Mock the service response
	mockUserService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(uint(12345), nil)
	mockUserService.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(nil, nil)

	handler.CreateUserHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	type CreateUserResponse struct {
		UserId string `json:"user_id"`
	}

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	var response CreateUserResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, "12345", response.UserId)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	req := httptest.NewRequest(http.MethodDelete, "/users/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	w := httptest.NewRecorder()

	// Мокаємо роль адміністратора
	ctx := context.WithValue(req.Context(), models.RoleContextKey, models.StrAdmin)
	req = req.WithContext(ctx)

	// Мокаємо відповідь сервісу
	deletedUser := &models.User{ID: 123, DeletedAt: time.Now()}
	mockUserService.EXPECT().DeleteUser(gomock.Any(), "123").Return(deletedUser, nil)

	handler.DeleteUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	w := httptest.NewRecorder()

	// Мокаємо відповідь сервісу
	expectedUser := &models.User{ID: 123, Email: "test@example.com"}
	mockUserService.EXPECT().GetUser(gomock.Any(), "123").Return(expectedUser, nil)

	handler.GetUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var user models.User
	_ = json.NewDecoder(res.Body).Decode(&user)
	assert.Equal(t, uint(123), user.ID)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	// Мокаємо відповідь сервісу
	users := []models.User{
		{ID: 1, Email: "test1@example.com"},
		{ID: 2, Email: "test2@example.com"},
	}
	mockUserService.EXPECT().ListUsers(gomock.Any(), defaultPage, defaultPageSize).Return(users, nil)

	handler.ListUsers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var returnedUsers []*models.User
	_ = json.NewDecoder(res.Body).Decode(&returnedUsers)
	assert.Len(t, returnedUsers, 2)
}

func TestCountUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	req := httptest.NewRequest(http.MethodGet, "/users/count", nil)
	w := httptest.NewRecorder()

	// Мокаємо відповідь сервісу
	mockUserService.EXPECT().CountUsers(gomock.Any()).Return(123, nil)

	handler.CountUsers(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	type CreateUserResponse struct {
		Count int `json:"count"`
	}
	var response CreateUserResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	assert.Equal(t, 123, response.Count)
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)

	logger := zap.NewExample().Sugar()
	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	cfg := &config.Config{}

	handler := NewUserHandler(mockUserService, logger, validate, cfg)

	reqBody := &CreateUserRequest{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password@123",
		RoleID:    1,
	}

	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	w := httptest.NewRecorder()

	// Мокаємо роль адміністратора
	ctx := context.WithValue(req.Context(), models.RoleContextKey, models.StrAdmin)
	req = req.WithContext(ctx)

	// Мокаємо відповідь сервісу
	mockUserService.EXPECT().UpdateUser(gomock.Any(), "123", gomock.Any()).Return(nil, nil)

	handler.UpdateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusCreated, res.StatusCode)
}
