package myValidate

import (
	"regexp"

	"github.com/go-playground/validator"
)

// Custom password validation
func Password(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	hasNumber := false
	hasSpecial := false
	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasNumber = true
		case regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(string(char)):
			hasSpecial = true
		}
	}
	return hasNumber && hasSpecial
}
