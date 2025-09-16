package auth

import (
	"time"
)

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	FullName     string    `db:"full_name"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"omitempty,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type ResponsGetUser struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	FullName string `db:"full_name"`
}
