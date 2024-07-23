package server

import (
	"encoding/json"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
)

type ErrorResponse struct {
	message string
}

func (srv *server) createUserHandler() http.HandlerFunc {
	type CreateUserRequest struct {
		FirstName string
		LastName  string
		Nickname  string
	}

	type CreateUserResponse struct {
		UserId string `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		createUserRequest := &CreateUserRequest{}
		err := srv.decode(r, createUserRequest)
		if err != nil {
			srv.logger.Error(err)
			srv.respond(w, &ErrorResponse{message: err.Error()}, http.StatusBadRequest)
			return
		}

		userService := services.NewUserService(srv.db, srv.logger)
		user := &models.User{
			FirstName: createUserRequest.FirstName,
			LastName:  createUserRequest.LastName,
			Nickname:  createUserRequest.Nickname,
		}
		userId, err := userService.CreateUser(r.Context(), user)
		if err != nil {
			srv.logger.Error(err)
			appErrors := err.(*apperrors.AppError)
			srv.respond(w, &ErrorResponse{message: appErrors.Message}, http.StatusInternalServerError)
			return
		}

		createUserResponse := &CreateUserResponse{UserId: userId}
		srv.respond(w, createUserResponse, http.StatusCreated)
	}
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
