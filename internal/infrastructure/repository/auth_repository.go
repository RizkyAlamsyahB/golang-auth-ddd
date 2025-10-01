package repository

import (
	"context"
	"database/sql"

	"github.com/rizkyalamsyahb/library-golang/internal/domain/auth"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auth.Repository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Create(ctx context.Context, user *auth.User) error {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *AuthRepository) FindByID(ctx context.Context, id int64) (*auth.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ?`

	user := &auth.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?`

	user := &auth.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) Update(ctx context.Context, user *auth.User) error {
	query := `UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.ID)
	return err
}

func (r *AuthRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}