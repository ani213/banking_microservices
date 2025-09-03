package auth

import (
	"database/sql"

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
        INSERT INTO users (id, email, password_hash) 
        VALUES ($1, $2, $3)`, u.ID, u.Email, u.PasswordHash)
	return err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}
