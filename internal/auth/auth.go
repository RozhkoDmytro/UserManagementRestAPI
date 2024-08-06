package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateTokenHandler(username, role string, JwtKey []byte) []byte {
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return []byte(tokenString)
}

func Access(username, password string, user *models.User) error {
	// Check the password
	if user == nil {
		return errors.New("Username is not fount in DB")
	}

	if !passwords.CheckPasswordHash(password, user.Password) {
		return errors.New("Invalid username or password")
	}

	return nil
}
