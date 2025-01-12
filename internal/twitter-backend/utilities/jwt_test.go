package utilities

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT_Success(t *testing.T) {
	jwtUtils := NewJWTUtils()
	token, err := jwtUtils.GenerateJWT(1, "testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT_Success(t *testing.T) {
	jwtUtils := NewJWTUtils()
	token, _ := jwtUtils.GenerateJWT(1, "testuser")
	claims, err := jwtUtils.ValidateJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), claims["user_id"])
	assert.Equal(t, "testuser", claims["name"])
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	jwtUtils := NewJWTUtils()
	claims, err := jwtUtils.ValidateJWT("invalidToken")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	jwtUtils := &JWTUtils{secret: []byte(jwtSecretKey)}
	claims := jwt.MapClaims{
		"user_id": 1,
		"name":    "testuser",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtUtils.secret)

	parsedClaims, err := jwtUtils.ValidateJWT(tokenString)

	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}
