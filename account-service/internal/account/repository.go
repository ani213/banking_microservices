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
