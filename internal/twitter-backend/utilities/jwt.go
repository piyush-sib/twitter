package utilities

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecretKey = "twitterAspire"

type JWTUtils struct {
	secret []byte
}

func NewJWTUtils() *JWTUtils {
	return &JWTUtils{secret: []byte(jwtSecretKey)}
}

func (j *JWTUtils) GenerateJWT(userID int, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateJWT validates a JWT token and returns its claims.
func (j *JWTUtils) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
