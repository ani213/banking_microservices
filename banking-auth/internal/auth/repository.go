package auth

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(u *User) error {
	_, err := r.db.Exec(`
        INSERT INTO users ( email, password_hash, full_name) 
        VALUES ($1, $2, $3)`, u.Email, u.PasswordHash, u.FullName)
	return err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		fmt.Println("Error fetching user by email:", err)

		return nil, err
	}
	return &u, err
}
