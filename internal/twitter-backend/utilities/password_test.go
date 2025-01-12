package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordUtils_HashPassword_Success(t *testing.T) {
	passwordUtils := NewPasswordUtils()
	password := "piyush"

	hashedPassword, err := passwordUtils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestPasswordUtils_CheckPasswordHash_ValidPassword(t *testing.T) {
	passwordUtils := NewPasswordUtils()
	password := "piyush"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	isValid := passwordUtils.CheckPasswordHash(password, string(hashedPassword))

	assert.True(t, isValid)
}

func TestPasswordUtils_CheckPasswordHash_InvalidPassword(t *testing.T) {
	passwordUtils := NewPasswordUtils()
	password := "piyush"
	invalidPassword := "piyush2"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	isValid := passwordUtils.CheckPasswordHash(invalidPassword, string(hashedPassword))

	assert.False(t, isValid)
}
