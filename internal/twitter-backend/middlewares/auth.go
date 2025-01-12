// AuthMiddleware validates the JWT token and adds user info to the context.
package middleware

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/dig"
	"net/http"
	"strings"
	"twitter/internal/twitter-backend/utilities"
)

type AuthMiddlewares struct {
	validateJWT func(tokenString string) (jwt.MapClaims, error)
}
type AuthMiddlewaresParams struct {
	dig.In

	JwtUtils *utilities.JWTUtils
}

func NewAuthMiddlewares(params AuthMiddlewaresParams) *AuthMiddlewares {
	return &AuthMiddlewares{
		validateJWT: params.JwtUtils.ValidateJWT,
	}
}

// AuthMiddleware validates the JWT token and adds user info to the context.
func (a *AuthMiddlewares) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate the token
			claims, err := a.validateJWT(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}
			userIDInt := claims["user_id"]
			// Check and convert userID to int to match with mysql db column type
			if floatVal, ok := claims["user_id"].(float64); ok {
				userIDInt = int(floatVal)
			}
			// Add user info to the context
			ctx := context.WithValue(r.Context(), "user_id", userIDInt)
			ctx = context.WithValue(ctx, "name", claims["name"])

			// Pass the request to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves user details from the context.
func (a *AuthMiddlewares) GetUserFromContext(ctx context.Context) (userID int, username string, ok bool) {
	userIDRaw := ctx.Value("user_id")
	usernameRaw := ctx.Value("name")

	if userIDRaw == nil || usernameRaw == nil {
		return 0, "", false
	}

	userID, ok = userIDRaw.(int)
	username, _ = usernameRaw.(string)

	return userID, username, true
}
