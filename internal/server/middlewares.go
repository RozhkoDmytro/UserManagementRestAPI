package server

import (
	"context"
	"net/http"
	"time"
)

func (srv *server) contextExpire(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		r.WithContext(ctx)
		h(w, r)
	}
}
