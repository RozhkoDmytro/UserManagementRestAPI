package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

type loginHandler struct {
	userService services.UserServiceInterface
	logger      *zap.SugaredLogger
	cfg         *config.Config
}

func NewLoginHandler(userService services.UserServiceInterface, logger *zap.SugaredLogger, cfg *config.Config) *loginHandler {
	return &loginHandler{
		userService: userService,
		logger:      logger,
		cfg:         cfg,
	}
}

func (srv *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	w.Write(auth.GenerateTokenHandler(email, user.Role.Name, user.ID, []byte(srv.cfg.JwtKey)))
}

func (srv *loginHandler) sendError(w http.ResponseWriter, err error, httpStatus int) {
	srv.logger.Error(err)
	srv.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}

func (srv *loginHandler) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		srv.logger.Error(err)
	}
}
