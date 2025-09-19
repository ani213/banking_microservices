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
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	// Step 1 create user
	var newUser User

	err = tx.QueryRowx(`
        INSERT INTO users ( email, password_hash, full_name) 
        VALUES ($1, $2, $3) RETURNING *`, u.Email, u.PasswordHash, u.FullName).StructScan(&newUser)
	if err != nil {
		tx.Rollback()
		return err
	}
	// step 2 add role
	var role int = 1
	_, err = tx.Exec(`INSERT INTO user_role_mapping (user_id ,role_id) VALUES ($1, $2)`, newUser.ID, role)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
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

func (r *Repository) FindUsers() ([]ResponsGetUser, error) {
	var u []ResponsGetUser
	err := r.db.Select(&u, "SELECT id, email, full_name FROM users")
	if err != nil {
		fmt.Println("Error fetching users:", err)

		return nil, err
	}
	return u, err
}
