package passwords

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "mySecureP@ssw0rd"

	hash, err := HashPassword(password)
	assert.NoError(t, err, "HashPassword should not return an error")

	assert.NotEmpty(t, hash, "Hash should not be empty")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecureP@ssw0rd"

	hash, err := HashPassword(password)
	assert.NoError(t, err, "HashPassword should not return an error")

	valid := CheckPasswordHash(password, hash)
	assert.True(t, valid, "Password should match the hash")

	invalid := CheckPasswordHash("wrongpassword", hash)
	assert.False(t, invalid, "Password should not match the hash")
}
