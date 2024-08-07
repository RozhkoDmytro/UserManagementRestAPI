package server

import (
	"bytes"
	"context"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/database"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	myValidate "gitlab.com/jkozhemiaka/web-layout/internal/validate"
	"go.uber.org/zap"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	domain  = "@example.com"
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
	userService := services.NewUserService(repositories.NewUserRepo(db, logger.Sugar()), logger.Sugar())

	// Initialize validator
	validate := validator.New()
	validate.RegisterValidation("password", myValidate.Password)

	srvRouter := &router{mux: mux.NewRouter()}
	srv := &server{
		db:          db,
		router:      srvRouter,
		logger:      logger.Sugar(),
		validate:    validate,
		cfg:         cfg,
		userService: userService,
	}
	return mockUserService, srv
}

// GenerateEmail creates a random email address
func GenerateEmail() string {
	// Create a new random source with the current time as a seed
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := strings.Builder{}
	for i := 0; i < 10; i++ {
		b.WriteByte(letters[r.Intn(len(letters))])
	}
	return b.String() + domain
}

func TestCreateUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)
	newEmail := ""
	for i := 0; i < 10; i++ {
		t.Run("success", func(t *testing.T) {
			newEmail = GenerateEmail()
			token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))
			mockUserService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return("2", nil)

			reqBody := `{"email":"` + GenerateEmail() + `","first_name":"John","last_name":"Doe","password":"passwoSrd123!", "role_id":3}`
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
			req.Header.Set("Authorization", "Bearer "+string(token))
			w := httptest.NewRecorder()

			srv.createUserHandler(w, req)
			srv.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
		})
	}
	t.Run("invalid request", func(t *testing.T) {
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

		reqBody := `{"email":"` + GenerateEmail() + `","first_name":"John","last_name":"Doe","password":"passwo", "role_id":3}`
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", "Bearer "+string(token))
		w := httptest.NewRecorder()

		srv.createUserHandler(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}

// Similar tests for getUserHandler, updateUserHandler, listUsersHandler

func TestGetUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	t.Run("success", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))
		mockUserService.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&models.User{
			ID:        7,
			Email:     newEmail,
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/user/8", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		srv.getUser(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))
		mockUserService.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, &apperrors.AppError{Message: "user not found"})

		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		srv.getUser(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Add more test cases as needed
}

func TestUpdateUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	t.Run("success", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

		mockUserService.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.User{
			ID:        7,
			Email:     newEmail,
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			RoleID:    3,
		}, nil)

		reqBody := `{"email":"` + newEmail + `","first_name":"John","last_name":"Doe","password":"passwor!!Gd123", "role_id":3}`
		req := httptest.NewRequest(http.MethodPut, "/users/8", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		srv.updateUser(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

		reqBody := `{"email":"` + GenerateEmail() + `"}`
		req := httptest.NewRequest(http.MethodPut, "/users/8", bytes.NewBufferString(reqBody))
		req.Header.Set("Authorization", "Bearer "+string(token))
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		w := httptest.NewRecorder()

		srv.updateUser(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Add more test cases as needed
}

func TestListUsersHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)

	t.Run("success", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

		mockUserService.EXPECT().ListUsers(gomock.Any(), "1", "10").Return([]*models.User{
			{
				ID:        1,
				Email:     GenerateEmail(),
				FirstName: "John",
				LastName:  "Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Email:     GenerateEmail(),
				FirstName: "Jane",
				LastName:  "Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users?page=1&page_size=10", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		srv.listUsers(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid request, but return ok!", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

		req := httptest.NewRequest(http.MethodGet, "/users?page=bad&page_size=10", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		srv.listUsers(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Add more test cases as needed
}

func TestDeleteUserHandler(t *testing.T) {
	mockUserService, srv := InitializeMock(t)
	/*
		t.Run("success", func(t *testing.T) {
			newEmail := GenerateEmail()
			token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))

			mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(&models.User{ID: 9, DeletedAt: time.Now()}, nil)

			req := httptest.NewRequest(http.MethodDelete, "/users/9", nil)
			req.Header.Set("Authorization", "Bearer "+string(token))
			req = mux.SetURLVars(req, map[string]string{"id": "9"})
			ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
			ctx = context.WithValue(ctx, EmailContextKey, newEmail)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			srv.deleteUser(w, req)
			srv.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}) */

	t.Run("not found", func(t *testing.T) {
		newEmail := GenerateEmail()
		token := auth.GenerateTokenHandler(newEmail, "admin", []byte(srv.cfg.JwtKey))
		mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil, &apperrors.AppError{Message: "user not found"})

		req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
		req.Header.Set("Authorization", "Bearer "+string(token))
		ctx := context.WithValue(req.Context(), RoleContextKey, "admin")
		ctx = context.WithValue(ctx, EmailContextKey, newEmail)
		req = req.WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		srv.deleteUser(w, req)
		srv.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// Add more test cases as needed
}

func TestLogin(t *testing.T) {
	_, srv := InitializeMock(t)

	t.Run("success", func(t *testing.T) {
		newEmail := "admin@example.com" // Use a fixed email for predictability

		// Encode the form values
		form := url.Values{}
		form.Add("email", newEmail)
		form.Add("password", "securePassword7!")
		reqBody := form.Encode()

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		srv.login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
