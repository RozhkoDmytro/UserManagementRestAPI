package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
)

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
}

func (srv *server) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
	}

	userId, err := srv.userService.CreateUser(r.Context(), user)
	if err != nil {
		srv.sendError(w, err, http.StatusInternalServerError)
		return
	}

	createUserResponse := &CreateUserResponse{UserId: userId}
	srv.respond(w, createUserResponse, http.StatusCreated)
}

func (srv *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		UserID    uint      `json:"user_id"`
		DeletedAt time.Time `json:"deleted_at"`
	}

	vars := mux.Vars(r)
	userID := vars["id"]

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

func (srv *server) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := srv.userService.GetUser(r.Context(), userID)
	if err != nil {
		srv.sendError(w, err, http.StatusNotFound)
		return
	}

	srv.respond(w, user, http.StatusCreated)
}

func (srv *server) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	createUserRequest := &CreateUserRequest{}
	err := srv.decode(r, createUserRequest)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}
	// Validate the User struct
	err = srv.validate.Struct(createUserRequest)
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

	_, err = srv.userService.UpdateUser(r.Context(), userID, updatedData)
	if err != nil {
		srv.sendError(w, err, http.StatusNotFound)
		return
	}

	srv.respond(w, nil, http.StatusCreated)
}

func (srv *server) listUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	page := queryParams.Get("page")
	pageSize := queryParams.Get("page_size")

	intPage, intPageSize, err := srv.validateListUsersParam(page, pageSize)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}
	users, err := srv.userService.ListUsers(r.Context(), intPage, intPageSize)
	if err != nil {
		srv.sendError(w, err, http.StatusBadRequest)
		return
	}

	srv.respond(w, users, http.StatusOK)
}

func (srv *server) validateListUsersParam(page, pageSize string) (validPage, validPageSize int, err error) {
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

func (srv *server) countUsers(w http.ResponseWriter, r *http.Request) {
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

func (srv *server) login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := srv.userService.GetUserByEmail(r.Context(), email)
	if err != nil {
		srv.sendError(w, err, http.StatusInternalServerError)
		return
	}

	err = auth.Access(email, password, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write(auth.GenerateTokenHandler(email, user.Role.Name, []byte(srv.cfg.JwtKey)))
}

func (srv *server) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (srv *server) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		srv.logger.Error(err)
	}
}

func (srv *server) ValidateUserStruct(ctx context.Context, createUserRequest *CreateUserRequest) error {
	// Validate the User struct
	err := srv.validate.Struct(createUserRequest)
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

func (srv *server) sendError(w http.ResponseWriter, err error, httpStatus int) {
	srv.logger.Error(err)
	srv.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}
