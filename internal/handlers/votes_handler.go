package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
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
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userService.GetUserByEmail(r.Context(), email)
	if err != nil {
		h.sendError(w, err, http.StatusInternalServerError)
		return
	}

	err = auth.Access(email, password, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write(auth.GenerateTokenHandler(email, user.Role.Name, user.ID, []byte(h.cfg.JwtKey)))
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
