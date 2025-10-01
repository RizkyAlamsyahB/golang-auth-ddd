package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rizkyalamsyahb/library-golang/internal/application/dto"
	"github.com/rizkyalamsyahb/library-golang/internal/application/usecase"
	"github.com/rizkyalamsyahb/library-golang/internal/interfaces/http/middleware"
	"github.com/rizkyalamsyahb/library-golang/internal/interfaces/http/response"
	"github.com/rizkyalamsyahb/library-golang/pkg/validator"
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
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	// Call use case
	user, err := h.authUseCase.Register(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "registration failed", err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "user registered successfully", user)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	// Call use case
	result, err := h.authUseCase.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "login failed", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "login successful", result)
}

// GetProfile handles GET /api/auth/profile (protected route)
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r.Context())
	if userID == 0 {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "invalid user")
		return
	}

	// Get user profile
	user, err := h.authUseCase.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found", err.Error())
		return
	}

	response.Success(w, http.StatusOK, "profile retrieved successfully", user)
}