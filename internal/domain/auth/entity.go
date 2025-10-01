package auth

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // "-" means won't be serialized to JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Business logic methods
func (u *User) IsValid() bool {
	return u.Name != "" && u.Email != "" && u.Password != ""
}

func (u *User) ClearPassword() {
	u.Password = ""
}