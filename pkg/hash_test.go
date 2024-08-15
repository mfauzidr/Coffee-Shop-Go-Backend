package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err, "An error occurred while hashing the password")
	assert.NotEmpty(t, hashedPassword, "The hashed password is empty/unfilled")
	assert.NotEqual(t, password, hashedPassword, "The hashed password should be different from the original password")
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)

	assert.NoError(t, err, "An error occurred while hashing the password")

	err = VerifyPassword(hashedPassword, password)
	assert.NoError(t, err, "An error occurred even though the password entered was correct when verifying password.")

	wrongPassword := "wrongpassword"
	err = VerifyPassword(hashedPassword, wrongPassword)
	assert.Error(t, err, "An error occurred even though the password entered was wrong when verifying password.")
}
