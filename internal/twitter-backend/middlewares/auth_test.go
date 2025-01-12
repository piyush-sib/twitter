package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	authMiddleware := &AuthMiddlewares{
		validateJWT: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"user_id": 1, "name": "testuser"}, nil
		},
	}

	handler := authMiddleware.AuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id")
		name := r.Context().Value("name")
		assert.Equal(t, 1, userID)
		assert.Equal(t, "testuser", name)
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer validToken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	authMiddleware := &AuthMiddlewares{
		validateJWT: func(tokenString string) (jwt.MapClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := authMiddleware.AuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Authorization header missing\n", rr.Body.String())
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	authMiddleware := &AuthMiddlewares{
		validateJWT: func(tokenString string) (jwt.MapClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := authMiddleware.AuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid Authorization header format\n", rr.Body.String())
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	authMiddleware := &AuthMiddlewares{
		validateJWT: func(tokenString string) (jwt.MapClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	handler := authMiddleware.AuthMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalidToken")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid or expired token\n", rr.Body.String())
}

func TestGetUserFromContext_ValidContext(t *testing.T) {
	authMiddleware := &AuthMiddlewares{}
	ctx := context.WithValue(context.Background(), "user_id", 1)
	ctx = context.WithValue(ctx, "name", "testuser")

	userID, username, ok := authMiddleware.GetUserFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, 1, userID)
	assert.Equal(t, "testuser", username)
}

func TestGetUserFromContext_InvalidContext(t *testing.T) {
	authMiddleware := &AuthMiddlewares{}
	ctx := context.Background()

	userID, username, ok := authMiddleware.GetUserFromContext(ctx)

	assert.False(t, ok)
	assert.Equal(t, 0, userID)
	assert.Equal(t, "", username)
}
