package server

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
)

// Define a custom type for the context key
type contextKey string

const (
	RoleContextKey  contextKey = "role"
	EmailContextKey contextKey = "email"
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
		if len(userPass) != 2 {
			http.Error(w, "Invalid username or password format", http.StatusUnauthorized)
			return
		}

		user, err := srv.userService.GetUserByEmail(ctx, userPass[0])
		if err != nil {
			srv.sendError(w, err, http.StatusBadRequest)
			return
		}

		err = auth.Access(userPass[0], userPass[1], user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
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

func (srv *server) jwtMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(srv.cfg.JwtKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RoleContextKey, claims.Role)
		ctx = context.WithValue(ctx, EmailContextKey, claims.Email)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
