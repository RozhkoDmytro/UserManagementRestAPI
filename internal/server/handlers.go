package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
)

type ErrorResponse struct {
	message string
}

func (srv *server) createUserHandler() http.HandlerFunc {
	type CreateUserRequest struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
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

		hash, err := passwords.HashPassword(createUserRequest.Password)
		if err != nil {
			srv.logger.Error(err)
			srv.respond(w, &ErrorResponse{message: err.Error()}, http.StatusBadRequest)
			return
		}

		userService := services.NewUserService(srv.db, srv.logger)
		user := &models.User{
			Email:     createUserRequest.Email,
			FirstName: createUserRequest.FirstName,
			LastName:  createUserRequest.LastName,
			Password:  hash,
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

func (srv *server) deleteUser() http.HandlerFunc {
	type CreateUserResponse struct {
		UserID    uint
		DeletedAt time.Time
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		userService := services.NewUserService(srv.db, srv.logger)
		user, err := userService.DeleteUser(r.Context(), userID)
		if err != nil {
			srv.logger.Error(err)
			appErrors := err.(*apperrors.AppError)
			srv.respond(w, &ErrorResponse{message: appErrors.Message}, http.StatusInternalServerError)
			return
		}
		res := &CreateUserResponse{
			UserID:    user.ID,
			DeletedAt: user.DeletedAt,
		}
		srv.respond(w, res, http.StatusOK)
	}
}

func (srv *server) getUser() http.HandlerFunc {
	type CreateUserResponse struct {
		UserID    uint
		Email     string
		FirstName string
		LastName  string

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		userService := services.NewUserService(srv.db, srv.logger)
		user, err := userService.GetUser(r.Context(), userID)
		if err != nil {
			srv.logger.Error(err)
			appErrors := err.(*apperrors.AppError)
			srv.respond(w, &ErrorResponse{message: appErrors.Message}, http.StatusInternalServerError)
			return
		}

		res := &CreateUserResponse{
			UserID:    user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		srv.respond(w, res, http.StatusCreated)
	}
}

func (srv *server) updateUser() http.HandlerFunc {
	type CreateUserRequet struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		createUserRequest := &CreateUserRequet{}
		err := srv.decode(r, createUserRequest)
		if err != nil {
			srv.logger.Error(err)
			srv.respond(w, &ErrorResponse{message: err.Error()}, http.StatusBadRequest)
			return
		}

		hash, err := passwords.HashPassword(createUserRequest.Password)
		if err != nil {
			srv.logger.Error(err)
			srv.respond(w, &ErrorResponse{message: err.Error()}, http.StatusBadRequest)
			return
		}

		updatedData := &models.User{
			Email:     createUserRequest.Email,
			FirstName: createUserRequest.FirstName,
			LastName:  createUserRequest.LastName,
			Password:  hash,
		}

		userService := services.NewUserService(srv.db, srv.logger)
		_, err = userService.UpdateUser(r.Context(), userID, updatedData)
		if err != nil {
			srv.logger.Error(err)
			appErrors := err.(*apperrors.AppError)
			srv.respond(w, &ErrorResponse{message: appErrors.Message}, http.StatusInternalServerError)
			return
		}

		srv.respond(w, nil, http.StatusCreated)
	}
}

func (srv *server) listUsers() http.HandlerFunc {
	type CreateUserResponse struct {
		UserID    uint
		Email     string
		FirstName string
		LastName  string

		CreatedAt time.Time
		UpdatedAt time.Time
	}
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		page := queryParams.Get("page")
		pageSize := queryParams.Get("pageSize")

		userService := services.NewUserService(srv.db, srv.logger)
		users, err := userService.ListUsers(r.Context(), page, pageSize)
		if err != nil {
			srv.logger.Error(err)
			appErrors := err.(*apperrors.AppError)
			srv.respond(w, &ErrorResponse{message: appErrors.Message}, http.StatusInternalServerError)
			return
		}

		// convert models into struct
		var res []CreateUserResponse
		for _, item := range users {
			res = append(res, CreateUserResponse{
				UserID:    item.ID,
				Email:     item.Email,
				FirstName: item.FirstName,
				LastName:  item.LastName,

				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}
		srv.respond(w, res, http.StatusOK)
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
