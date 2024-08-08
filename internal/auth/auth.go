package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/passwords"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	ID    uint   `json:"user_id"`
	jwt.StandardClaims
}

func GenerateTokenHandler(email, role string, ID uint, JwtKey []byte) []byte {
	claims := &Claims{
		Email: email,
		Role:  role,
		ID:    ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return nil
	}

	return []byte(tokenString)
}

func Access(username, password string, user *models.User) error {
	// Check the password
	if user == nil {
		return errors.New("username is not fount in DB")
	}

	if !passwords.CheckPasswordHash(password, user.Password) {
		return errors.New("invalid username or password")
	}

	return nil
}
