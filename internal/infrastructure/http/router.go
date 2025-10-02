package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/handler"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/middleware"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Router {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode) // Change to gin.DebugMode for development

	return &Router{
		engine:         gin.New(),
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Add default middleware
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1
	api := r.engine.Group("/api")

	// Auth routes (public)
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", r.authHandler.Register)
		authRoutes.POST("/login", r.authHandler.Login)
	}

	// Protected routes (require authentication)
	protectedRoutes := api.Group("/auth")
	protectedRoutes.Use(r.authMiddleware.Authenticate())
	{
		protectedRoutes.GET("/profile", r.authHandler.GetProfile)
	}

	return r.engine
}
