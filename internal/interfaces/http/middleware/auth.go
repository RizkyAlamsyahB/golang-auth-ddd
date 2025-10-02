package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/response"
	"github.com/rizkyalamsyah_dev/library-golang/pkg/jwt"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorGin(c, http.StatusUnauthorized, "unauthorized", "missing authorization header")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorGin(c, http.StatusUnauthorized, "unauthorized", "invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			response.ErrorGin(c, http.StatusUnauthorized, "unauthorized", err.Error())
			c.Abort()
			return
		}

		// Add user info to Gin context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		// Continue to next handler
		c.Next()
	}
}

// Helper function to get user ID from Gin context
func GetUserIDFromGin(c *gin.Context) int64 {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}

	if id, ok := userID.(int64); ok {
		return id
	}

	return 0
}
