package server

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gitlab.com/jkozhemiaka/web-layout/internal/auth"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
)

type CacheKeyGenerator func(r *http.Request) string

func (srv *server) contextExpire(h http.HandlerFunc, keyGen CacheKeyGenerator, cacheTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()
		if r.Method != http.MethodGet {
			r = r.WithContext(ctx)
			h(w, r)
			return
		}

		// Generate a cacheKey based on a custom function
		cacheKey := keyGen(r)

		cachedData, err := srv.cache.Get(ctx, cacheKey, cacheTTL)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cachedData))
			return
		}

		r = r.WithContext(ctx) // Use the returned request with the new context

		// Буфер для зберігання відповіді
		responseBuffer := new(bytes.Buffer)
		// Створюємо кастомний writer для зберігання відповіді в буфер
		bufferedWriter := &bufferedResponseWriter{
			ResponseWriter: w,
			buffer:         responseBuffer,
		}
		h(bufferedWriter, r)
		if bufferedWriter.statusCode == http.StatusOK || bufferedWriter.statusCode == http.StatusCreated {
			err := srv.cache.Set(ctx, cacheKey, responseBuffer.String(), cacheTTL)
			if err != nil {
				log.Printf("Error caching response: %v", err)
			}
		}
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

// bufferedResponseWriter використовується для зберігання тіла відповіді
type bufferedResponseWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	statusCode int
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	bw.buffer.Write(b)
	return bw.ResponseWriter.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(statusCode int) {
	bw.statusCode = statusCode
	bw.ResponseWriter.WriteHeader(statusCode)
}

func CacheGenId(r *http.Request) string {
	vars := mux.Vars(r)
	return "user:" + vars["id"]
}
