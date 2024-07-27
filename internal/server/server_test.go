package server

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	myValidate "gitlab.com/jkozhemiaka/web-layout/internal/validate"
	"go.uber.org/zap"
)

func InitializeMock(t *testing.T) (*services.MockUserServiceInterface, *server) {
	os.Setenv("CONFIG_PATH", "../../configs/.sample.env")
	defer os.Unsetenv("CONFIG_PATH")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := services.NewMockUserServiceInterface(ctrl)
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(apperrors.LoggerInitError.AppendMessage(err))
	}
	defer logger.Sync()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	db, err := database.SetupDatabase(cfg)
	if err != nil {
		logger.Sugar().Fatal(err)
	}

	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	srvRouter := &router{mux: mux.NewRouter()}
	srv := &server{
		db:       db,
		router:   srvRouter,
		logger:   logger.Sugar(),
		validate: validate,
		cfg:      cfg,
	}
	return mockUserService, srv
}

func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestCreateUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	handler := srv.createUserHandler()

	t.Run("success", func(t *testing.T) {
		mockUserService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return("1", nil)

		reqBody := `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":"passwoSrd123!"}`
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}

func TestDeleteUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	handler := srv.deleteUser()

	t.Run("success", func(t *testing.T) {
		mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(&models.User{ID: 7, DeletedAt: time.Now()}, nil)

		req := httptest.NewRequest(http.MethodDelete, "/user/7", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil, &apperrors.AppError{Message: "user not found"})

		req := httptest.NewRequest(http.MethodDelete, "/user/999", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// Add more test cases as needed
}

// Similar tests for getUserHandler, updateUserHandler, listUsersHandler

func TestGetUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	handler := srv.getUser()

	t.Run("success", func(t *testing.T) {
		mockUserService.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/user/8", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUserService.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, &apperrors.AppError{Message: "user not found"})

		req := httptest.NewRequest(http.MethodGet, "/user/999", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// Add more test cases as needed
}

func TestUpdateUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	handler := srv.updateUser()

	t.Run("success", func(t *testing.T) {
		mockUserService.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		reqBody := `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":"passwor!!Gd123"}`
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		reqBody := `{"email":"test@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}

func TestListUsersHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	handler := srv.listUsers()

	t.Run("success", func(t *testing.T) {
		mockUserService.EXPECT().ListUsers(gomock.Any(), "1", "10").Return([]*models.User{
			{
				ID:        1,
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Email:     "jane@example.com",
				FirstName: "Jane",
				LastName:  "Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users?page=1&page_size=10", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users?page=bad&page_size=10", nil)
		req.Header.Set("Authorization", encodeBasicAuth(srv.cfg.Baseauth.Username, srv.cfg.Baseauth.Password))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}
