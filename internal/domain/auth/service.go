package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Register - business logic untuk register user baru
func (s *Service) Register(ctx context.Context, name, email, password string) (*User, error) {
	// 1. Check if email already exists
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user entity
	user := &User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 4. Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Clear password before returning
	user.ClearPassword()

	return user, nil
}

// Login - business logic untuk authenticate user
func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	// 1. Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Clear password before returning
	user.ClearPassword()

	return user, nil
}

// GetUserByID - get user by ID
func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	user.ClearPassword()
	return user, nil
}