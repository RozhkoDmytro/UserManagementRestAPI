package server

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

func (srv *server) basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		authType, authValue, ok := strings.Cut(authHeader, " ")
		if !ok || authType != "Basic" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		decodedValue, err := base64.StdEncoding.DecodeString(authValue)
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusUnauthorized)
			return
		}

		userPass := strings.SplitN(string(decodedValue), ":", 2)
		if len(userPass) != 2 || userPass[0] != srv.cfg.Baseauth.Username || userPass[1] != srv.cfg.Baseauth.Password {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(ctx) // Use the returned request with the new context
		h(w, r)
	}
}

func (srv *server) contextExpire(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		r = r.WithContext(ctx) // Use the returned request with the new context
		h(w, r)
	}
}
