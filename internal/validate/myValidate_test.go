package myValidate

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("password", Password)

	tests := []struct {
		password string
		valid    bool
	}{
		{"", false},              // Empty password
		{"short", false},         // Too short
		{"n0specialchar", false}, // No special character
		{"N0Special!", true},     // Valid password
		{"12345678", false},      // Only numbers
		{"abcdefgh", false},      // Only letters
		{"abcd1234", false},      // Letters and numbers, no special character
		{"!@#$%^&*", false},      // Only special characters
		{"Abcdef1!", true},       // Valid password with mixed characters
	}

	for _, test := range tests {
		err := validate.Var(test.password, "password")
		if test.valid {
			assert.NoError(t, err, "Password %q should be valid", test.password)
		} else {
			assert.Error(t, err, "Password %q should be invalid", test.password)
		}
	}
}
