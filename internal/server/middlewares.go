package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
)

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
		ID := strconv.FormatUint(uint64(claims.ID), 10)
		if claims.Role == "" || claims.Email == "" || ID == "" {
			http.Error(w, "token haven't info about Role,Email,ID", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.RoleContextKey, claims.Role)
		ctx = context.WithValue(ctx, models.EmailContextKey, claims.Email)
		ctx = context.WithValue(ctx, models.IDContextKey, ID)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
