package usecase

import (
	"context"

"github.com/rizkyalamsyah_dev/library-golang/internal/application/dto"
    "github.com/rizkyalamsyah_dev/library-golang/internal/domain/auth"
    "github.com/rizkyalamsyah_dev/library-golang/pkg/jwt"
)

type AuthUseCase struct {
	authService *auth.Service
	jwtService  *jwt.JWTService
}

func NewAuthUseCase(authService *auth.Service, jwtService *jwt.JWTService) *AuthUseCase {
	return &AuthUseCase{
		authService: authService,
		jwtService:  jwtService,
	}
}

// Register handles user registration
func (uc *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Call domain service to register
	user, err := uc.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Convert entity to DTO
	response := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

// Login handles user authentication and generates JWT token
func (uc *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Authenticate user
	user, err := uc.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Prepare response
	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return response, nil
}

// GetProfile gets user profile by ID
func (uc *AuthUseCase) GetProfile(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := uc.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}