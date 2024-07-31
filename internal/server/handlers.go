package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
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
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	// Validate the User struct
	err = srv.validate.Struct(createUserRequest)
	if err != nil {

		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	/* 	// Check if email is unique
	   	existingUser, _ := srv.userService.GetUserByEmail(r.Context(), createUserRequest.Email)
	   	if existingUser != nil {
	   		http.Error(w, "Email already in use", http.StatusBadRequest)
	   		return
	   	} */

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
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
		srv.logger.Error(err)
		appErrors := err.(*apperrors.AppError)
		srv.respond(w, &ErrorResponse{Message: appErrors.Message}, http.StatusInternalServerError)
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
		srv.logger.Error(err)
		appErrors := err.(*apperrors.AppError)
		srv.respond(w, &ErrorResponse{Message: appErrors.Message}, http.StatusInternalServerError)
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
		srv.logger.Error(err)
		appErrors := err.(*apperrors.AppError)
		srv.respond(w, &ErrorResponse{Message: appErrors.Message}, http.StatusNotFound)
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
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	// Validate the User struct
	err = srv.validate.Struct(createUserRequest)
	if err != nil {
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	hash, err := passwords.HashPassword(createUserRequest.Password)
	if err != nil {
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
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
		srv.logger.Error(err)
		appErrors := err.(*apperrors.AppError)
		srv.respond(w, &ErrorResponse{Message: appErrors.Message}, http.StatusNotFound)
		return
	}

	srv.respond(w, nil, http.StatusCreated)
}

func (srv *server) listUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	page := queryParams.Get("page")
	pageSize := queryParams.Get("page_size")

	users, err := srv.userService.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	srv.respond(w, users, http.StatusOK)
}

func (srv *server) countUsers(w http.ResponseWriter, r *http.Request) {
	type CreateUserResponse struct {
		Count uint `json:"count"`
	}
	ctx := r.Context()
	count, err := srv.userService.CountUsers(ctx)
	if err != nil {
		srv.logger.Error(err)
		srv.respond(w, &ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	res := &CreateUserResponse{
		Count: uint(count),
	}
	srv.respond(w, res, http.StatusOK)
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
