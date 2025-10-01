package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rizkyalamsyahb/library-golang/internal/interfaces/http/handler"
	"github.com/rizkyalamsyahb/library-golang/internal/interfaces/http/middleware"
)

type Router struct {
	mux            *mux.Router
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Router {
	return &Router{
		mux:            mux.NewRouter(),
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Setup() *mux.Router {
	// API v1
	api := r.mux.PathPrefix("/api").Subrouter()

	// Health check
	r.mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	// Auth routes (public)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", r.authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", r.authHandler.Login).Methods("POST")

	// Protected routes (require authentication)
	protectedRoutes := api.PathPrefix("/auth").Subrouter()
	protectedRoutes.Use(r.authMiddleware.Authenticate)
	protectedRoutes.HandleFunc("/profile", r.authHandler.GetProfile).Methods("GET")

	return r.mux
}