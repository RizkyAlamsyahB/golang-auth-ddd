package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rizkyalamsyah_dev/library-golang/internal/application/dto"
	"github.com/rizkyalamsyah_dev/library-golang/internal/application/usecase"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/middleware"
	"github.com/rizkyalamsyah_dev/library-golang/internal/interfaces/http/response"
	"github.com/rizkyalamsyah_dev/library-golang/pkg/validator"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	validator   *validator.Validator
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator.New(),
	}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorGin(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		response.ErrorGin(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	// Call use case
	user, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		response.ErrorGin(c, http.StatusBadRequest, "registration failed", err.Error())
		return
	}

	response.SuccessGin(c, http.StatusCreated, "user registered successfully", user)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorGin(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		response.ErrorGin(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	// Call use case
	result, err := h.authUseCase.Login(c.Request.Context(), req)
	if err != nil {
		response.ErrorGin(c, http.StatusUnauthorized, "login failed", err.Error())
		return
	}

	response.SuccessGin(c, http.StatusOK, "login successful", result)
}

// GetProfile handles GET /api/auth/profile (protected route)
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == 0 {
		response.ErrorGin(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return
	}

	// Get user profile
	user, err := h.authUseCase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.ErrorGin(c, http.StatusNotFound, "user not found", err.Error())
		return
	}

	response.SuccessGin(c, http.StatusOK, "profile retrieved successfully", user)
}
