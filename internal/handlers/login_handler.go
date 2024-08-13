package handlers

import (
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gitlab.com/jkozhemiaka/web-layout/internal/services"
	"go.uber.org/zap"
)

type loginHandler struct {
	*BaseHandler
	userService services.UserServiceInterface
	logger      *zap.SugaredLogger
	cfg         *config.Config
}

func NewLoginHandler(userService services.UserServiceInterface, logger *zap.SugaredLogger, cfg *config.Config) *loginHandler {
	return &loginHandler{
		BaseHandler: NewBaseHandler(logger),
		userService: userService,
		logger:      logger,
		cfg:         cfg,
	}
}

func (h *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
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
