package account

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, acc *Account) error {
	query := `
		INSERT INTO accounts (user_id, account_number, balance, account_type_id, status_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		acc.UserID,
		acc.AccountNumber,
		acc.Balance,
		acc.AccountTypeID,
		acc.StatusID,
	).Scan(&acc.ID, &acc.CreatedAt, &acc.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Account, error) {
	var acc Account
	err := r.db.GetContext(ctx, &acc, "SELECT * FROM accounts WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// Deposit updates balance
func (r *Repository) Deposit(ctx context.Context, accountID int64, amount decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = now()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, amount, accountID)
	return err
}

// Withdraw ensures no negative balance
func (r *Repository) Withdraw(ctx context.Context, accountID int64, amount decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = balance - $1, updated_at = now()
		WHERE id = $2 AND balance >= $1
	`
	res, err := r.db.ExecContext(ctx, query, amount, accountID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("insufficient balance")
	}
	return nil
}

func (r *Repository) AccountsByUserID(user_id string) ([]ResponseAccount, error) {
	var accounts []ResponseAccount
	query := `
        SELECT 
			users.id AS user_id,
            users.full_name AS full_name,
            users.email AS email,
            accounts.account_number AS account_number,
            accounts.balance AS balance,
            account_types.name AS account_type,
            account_status.name AS status
        FROM accounts
        JOIN users ON users.id = accounts.user_id
        JOIN account_types ON account_types.id = accounts.account_type_id
        JOIN account_status ON account_status.id = accounts.status_id
        WHERE users.id = $1;`
	err := r.db.Select(&accounts, query, user_id)
	if err != nil {
		return []ResponseAccount{}, err
	}
	return accounts, nil
}

func (r *Repository) GetAllUserWithAccounts() ([]UserAccount, error) {
	var usersAccount []UserAccount
	query := `SELECT
	        u.id AS user_id,
	        u.email as email,
	        a.account_number as account_number,
	        a.balance as balance,
	        at.name AS account_type,
	        account_status.name AS status
	    FROM users u
	    LEFT JOIN accounts a ON u.id = a.user_id
	    LEFT JOIN account_types at ON a.account_type_id = at.id
	    LEFT JOIN account_status ON a.status_id = account_status.id order by u.id;`

	err := r.db.Select(&usersAccount, query)
	if err != nil {
		return []UserAccount{}, err
	}
	return usersAccount, nil
}
