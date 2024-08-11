package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

type userHandler struct {
	userService services.UserServiceInterface
	logger      *zap.SugaredLogger
	validator   *validator.Validate
	cfg         *config.Config
}

func NewUserHandler(userService services.UserServiceInterface, logger *zap.SugaredLogger, validator *validator.Validate, cfg *config.Config) *userHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
		validator:   validator,
		cfg:         cfg,
	}
}

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 1000
)

type ErrorResponse struct {
	Message string `json:"message"`
}
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required,min=8,password"`
	RoleID    uint   `json:"role_id" validate:"required,oneof=1 2 3"`
}

func (srv *userHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		UserId string `json:"user_id"`
	}

	createUserRequest := &CreateUserRequest{}
	err := srv.decode(r, createUserRequest)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	err = srv.ValidateUserStruct(r.Context(), createUserRequest)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return

	}

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:     createUserRequest.Email,
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Password:  hash,
		RoleID:    1, // bad approach
	}

	userId, err := srv.userService.CreateUser(r.Context(), user)
	if err != nil {
		srv.sendError(w, err, http.StatusInternalServerError)
		return
	}

	createUserResponse := &CreateUserResponse{UserId: userId}
	srv.respond(w, createUserResponse, http.StatusCreated)
}

func (srv *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		UserID    uint      `json:"user_id"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	ctx := r.Context()
	role := roleFromContext(ctx)

	if role != models.StrAdmin {
		srv.sendError(w, errors.New("premission is denided"), http.StatusBadRequest)
		return
	}

	user, err := srv.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		srv.sendError(w, err, http.StatusInternalServerError)
		return
	}
	res := &CreateUserResponse{
		UserID:    user.ID,
		DeletedAt: user.DeletedAt,
	}
	srv.respond(w, res, http.StatusOK)
}

func (srv *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	ctx := r.Context()
	role := roleFromContext(ctx)

	if role != models.StrAdmin && userID != vars["id"] {
		srv.sendError(w, errors.New("premission is denided"), http.StatusBadRequest)
		return
	}

	createUserRequest := &CreateUserRequest{}
	err := srv.decode(r, createUserRequest)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}
	// Validate the User struct
	err = srv.validator.Struct(createUserRequest)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	updatedData := &models.User{
		Email:     createUserRequest.Email,
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Password:  hash,
	}

	if role == models.StrAdmin && createUserRequest.RoleID > 0 {
		updatedData.RoleID = createUserRequest.RoleID
	}

	_, err = srv.userService.UpdateUser(ctx, userID, updatedData)
	if err != nil {
		srv.sendError(w, err, http.StatusNotFound)
		return
	}

	srv.respond(w, nil, http.StatusCreated)
}

func (srv *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	ctx := r.Context()

	user, err := srv.userService.GetUser(ctx, userID)
	if err != nil {
		srv.sendError(w, err, http.StatusNotFound)
		return
	}

	srv.respond(w, user, http.StatusCreated)
}

func (srv *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ctx := r.Context()

	page := queryParams.Get("page")
	pageSize := queryParams.Get("page_size")

	intPage, intPageSize, err := srv.validateListUsersParam(page, pageSize)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}
	users, err := srv.userService.ListUsers(ctx, intPage, intPageSize)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	srv.respond(w, users, http.StatusOK)
}

func (srv *userHandler) CountUsers(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		Count uint `json:"count"`
	}
	ctx := r.Context()
	count, err := srv.userService.CountUsers(ctx)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}
	res := &CreateUserResponse{
		Count: uint(count),
	}
	srv.respond(w, res, http.StatusOK)
}

func (srv *userHandler) validateListUsersParam(page, pageSize string) (validPage, validPageSize int, err error) {
	validPage, err = strconv.Atoi(page)
	if err != nil {
		if page != "" {
			srv.logger.Error(err)
		}
		validPage = defaultPage
	}

	validPageSize, err = strconv.Atoi(pageSize)
	if err != nil {
		if pageSize != "" {
			srv.logger.Error(err)
		}
		validPageSize = defaultPageSize
	}

	if validPage < defaultPage {
		return validPage, validPageSize, errors.New("incorrect page number")
	}

	if validPageSize > maxPageSize || validPageSize <= 0 {
		return validPage, validPageSize, errors.New("the number of objects on the page should be in the range from 1 to " + strconv.Itoa(maxPageSize))
	}
	return validPage, validPageSize, nil
}

func (srv *userHandler) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (srv *userHandler) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		srv.logger.Error(err)
	}
}

func (srv *userHandler) ValidateUserStruct(ctx context.Context, createUserRequest *CreateUserRequest) error {
	// Validate the User struct
	err := srv.validator.Struct(createUserRequest)
	if err != nil {
		return err
	}

	// Check if email is unique
	existingUser, _ := srv.userService.GetUserByEmail(ctx, createUserRequest.Email)
	if existingUser != nil {
		err := errors.New("email already in use")
		return err
	}
	return nil
}

func (srv *userHandler) sendError(w http.ResponseWriter, err error, httpStatus int) {
	srv.logger.Error(err)
	srv.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}

func roleFromContext(ctx context.Context) string {
	role, _ := ctx.Value(models.RoleContextKey).(string)
	return role
}

func IDFromContext(ctx context.Context) string {
	ID, _ := ctx.Value(models.IDContextKey).(string)
	return ID
}
