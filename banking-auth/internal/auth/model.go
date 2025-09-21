package auth

import (
	"time"
)

type User struct {
	ID           int64     `db:"id"`
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

type EmailRequestBody struct {
	To      string `json:"to"`
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

type UserRoleRequestBody struct {
	UserID int64   `json:"userId" validate:"required"`
	Roles  []int64 `json:"roles" validate:"required"`
}
