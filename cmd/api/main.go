package main

import (
	"fmt"
	"log"

	"github.com/rizkyalamsyah_dev/library-golang/config"
	"github.com/rizkyalamsyah_dev/library-golang/internal/application/usecase"
	"github.com/rizkyalamsyah_dev/library-golang/internal/domain/auth"
	"github.com/rizkyalamsyah_dev/library-golang/internal/infrastructure/database"
	httpRouter "github.com/rizkyalamsyah_dev/library-golang/internal/infrastructure/http"
	"github.com/rizkyalamsyah_dev/library-golang/internal/infrastructure/repository"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/handler"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/middleware"
	"github.com/rizkyalamsyah_dev/library-golang/pkg/jwt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("🚀 Starting %s...\n", cfg.App.Name)

	// Initialize database
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)

	// Initialize domain services
	authService := auth.NewService(authRepo)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(authService, jwtService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup router
	router := httpRouter.NewRouter(authHandler, authMiddleware)
	engine := router.Setup()

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("✅ Server running on http://%s\n", addr)
	fmt.Println("📚 API Endpoints:")
	fmt.Println("   POST   http://localhost:8080/api/auth/register")
	fmt.Println("   POST   http://localhost:8080/api/auth/login")
	fmt.Println("   GET    http://localhost:8080/api/auth/profile (Protected)")
	fmt.Println("   GET    http://localhost:8080/health")

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
