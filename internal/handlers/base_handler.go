package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"go.uber.org/zap"
)

type BaseHandler struct {
	logger *zap.SugaredLogger
}

func NewBaseHandler(logger *zap.SugaredLogger) *BaseHandler {
	return &BaseHandler{
		logger: logger,
	}
}

func (h *BaseHandler) sendError(w http.ResponseWriter, err error, httpStatus int) {
	h.logger.Error(err.Error())
	h.respond(w, &ErrorResponse{Message: err.Error()}, httpStatus)
}

func (h *BaseHandler) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *BaseHandler) respond(w http.ResponseWriter, data interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func (h *BaseHandler) GetAuthenticatedUserID(ctx context.Context) string {
	ID, _ := ctx.Value(models.IDContextKey).(string)
	return ID
}

func (h *BaseHandler) GetAuthenticatedRole(ctx context.Context) string {
	role, _ := ctx.Value(models.RoleContextKey).(string)
	return role
}
