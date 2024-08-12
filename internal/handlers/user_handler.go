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

func (h *userHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		UserId string `json:"user_id"`
	}

	createUserRequest := &CreateUserRequest{}
	err := h.decode(r, createUserRequest)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	err = h.ValidateUserStruct(r.Context(), createUserRequest)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return

	}

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	user := &models.User{
		Email:     createUserRequest.Email,
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Password:  hash,
		RoleID:    1, // bad approach
	}

	userId, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		h.sendError(w, err, http.StatusInternalServerError)
		return
	}

	createUserResponse := &CreateUserResponse{UserId: userId}
	h.respond(w, createUserResponse, http.StatusCreated)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		UserID    uint      `json:"user_id"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	vars := mux.Vars(r)
	userID := vars["id"]

	ctx := r.Context()
	role := roleFromContext(ctx)

	if role != models.StrAdmin {
		h.sendError(w, errors.New("premission is denided"), http.StatusBadRequest)
		return
	}

	user, err := h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		h.sendError(w, err, http.StatusInternalServerError)
		return
	}
	res := &CreateUserResponse{
		UserID:    user.ID,
		DeletedAt: user.DeletedAt,
	}
	h.respond(w, res, http.StatusOK)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	ctx := r.Context()
	role := roleFromContext(ctx)

	if role != models.StrAdmin && userID != vars["id"] {
		h.sendError(w, errors.New("premission is denided"), http.StatusBadRequest)
		return
	}

	createUserRequest := &CreateUserRequest{}
	err := h.decode(r, createUserRequest)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}
	// Validate the User struct
	err = h.validator.Struct(createUserRequest)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
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

	_, err = h.userService.UpdateUser(ctx, userID, updatedData)
	if err != nil {
		h.sendError(w, err, http.StatusNotFound)
		return
	}

	h.respond(w, nil, http.StatusCreated)
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	ctx := r.Context()

	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		h.sendError(w, err, http.StatusNotFound)
		return
	}

	h.respond(w, user, http.StatusCreated)
}

func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ctx := r.Context()

	page := queryParams.Get("page")
	pageSize := queryParams.Get("page_size")

	intPage, intPageSize, err := h.validateListUsersParam(page, pageSize)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}
	users, err := h.userService.ListUsers(ctx, intPage, intPageSize)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}

	h.respond(w, users, http.StatusOK)
}

func (h *userHandler) CountUsers(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		Count uint `json:"count"`
	}
	ctx := r.Context()
	count, err := h.userService.CountUsers(ctx)
	if err != nil {
		h.sendError(w, err, http.StatusBadRequest)
		return
	}
	res := &CreateUserResponse{
		Count: uint(count),
	}
	h.respond(w, res, http.StatusOK)
}

func (h *userHandler) validateListUsersParam(page, pageSize string) (validPage, validPageSize int, err error) {
	validPage, err = strconv.Atoi(page)
	if err != nil {
		if page != "" {
			h.logger.Error(err)
		}
		validPage = defaultPage
	}

	validPageSize, err = strconv.Atoi(pageSize)
	if err != nil {
		if pageSize != "" {
			h.logger.Error(err)
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

func (h *userHandler) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *userHandler) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *userHandler) ValidateUserStruct(ctx context.Context, createUserRequest *CreateUserRequest) error {
	// Validate the User struct
	err := h.validator.Struct(createUserRequest)
	if err != nil {
		return err
	}

	// Check if email is unique
	existingUser, _ := h.userService.GetUserByEmail(ctx, createUserRequest.Email)
	if existingUser != nil {
		err := errors.New("email already in use")
		return err
	}
	return nil
}

func (h *userHandler) sendError(w http.ResponseWriter, err error, httpStatus int) {
	h.logger.Error(err)
	h.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}

func roleFromContext(ctx context.Context) string {
	role, _ := ctx.Value(models.RoleContextKey).(string)
	return role
}

func IDFromContext(ctx context.Context) string {
	ID, _ := ctx.Value(models.IDContextKey).(string)
	return ID
}
