package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rizkyalamsyahb/library-golang/internal/interfaces/http/response"
	"github.com/rizkyalamsyahb/library-golang/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserEmailKey contextKey = "user_email"

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "unauthorized", "missing authorization header")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "unauthorized", "invalid authorization format")
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "unauthorized", err.Error())
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to get user ID from context
func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}