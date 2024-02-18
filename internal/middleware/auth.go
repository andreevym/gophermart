// middleware/auth.go

package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/andreevym/gophermart/internal/services"
	"github.com/andreevym/gophermart/pkg/logger"
	"go.uber.org/zap"
)

type ContextKey int

const (
	UserIDContextKey ContextKey = iota
)

var ErrAuthUnauthorized = errors.New("unauthorized")

// AuthMiddleware is a middleware for authentication using JWT tokens.
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware creates a new instance of AuthMiddleware with the given AuthService.
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService}
}

// WithAuthentication implements the http.HandlerFunc interface for the AuthMiddleware.
func (am *AuthMiddleware) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/api/user/register") ||
			strings.Contains(r.RequestURI, "/api/user/login") {
			next.ServeHTTP(w, r)
			return
		}

		// Extract the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		// Validate the token and extract user ID
		userID, err := am.authService.ValidateToken(tokenString)
		if err != nil {
			logger.Logger().Warn("authService.ValidateToken", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the user ID from the token in the request context
		ctx := setUserID(r, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setUserID(r *http.Request, userID int64) context.Context {
	return context.WithValue(r.Context(), UserIDContextKey, userID)
}

func GetUserID(ctx context.Context) (int64, error) {
	userID := ctx.Value(UserIDContextKey)
	if userID == nil {
		return -1, ErrAuthUnauthorized
	}
	return userID.(int64), nil
}
