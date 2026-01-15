package account

import (
	"context"
	"errors"
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
func (r *Repository) Deposit(ctx context.Context, accountNumber string, amount float64) (BalancEmail, error) {
	// 1. Start a transation
	tx, err := r.db.Begin()
	if err != nil {
		return BalancEmail{}, err
	}
	// Ensure rollback if we return early due to an error
	defer tx.Rollback()
	// 2. Read with FOR UPDATE
	// This locks the account row until tx.Commit() or tx.Rollback()
	var statusName string
	var balance float64
	var email string
	query := `
        SELECT  t.name,a.balance,u.email
        FROM accounts a 
        JOIN account_status t ON a.status_id = t.id 
		JOIN users u on a.user_id =u.id 
        WHERE a.account_number = $1 
        FOR UPDATE OF a;` // Lock only the 'accounts' table row

	err = tx.QueryRow(query, accountNumber).Scan(&statusName, &balance, &email)
	if err != nil {
		return BalancEmail{}, err
	}

	if statusName == "active" {
		updateQuery := `UPDATE accounts SET balance = balance + $1 WHERE account_number = $2`
		_, err := tx.Exec(updateQuery, amount, accountNumber)
		if err != nil {
			return BalancEmail{}, err
		}

	} else {
		return BalancEmail{}, fmt.Errorf("account is not active")
	}
	// 4. Commit the transaction (this releases the lock)
	err = tx.Commit()
	if err != nil {
		return BalancEmail{}, err
	}
	return BalancEmail{Balance: balance + amount, Email: email, AccountNo: accountNumber}, nil
	// query:=`select a.account_number,t.name from accounts a join account_status t on a.status_id =t.id where a.account_number=$1;`

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

func (r *Repository) GetEmailByUserId(userId string) (string, error) {
	var email []string
	query := `select email from users where id=$1;`
	err := r.db.Select(&email, query, userId)
	if err != nil {
		return "", err
	}
	if len(email) == 0 {
		return "", errors.New("no email found")
	}
	return email[0], nil
}

func (r *Repository) UpdateBalance(accountNumber string, amount float64) error {
	// 1. Start a transation
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Ensure rollback if we return early due to an error
	defer tx.Rollback()
	// 2. Read with FOR UPDATE
	// This locks the account row until tx.Commit() or tx.Rollback()
	var statusName string
	var balance float64
	query := `
        SELECT  t.name,a.balance 
        FROM accounts a 
        JOIN account_status t ON a.status_id = t.id 
        WHERE a.account_number = $1 
        FOR UPDATE OF a;` // Lock only the 'accounts' table row

	err = tx.QueryRow(query, accountNumber).Scan(&statusName, &balance)
	if err != nil {
		return err
	}

	if statusName == "active" {
		updateQuery := `UPDATE accounts SET balance = balance + $1 WHERE account_number = $2`
		_, err := tx.Exec(updateQuery, amount, accountNumber)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("account is not active")
	}
	// 4. Commit the transaction (this releases the lock)
	return tx.Commit()
	// query:=`select a.account_number,t.name from accounts a join account_status t on a.status_id =t.id where a.account_number=$1;`

}
