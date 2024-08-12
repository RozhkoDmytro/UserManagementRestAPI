package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

type votesHandler struct {
	userService services.UserServiceInterface
	logger      *zap.SugaredLogger
	cfg         *config.Config
}

func NewVotesHandler(userService services.UserServiceInterface, logger *zap.SugaredLogger, cfg *config.Config) *votesHandler {
	return &votesHandler{
		userService: userService,
		logger:      logger,
		cfg:         cfg,
	}
}

func (h *votesHandler) Like(w http.ResponseWriter, r *http.Request) {
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

func (h *votesHandler) DisLike(w http.ResponseWriter, r *http.Request) {
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

func (h *votesHandler) sendError(w http.ResponseWriter, err error, httpStatus int) {
	h.logger.Error(err)
	h.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}

func (h *votesHandler) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		h.logger.Error(err)
	}
}
